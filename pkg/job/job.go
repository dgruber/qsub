package job

import (
	"github.com/dgruber/drmaa2interface"
	"github.com/dgruber/wfl"
	"github.com/dgruber/wfl/pkg/context/kubernetes"
)

func fp(e error) {
	panic(e)
}

// Submit runs a batch job in the k8s cluster the current context points to.
func Submit(t drmaa2interface.JobTemplate) (string, error) {
	job := wfl.NewWorkflow(kubernetes.NewKubernetesContext().OnError(fp)).OnError(fp).RunT(t)
	if job.Errored() {
		return "", job.LastError()
	}
	return job.JobID(), nil
}
