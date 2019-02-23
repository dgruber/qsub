package template

import (
	"errors"
	"strings"

	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/qsub/pkg/cli"
)

// Create returns a JobTemplate based on the given parameters.
func Create(request cli.Commandline) (t drmaa2interface.JobTemplate, err error) {
	if request.Image == "" {
		return t, errors.New("container image is not specified")
	}
	if request.Cmd == nil {
		return t, errors.New("Command to run is not specified")
	}
	t.JobCategory = request.Image
	t.JobName = request.Jobname
	if request.Hostname != "" {
		t.CandidateMachines = []string{request.Hostname}
	}
	if request.Namespace != "" {
		if t.ExtensionList == nil {
			t.ExtensionList = make(map[string]string)
		}
		t.ExtensionList["namespace"] = request.Namespace
	}
	if request.Labels != nil {
		if t.ExtensionList == nil {
			t.ExtensionList = make(map[string]string)
		}
		t.ExtensionList["labels"] = strings.Join(request.Labels, ",")
	}
	if request.Scheduler != "" {
		if t.ExtensionList == nil {
			t.ExtensionList = make(map[string]string)
		}
		t.ExtensionList["scheduler"] = request.Scheduler
	}
	t.JobEnvironment = request.Envs
	t.RemoteCommand = request.Cmd[0]
	if len(request.Cmd) > 1 {
		t.Args = request.Cmd[1:]
	}
	return t, err
}
