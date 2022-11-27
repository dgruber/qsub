package job

import (
	"fmt"
	"os"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/wfl"
	"github.com/dgruber/wfl/pkg/context/docker"
	"github.com/dgruber/wfl/pkg/context/googlebatch"
	"github.com/dgruber/wfl/pkg/context/kubernetes"
)

func fp(e error) {
	fmt.Printf("Error: %s\n", e.Error())
	os.Exit(1)
}

// Submit runs a batch job in the k8s cluster the current context points to.
func Submit(jt drmaa2interface.JobTemplate) (string, error) {
	flow := wfl.NewWorkflow(kubernetes.NewKubernetesContext().OnError(fp)).OnError(fp)
	return Run(flow, jt)
}

func SubmitToBackend(backend string, jt drmaa2interface.JobTemplate) (string, error) {
	switch backend {
	case "kubernetes":
		return Submit(jt)
	case "process":
		return SubmitProcess(jt)
	case "docker":
		return SubmitDocker(jt)
	case "googlebatch":
		return SubmitGoogleBatch(jt)
	case "mpioperator":
		return "", fmt.Errorf("not implemented yet")
	default:
		return "", fmt.Errorf("Backend %s not supported", backend)
	}
}

func SubmitProcess(jt drmaa2interface.JobTemplate) (string, error) {
	flow := wfl.NewWorkflow(wfl.NewProcessContext().OnError(fp)).OnError(fp)
	return Run(flow, jt)
}

func SubmitDocker(jt drmaa2interface.JobTemplate) (string, error) {
	flow := wfl.NewWorkflow(docker.NewDockerContext().OnError(fp)).OnError(fp)
	return Run(flow, jt)
}

func SubmitGoogleBatch(jt drmaa2interface.JobTemplate) (string, error) {
	region := os.Getenv("GOOGLE_REGION")
	if region == "" {
		region = "us-central1"
	}
	project := os.Getenv("GOOGLE_PROJECT")
	if project == "" {
		return "", fmt.Errorf("GOOGLE_PROJECT environment variable not set")
	}
	flow := wfl.NewWorkflow(googlebatch.NewGoogleBatchContext(
		region,
		project,
	).OnError(fp)).OnError(fp)
	return Run(flow, jt)
}

func Run(flow *wfl.Workflow, jt drmaa2interface.JobTemplate) (string, error) {
	job := flow.RunT(jt)
	if job.Errored() {
		return "", job.LastError()
	}
	return job.JobID(), nil
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
