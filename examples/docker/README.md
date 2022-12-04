# Example of Submitting Docker Jobs with qsub

The job script is executed within the local mount inside the
container. Hence the local "stageInFiles" needs to be adapted
to your absolute location of this directory. Docker does not
allow relative paths and env variable placeholders.

After adapting the path:

    qsub -b docker -j ./blast.json

    docker ps
    docker logs ..

Finally the results should appear in the local directory.

