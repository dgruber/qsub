package main

import (
	"fmt"
	"os"

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
	jt, err := template.Create(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(2)
	}
	jobid, err := job.Submit(jt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(3)
	}
	fmt.Printf("Submitted job with ID %s\n", jobid)
}
