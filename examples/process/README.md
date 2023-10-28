# Running Jobs as OS Processes

## Starting a process which runs in background

    qsub -b process -j ./process.json
    ps

## Waiting until job is finished

The _-s_ flag also forwards the exit code of the process to the local shell.

    qsub -b process -s -j ./counter.json
    echo $?

## Pull a Docker image and execute it

The Docker backend does not pull a container image hence it needs to be pulled
before using the process backend.

    qsub -b process -s -j pull.json && qsub -b docker -s -j container.json
