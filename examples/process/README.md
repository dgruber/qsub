# Running Jobs as OS Processes

## Starting a process which runs in background

    qsub -b process -j ./process.json
    ps

## Waiting until job is finished

The _-s_ flag also forwards the exit code of the process to the local shell.
If the application terminates with exit code > 0 then 128 is added. Hence 
an exit code of 129 means that the application terminated with exit code 1.

    qsub -b process -s -j ./counter.json
    echo $?

## Pull a Docker image and execute it

The Docker backend does not pull a container image hence it needs to be pulled
before using the process backend.

    qsub -b process -s -j pull.json && qsub -b docker -s -j container.json

## Pipe stdout to stdin

Piping the output from one process to the other. Here _--quiet_ is used to 
suppress the output of _qsub_ itself. This requires _inputPath_ to be set 
to _/dev/stdin_.

    qsub -b process -s --quiet -j ./numbers.json | qsub -b process -s --quiet -j ./multiply.json