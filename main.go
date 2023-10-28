package main

import (
	"fmt"
	"os"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/qsub/pkg/cli"
	"github.com/dgruber/qsub/pkg/job"
	"github.com/dgruber/qsub/pkg/template"
)

func main() {
	request, err := cli.ParseCommandline(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to submit job: %s\n", err.Error())
		os.Exit(1)
	}
	var jt drmaa2interface.JobTemplate
	if request.Backend == "" {
		// default backend with all parameters from cli
		jt, err = template.Create(request)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(2)
		}
		request.Backend = "kubernetes"
	} else {
		// job template from file
		jt, err = job.ReadJobTemplateFromJSONFile(request.JobTemplatePath)
		if err != nil {
			fmt.Printf("Failed parsing job template: %s\n", err.Error())
			os.Exit(3)
		}
	}

	_, job, err := job.SubmitToBackend(request.Backend, jt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(4)
	}
	fmt.Printf("Job submitted\n")

	if job != nil {
		fmt.Printf("Job ID: %s\n", job.JobID())
	}
	if request.Sync {
		// wait for the job to finish
		if job == nil {
			fmt.Fprintf(os.Stderr,
				fmt.Sprintf("Backend %s does not yet support sync\n",
					request.Backend))
		} else {
			fmt.Printf("Waiting for job to finish...\n")
			// forwards the exit code of the job
			os.Exit(job.ExitStatus())
		}
	}
}
