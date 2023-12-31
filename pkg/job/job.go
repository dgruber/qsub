package job

import (
	"context"
	"fmt"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client"
	genclient "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client/generated"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/simpletracker"
	"github.com/dgruber/qsub/pkg/cli"
	"github.com/dgruber/qsub/pkg/server"
	"github.com/dgruber/wfl"
	"github.com/dgruber/wfl/pkg/context/docker"
	"github.com/dgruber/wfl/pkg/context/googlebatch"
	"github.com/dgruber/wfl/pkg/context/kubernetes"

	cepubsub "github.com/cloudevents/sdk-go/protocol/pubsub/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func fp(e error) {
	fmt.Printf("Error: %s\n", e.Error())
	os.Exit(1)
}

// Submit runs a batch job in the k8s cluster the current context points to.
func Submit(jt drmaa2interface.JobTemplate) (*wfl.Workflow, *wfl.Job, error) {
	flow := wfl.NewWorkflow(kubernetes.NewKubernetesContext().
		OnError(fp)).OnError(fp)
	job, err := Run(flow, jt)
	return flow, job, err
}

func SubmitToBackend(request cli.Commandline, jt drmaa2interface.JobTemplate) (*wfl.Workflow, *wfl.Job, error) {
	switch request.Backend {
	case "kubernetes":
		return Submit(jt)
	case "process":
		// let's ignore the internal ID as it is not for
		// the user relevant
		return SubmitProcess(jt)
	case "docker":
		return SubmitDocker(jt)
	case "googlebatch":
		return SubmitGoogleBatch(jt)
	case "pubsub":
		return SubmitToPubSub(jt)
	case "server":
		return SubmitToQsubServer(request.ServeHost, request.ServePort, jt)
	case "mpioperator":
		return nil, nil, fmt.Errorf("not implemented yet")
	default:
		return nil, nil, fmt.Errorf("Backend %s not supported", request.Backend)
	}
}

func SubmitProcess(jt drmaa2interface.JobTemplate) (*wfl.Workflow, *wfl.Job, error) {
	// for better performance this could be tuned
	ctx := wfl.NewProcessContextByCfgWithInitParams(wfl.ProcessConfig{
		/*
			DBFile:               "drmaa2session.db",
			JobDBFile:            "drmaa2job.db",
			PersistentJobStorage: true,
		*/
	}, simpletracker.SimpleTrackerInitParams{
		/*
			UsePersistentJobStorage: true,
			DBFilePath:              "drmaa2job.db",
		*/
	}).OnError(fp)
	flow := wfl.NewWorkflow(ctx).OnError(fp)
	job, err := Run(flow, jt)
	return flow, job, err
}

func SubmitDocker(jt drmaa2interface.JobTemplate) (*wfl.Workflow, *wfl.Job, error) {
	flow := wfl.NewWorkflow(
		docker.NewDockerContext().
			OnError(fp)).OnError(fp)
	job, err := Run(flow, jt)
	return flow, job, err
}

func SubmitGoogleBatch(jt drmaa2interface.JobTemplate) (*wfl.Workflow, *wfl.Job, error) {
	region := os.Getenv("GOOGLE_REGION")
	if region == "" {
		region = "us-central1"
	}
	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		return nil, nil, fmt.Errorf("GOOGLE_PROJECT environment variable not set")
	}
	flow := wfl.NewWorkflow(googlebatch.NewGoogleBatchContext(
		region,
		project,
	).OnError(fp)).OnError(fp)
	job, err := Run(flow, jt)
	return flow, job, err
}

func SubmitToQsubServer(host string, port int, jt drmaa2interface.JobTemplate) (*wfl.Workflow, *wfl.Job, error) {
	password, err := server.GetOrCreateSecret()
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("Using password stored in ~/.qsub/secret")
	basicAuthProvider, err := securityprovider.NewSecurityProviderBasicAuth(
		"qsub", password)
	if err != nil {
		return nil, nil, err
	}

	srv := fmt.Sprintf("http://%s:%d",
		host, port)

	initParams := client.ClientTrackerParams{
		Server: srv,
		Path:   "/qsub",
		Opts: []genclient.ClientOption{
			genclient.WithRequestEditorFn(basicAuthProvider.Intercept),
		},
	}

	ctx := wfl.NewRemoteContext(wfl.RemoteConfig{}, &initParams)
	if ctx.HasError() {
		return nil, nil, ctx.CtxCreationErr
	}

	flow := wfl.NewWorkflow(ctx)
	if flow.HasError() {
		return nil, nil, flow.Error()
	}

	job, err := Run(flow, jt)
	return flow, job, err
}

// SubmitToPubSub sends job template a CloudEvent to a PubSub topic. It requires
// the job template to contain the extension "googleProjectID" which is the
// Google Cloud project ID. If the extension is not set the environment variable
// GOOGLE_PROJECT is used. The queueName of the job template is used as the
// PubSub topic name.
func SubmitToPubSub(jt drmaa2interface.JobTemplate) (*wfl.Workflow, *wfl.Job, error) {
	ctx := context.Background()

	googleProject := os.Getenv("GOOGLE_PROJECT")
	if jt.ExtensionList != nil {
		googleProject, _ = jt.ExtensionList["googleProjectID"]
	}
	if googleProject == "" {
		return nil, nil,
			fmt.Errorf("job template does not contain googleProjectID extension")
	}

	pubsubTopic := jt.QueueName
	if pubsubTopic == "" {
		return nil, nil,
			fmt.Errorf("job template does not contain the PubSub topic specified as the queueName")
	}

	sender, err := cepubsub.New(ctx, cepubsub.WithProjectID(googleProject),
		cepubsub.WithTopicID(pubsubTopic))
	if err != nil {
		return nil, nil,
			fmt.Errorf("failed to create pubsub transport: %v", err)
	}
	defer sender.Close(ctx)

	client, err := cloudevents.NewClient(sender,
		cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create pubsub client: %v", err)
	}

	evt := cloudevents.NewEvent()
	evt.SetType("org.drmaa2.events.jobtemplate")
	evt.SetSource("qsub")
	err = evt.SetData("application/json", jt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set CloudEvent data: %v", err)
	}

	result := client.Send(ctx, evt)
	if cloudevents.IsUndelivered(result) {
		return nil, nil, fmt.Errorf("failed to send CloudEvent to topic %s: %v",
			pubsubTopic, result)
	}

	return nil, nil, nil
}

func Run(flow *wfl.Workflow, jt drmaa2interface.JobTemplate) (*wfl.Job, error) {
	job := flow.RunT(jt)
	if job.Errored() {
		return nil, job.LastError()
	}
	return job, nil
}

func ReadJobTemplateFromJSONFile(path string) (drmaa2interface.JobTemplate, error) {
	var jt drmaa2interface.JobTemplate
	jtBytes, err := os.ReadFile(path)
	if err != nil {
		return jt, fmt.Errorf("could not read job template from file %s: %v",
			path, err)
	}
	err = jt.UnmarshalJSON(jtBytes)
	if err != nil {
		return jt, fmt.Errorf("could not unmarshal job template: %v", err)
	}
	return jt, nil
}
