package main

import (
	"fmt"
	"os"

	"github.com/deepmap/oapi-codegen/pkg/securityprovider"
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client"
	genclient "github.com/dgruber/drmaa2os/pkg/jobtracker/remote/client/generated"
	"github.com/dgruber/qsub/pkg/server"
	"github.com/dgruber/wfl"
	"github.com/jedib0t/go-pretty/v6/table"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Use --help for syntax.\n")
		os.Exit(1)
	}

	cli, err := ParseCommandline(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse command line: %v\n", err)
		os.Exit(1)
	}
	jobs, err := getStatus(cli.JobIDs, cli.RemoteHost, cli.RemotePort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get status: %v\n", err)
		os.Exit(1)
	}

	t := table.Table{}
	t.AppendHeader(table.Row{
		"JobID",
		"State",
		"ExitCode",
		"DispatchTime",
		"Runtime"})
	rows := []table.Row{}
	for _, job := range jobs {
		i := job.JobInfo()
		if job.State() == drmaa2interface.Failed ||
			job.State() == drmaa2interface.Done {
			rows = append(rows, table.Row{
				job.JobID(), i.State.String(),
				i.ExitStatus, i.DispatchTime.String(),
				i.FinishTime.Sub(i.DispatchTime).String(),
			})
		} else if job.State() == drmaa2interface.Queued {
			rows = append(rows, table.Row{
				job.JobID(), i.State.String(),
				"-", "-",
				"-",
			})
		} else {
			rows = append(rows, table.Row{
				job.JobID(), i.State.String(),
				"-", i.DispatchTime.String(),
				"-",
			})
		}
	}
	t.AppendRows(rows)
	fmt.Printf("%s\n", t.Render())
}

func getStatus(jobIDs []string, host string, port int) ([]*wfl.Job, error) {
	password, err := server.GetOrCreateSecret()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Using password stored in ~/.qsub/secret\n")
	basicAuthProvider, err := securityprovider.NewSecurityProviderBasicAuth(
		"qsub", password)
	if err != nil {
		return nil, err
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
		return nil, ctx.CtxCreationErr
	}

	flow := wfl.NewWorkflow(ctx)
	if flow.HasError() {
		return nil, flow.Error()
	}

	jobs := flow.ListJobs()
	if flow.HasError() {
		return nil, flow.Error()
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("no jobs found")
	}
	filteredJobs := []*wfl.Job{}
	filter := map[string]bool{}
	for _, job := range jobIDs {
		filter[job] = true
	}
	for _, job := range jobs {
		if filter[job.JobID()] {
			filteredJobs = append(filteredJobs, job)
		}
	}
	return filteredJobs, nil
}
