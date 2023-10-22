# qsub _not just_ for Kubernetes

[![CircleCI](https://circleci.com/gh/dgruber/qsub.svg?style=svg)](https://circleci.com/gh/dgruber/qsub)
[![codecov](https://codecov.io/gh/dgruber/qsub/branch/master/graph/badge.svg)](https://codecov.io/gh/dgruber/qsub)

_qsub_ is a command line tool for submitting batch jobs to a
workload manager. Its basic functionality is described and specified in
the [POSIX standard](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/qsub.html). Several HPC job schedulers (like Grid Engine, SLURM, LSF, Torque) provide
a _qsub_ command line utility enhanced with an [huge amount of extensions](http://gridengine.eu/mangridengine/manuals.html).

This repository provides a simple qsub implementation for local processes, Docker,
Kubernetes, Google Batch, PubSub, ...

## Installation

_qsub_ can be build directly (GO111MODULE=on) from the sources or alternatively pre-build binaries for darwin and linux can be downloaded from the builds dir.

## Usage for Kubernetes

In order to let the _sleep_ command run in the Kubernetes cluster
which is selected in the current context (_./kube/config_) you need
to specify the container image as well as the command to be executed
in the container (potentially with its arguments).

    qsub --image busybox:latest sleep 123

It returns the ID of the job.

The container image can also be set beforehand as environment variable.

    export QSUB_IMAGE=busybox:latest
    qsub sleep 123

The corresponding pods can be showed with the ID returned back on command line:

    kubectl describe pod -l job-name=ID

### More Arguments for Kubernetes

The image name (_--img_) as well as the command which should be executed in the
container derived from the image are mandatory.

Following optional arguments are currently available.

#### Job Name

The job name must be unique otherwise job submission will fail.

    qsub --image busybox:latest -N unique sleep 123
    Submitted job with ID unique

    kubectl get jobs unique
    NAME   DESIRED   SUCCESSFUL   AGE
    unique 1         0            14s   

#### Namespace for Kubernetes Jobs

Jobs can be submitted to a specific kubernetes namespace.

    qsub --namespace default --image busybox:latest sleep 123

#### Labels for Kubernetes Jobs

Kubernetes allows to attach labels to pods. Labels are key-value pairs
which can be defined with the _-l_ argument.

#### Environment Variables for Kubernetes Jobs

Environment variables for the jobs can be set by passing them into _-v_. There are
two ways of doing so: as key-value pairs, or just by name.

Following example sets ENV1 to VALUE1 and ENV2 to VALUE2 (using the value from 
the current context).

    export ENV2=VALUE2
    qsub --image busybox:latest -v ENV1=VALUE1,ENV2 sleep 123

#### Scheduler for Kubernetes Jobs

In order to let the job be scheduled by a non-default scheduler (like poseidon 
or kube-batch) the _--scheduler_ argument can be used.

    qsub --scheduler poseidon --img busybox:latest sleep 123

## Using DRMAA2 JSON File

By using a [DRMAA2 compatible](https://github.com/dgruber/drmaa2interface) JSON file the
jobs can be submitted to other backends as well:

- Google Batch (-b googlebatch)
- Local Process (-b process)
- Docker (-b docker)
- Sending a DRMAA2 JobTemplated wrapped as CloudEvent into Google PubSub (-b pubsub)
