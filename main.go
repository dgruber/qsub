package main

import (
	"fmt"
	"os"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/qsub/pkg/cli"
	"github.com/dgruber/qsub/pkg/job"
	"github.com/dgruber/qsub/pkg/server"
	"github.com/dgruber/qsub/pkg/template"
)

func main() {
	request, err := cli.ParseCommandline(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to submit job: %s\n", err.Error())
		os.Exit(1)
	}

	if request.Client == false && request.ServeHost != "" {
		// start qsub as server

		password, err := server.GetOrCreateSecret()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Qsub server failed: %s\n",
				err.Error())
			os.Exit(1)
		}

		err = server.Serve(server.Config{
			Host:     request.ServeHost,
			Port:     request.ServePort,
			Password: password,
			Backend:  request.Backend,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Qsub server failed: %s\n",
				err.Error())
			os.Exit(1)
		}
		return
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

	_, job, err := job.SubmitToBackend(request, jt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(4)
	}
	if !request.Quiet {
		fmt.Printf("Job submitted\n")
	}

	if job != nil && !request.Quiet {
		fmt.Printf("Job ID: %s\n", job.JobID())
	}
	if request.Sync {
		// wait for the job to finish
		if job == nil {
			fmt.Fprintf(os.Stderr,
				fmt.Sprintf("Backend %s does not yet support sync\n",
					request.Backend))
		} else {
			if !request.Quiet {
				fmt.Printf("Waiting for job to finish...\n")
			}
			// forwards the exit code of the job + 128
			exitCode := job.ExitStatus()
			if exitCode == 0 {
				os.Exit(0)
			}
			os.Exit(exitCode + 128)
		}
	}
}
