package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Commandline contains the batch job submission instructions.
type Commandline struct {
	Image     string
	Jobname   string
	Hostname  string
	Namespace string
	Envs      map[string]string
	Labels    []string
	Cmd       []string
	Scheduler string
}

// ParseCommandline takes command line args and parses them.
func ParseCommandline(args []string) (Commandline, error) {
	var (
		err error
		cli Commandline
	)

	if args == nil {
		return cli, errors.New("Use --help for syntax.")
	}

	cli.Image = os.Getenv("QSUB_IMAGE")

argumentLoop:
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help":
			fmt.Fprintf(os.Stdout, "%s", Help())
			os.Exit(0)
		case "--hostname":
			if len(args) <= i+1 {
				err = errors.New("Option hostname requires an argument.")
				break argumentLoop
			}
			cli.Hostname = args[i+1]
			i++
			continue
		case "-I", "--image":
			if len(args) <= i+1 {
				err = errors.New("Option image requires an argument.")
				break argumentLoop
			}
			cli.Image = args[i+1]
			i++
			continue
		case "-l":
			cli.Labels = []string{}
			if len(args) <= i+1 {
				err = errors.New("Option -l requires an argument.")
				break argumentLoop
			}
			for _, v := range strings.Split(args[i+1], ",") {
				cli.Labels = append(cli.Labels, v)
			}
			i++
			continue
		case "-N", "--jobname":
			if len(args) <= i+1 {
				err = errors.New("Option jobname requires an argument.")
				break argumentLoop
			}
			cli.Jobname = args[i+1]
			i++
			continue
		case "-S", "--namespace":
			if len(args) <= i+1 {
				err = errors.New("Option namespace requires an argument.")
				break argumentLoop
			}
			cli.Namespace = args[i+1]
			i++
			continue
		case "--scheduler":
			if len(args) <= i+1 {
				err = errors.New("Option scheduler requires an argument.")
				break argumentLoop
			}
			cli.Scheduler = args[i+1]
			i++
		case "-v":
			if len(args) <= i+1 {
				err = errors.New("Option -v requires an argument.")
				break argumentLoop
			}
			for _, v := range strings.Split(args[i+1], ",") {
				if cli.Envs == nil {
					cli.Envs = make(map[string]string)
				}
				kv := strings.Split(v, "=")
				if len(kv) == 1 {
					cli.Envs[kv[0]] = os.Getenv(kv[0])
				} else {
					cli.Envs[kv[0]] = kv[1]
				}
			}
			i++
			continue
		default:
			if strings.HasPrefix(args[i], "-") {
				err = fmt.Errorf("Unknown argument %s", args[i])
				break argumentLoop
			}
			cli.Cmd = args[i:]
			break argumentLoop
		}
	}
	if cli.Cmd == nil {
		return cli, errors.New("No command given.")
	}
	if cli.Image == "" {
		return cli, errors.New("No container image given.")
	}
	return cli, err
}

// Help returns the help message.
func Help() string {
	usage := `qsub is a tool for submitting batch jobs to kubernetes.

Usage:
	Setting a container image and a command is mandatory. The container
	image can also be set by a QSUB_IMAGE environment variable.

	qsub [-N | --jobname unique_name_of_job]
	   	[-S | --namespace kubernetes_namespace]
		[-v env=content,...]
		[-l label1,...]
		[-I | --image] container_image
	   	command [args...]
`
	return usage
}
