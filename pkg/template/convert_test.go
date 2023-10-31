package template_test

import (
	. "github.com/dgruber/qsub/pkg/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/dgruber/qsub/pkg/cli"
)

var _ = Describe("Convert", func() {

	Context("JobTemplate created", func() {
		var request cli.Commandline

		BeforeEach(func() {
			request = cli.Commandline{
				Image: "image",
				Cmd:   []string{"cmd"},
			}
		})

		It("should contain RemoteCommand and Args", func() {
			jt, err := Create(request)
			Expect(err, BeNil())
			Expect(jt.RemoteCommand, Equal("cmd"))
			Expect(jt.Args, BeNil())

			request.Cmd = []string{"cmd", "arg1", "arg2"}
			jt, err = Create(request)
			Expect(err, BeNil())
			Expect(jt.RemoteCommand, Equal("cmd"))
			Expect(len(jt.Args), BeNumerically("==", 2))
			Expect(jt.Args[0], Equal("arg1"))
			Expect(jt.Args[1], Equal("arg2"))

			request.Cmd = nil
			jt, err = Create(request)
			Expect(err, Not(BeNil()))
		})

		It("should error when image is not specified", func() {
			request = cli.Commandline{Cmd: []string{"cmd"}}
			_, err := Create(request)
			Expect(err, Not(BeNil()))
		})

		It("should set the JobName", func() {
			request.Jobname = "jobname"
			j, err := Create(request)
			Expect(err, Not(BeNil()))
			Expect(j.JobName, Equal("jobname"))
		})

		It("should set the HostName", func() {
			request.Hostname = "host"
			j, err := Create(request)
			Expect(err, Not(BeNil()))
			Expect(j.CandidateMachines[0], Equal("host"))
		})

		It("should set the namespace", func() {
			request.Namespace = "namespace1"
			j, err := Create(request)
			Expect(err, BeNil())
			Expect(j.ExtensionList["namespace"], Equal("namespace1"))
		})

		It("should set the labels", func() {
			request.Labels = []string{"key1=label1", "key2=label2"}
			j, err := Create(request)
			Expect(err, BeNil())
			Expect(j.ExtensionList["labels"], Equal("key1=label1,key2=label2"))
		})

		It("should set the environment variables", func() {
			request.Envs = map[string]string{"key1": "env", "key2": "label2"}
			j, err := Create(request)
			Expect(err, BeNil())
			Expect(j.JobEnvironment["key1"], Equal("env"))
		})
	})

})
