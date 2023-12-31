# Running Stable Diffusion with DRMAA2 JSON Job Template

In this example, we will run Stable Diffusion using a DRMAA2 JSON Job Template. The goal is to execute Stable Diffusion in a container without GPUs on an Intel Mac.

To run this example, you will need a Hugging Face token, which can be generated in the Hugging Face portal. Set this token as the environment variable _HF_TOKEN_ in the _stable-diffusion.json_ file.

This example uses the _ghcr.io/fboulnois/stable-diffusion-docker_ image and mounts local folders input, output, cache, and tmp inside the container. For more details, please refer to the repository https://github.com/fboulnois/stable-diffusion-docker/tree/main.

Before running the job, make sure to pull the image first. The time taken to create the image (and its quality) depends on the resolution and iterations.

```bash
docker pull ghcr.io/fboulnois/stable-diffusion-docker:1.41.0
qsub -b docker -s -j stable-diffusion.json
```
