#!/bin/bash

echo "Starting server..."
qsub --serve localhost:8888 &

echo "Submitting jobs to server..."
qsub --server localhost:8888 --backend server -j ./sleep.json

echo "Get job status..."
qstat --server localhost:8888 1
