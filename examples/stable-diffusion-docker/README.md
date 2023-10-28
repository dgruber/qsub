# Running Stable Diffusion as DRMAA2 JSON Job Template

This example runs Stable Diffusion in a container without GPUs on an Intel
Mac. What is required is a _token.txt_ file in this directory which contains
a valid Huggingface token. You need to create an account and generate a new
access token in order to let the container download the official model files.

You need to change the absolute path of the _token.txt_ inside the 
_stable-diffusion.json_ file. The output image appears in the directory
_/tmp/qsub/output_ on the host.

It uses the _ghcr.io/fboulnois/stable-diffusion-docker_ image and mounts
local folders input/ouput/cache/tmp inside the container. For more details
please check https://github.com/fboulnois/stable-diffusion-docker/tree/main

The image creation time (and quality) highly depends on the resolution and 
iterations and of course on the absense of a GPU.

   docker pull ghcr.io/fboulnois/stable-diffusion-docker:1.41.0
   qsub -b docker -s -j stable-diffusion.json
