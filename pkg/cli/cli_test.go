package cli_test

import (
	"os"

	. "github.com/dgruber/qsub/pkg/cli"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cli", func() {

	Context("ParseCommandLine", func() {

		It("should set a job template, quiet, and synced", func() {
			cli, err := ParseCommandline([]string{"--sync", "--quiete", "-b",
				"process", "-j", "jobtemplate.json"})
			Expect(err).To(BeNil())
			Expect(cli.Sync).To(BeTrue())
			Expect(cli.Quiet).To(BeTrue())
			Expect(cli.Backend).To(Equal("process"))
			Expect(cli.JobTemplatePath).To(Equal("jobtemplate.json"))
		})

		It("should set the command and arguments", func() {
			cli, err := ParseCommandline([]string{"--image", "golang:latest",
				"cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(len(cli.Cmd)).To(BeNumerically("==", 2))
			Expect(cli.Cmd[0]).To(Equal("cmd"))
			Expect(cli.Cmd[1]).To(Equal("arg"))
		})

		It("should set other commands and arguments", func() {
			cli, err := ParseCommandline([]string{"--image", "busybox:latest",
				"sleep", "123"})
			Expect(err).To(BeNil())
			Expect(len(cli.Cmd)).To(BeNumerically("==", 2))
			Expect(cli.Cmd[0]).To(Equal("sleep"))
			Expect(cli.Cmd[1]).To(Equal("123"))
		})

		It("should set job container image", func() {
			cli, err := ParseCommandline([]string{"--image", "golang:latest",
				"cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Image).To(Equal("golang:latest"))

			os.Setenv("QSUB_IMAGE", "golang:latest2")
			defer os.Unsetenv("QSUB_IMAGE")
			cli, err = ParseCommandline([]string{"cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Image).To(Equal("golang:latest2"))
		})

		It("should set labels", func() {
			cli, err := ParseCommandline([]string{"--image", "golang:latest",
				"-l", "label1,label2", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(len(cli.Labels)).To(BeNumerically("==", 2))
			Expect(cli.Labels[0]).To(Equal("label1"))
			Expect(cli.Labels[1]).To(Equal("label2"))
		})

		It("should set jobname", func() {
			cli, err := ParseCommandline([]string{"--image", "golang:latest",
				"--jobname", "name", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Jobname).To(Equal("name"))

			cli, err = ParseCommandline([]string{"--image", "golang:latest",
				"-N", "name", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Jobname).To(Equal("name"))
		})

		It("should set hostname", func() {
			cli, err := ParseCommandline([]string{"--image", "golang:latest",
				"--hostname", "name", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Hostname).To(Equal("name"))

			cli, err = ParseCommandline([]string{"--image", "golang:latest",
				"-N", "name", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Jobname).To(Equal("name"))
		})

		It("should set the namespace", func() {
			cli, err := ParseCommandline([]string{"--image", "golang:latest",
				"--namespace", "name", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Namespace).To(Equal("name"))

			cli, err = ParseCommandline([]string{"--image", "golang:latest",
				"-S", "name", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Namespace).To(Equal("name"))
		})

		It("should set environment variables", func() {
			cli, err := ParseCommandline([]string{"--image", "golang:latest",
				"-v", "key=value,key2=value2", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Envs["key"]).To(Equal("value"))
			Expect(cli.Envs["key2"]).To(Equal("value2"))
			os.Setenv("TESTQSUBENV1", "X")
			defer os.Unsetenv("TESTQSUBENV1")
			os.Setenv("TESTQSUBENV2", "Y")
			defer os.Unsetenv("TESTQSUBENV2")
			cli, err = ParseCommandline([]string{"--image", "golang:latest",
				"-v", "TESTQSUBENV1,TESTQSUBENV2", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Envs["TESTQSUBENV1"]).To(Equal("X"))
			Expect(cli.Envs["TESTQSUBENV2"]).To(Equal("Y"))
		})

		It("should set scheduler", func() {
			cli, err := ParseCommandline([]string{"--image", "golang:latest",
				"--scheduler", "poseidon", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Scheduler).To(Equal("poseidon"))

			cli, err = ParseCommandline([]string{"--image", "golang:latest",
				"-N", "name", "cmd", "arg"})
			Expect(err).To(BeNil())
			Expect(cli.Scheduler).To(Equal(""))
		})

	})

	Context("ParseCommandLine errors", func() {

		It("should error when command is not set", func() {
			_, err := ParseCommandline([]string{})
			Expect(err).NotTo(BeNil())
			_, err = ParseCommandline(nil)
			Expect(err).NotTo(BeNil())
		})

		It("should error when image is not set", func() {
			_, err := ParseCommandline([]string{"cmd", "arg"})
			Expect(err).NotTo(BeNil())
		})

		It("should error when the jobname argument is not set", func() {
			_, err := ParseCommandline([]string{"--image", "golang:latest",
				"--jobname"})
			Expect(err).NotTo(BeNil())
		})

		It("should error when the hostname argument is not set", func() {
			_, err := ParseCommandline([]string{"--image", "golang:latest",
				"--hostname"})
			Expect(err).NotTo(BeNil())
		})

		It("should error when the image argument is not set", func() {
			_, err := ParseCommandline([]string{"--image"})
			Expect(err).NotTo(BeNil())
		})

		It("should error when the namespace argument is not set", func() {
			_, err := ParseCommandline([]string{"--image", "golang:latest",
				"--namespace"})
			Expect(err).NotTo(BeNil())
		})

		It("should error when the env argument is not set", func() {
			_, err := ParseCommandline([]string{"--image", "golang:latest",
				"-v"})
			Expect(err).NotTo(BeNil())
		})

		It("should error when the label argument is not set", func() {
			_, err := ParseCommandline([]string{"--image", "golang:latest",
				"-l"})
			Expect(err).NotTo(BeNil())
		})

		It("should error when the scheduler argument is not set", func() {
			_, err := ParseCommandline([]string{"--image", "golang:latest",
				"--scheduler"})
			Expect(err).NotTo(BeNil())
		})

		It("should error with unknown parameters", func() {
			_, err := ParseCommandline([]string{"--img", "golang:latest",
				"sleep", "123"})
			Expect(err).NotTo(BeNil())
		})

	})

	Context("Complex command line", func() {

		It("should find all settings", func() {
			cli, err := ParseCommandline([]string{"--image", "golang:latest",
				"--scheduler", "kube-batch", "--jobname", "jn", "sleep", "123",
				"-123"})
			Expect(err).To(BeNil())
			Expect(len(cli.Cmd)).To(BeNumerically("==", 3))
			Expect(cli.Cmd[0]).To(Equal("sleep"))
			Expect(cli.Cmd[1]).To(Equal("123"))
			Expect(cli.Jobname).To(Equal("jn"))
			Expect(cli.Image).To(Equal("golang:latest"))
			Expect(cli.Scheduler).To(Equal("kube-batch"))
		})

	})

})
