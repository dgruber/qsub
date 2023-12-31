package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Commandline struct {
	RemoteHost string
	RemotePort int
	JobIDs     []string
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

argumentLoop:
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help":
			fmt.Fprintf(os.Stdout, "%s", "Help()")
			os.Exit(0)
		case "--server":
			if len(args) <= i+1 {
				err = errors.New("Option server requires an argument.")
				break argumentLoop
			}
			s := strings.Split(args[i+1], ":")
			if len(s) == 1 {
				cli.RemoteHost = s[0]
				cli.RemotePort = 13177
			} else if len(s) == 2 {
				cli.RemoteHost = s[0]
				cli.RemotePort, err = strconv.Atoi(s[1])
				if err != nil {
					err = fmt.Errorf("Failed to parse port number: %s",
						err.Error())
					break argumentLoop
				}
			} else {
				err = fmt.Errorf("Failed to parse host and port: %s",
					args[i+1])
				break argumentLoop
			}
			i++
			continue
		}
		cli.JobIDs = append(cli.JobIDs, args[i:]...)
	}
	return cli, err
}
