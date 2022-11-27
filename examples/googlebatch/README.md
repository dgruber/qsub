# Example of Submitting Google Batch Jobs with qsub

Log into Google Cloud account:

    gcloud auth application-default login

Point GOOGLE_PROJECT to your Google project ID:

    export GOOGLE_PROJECT=<YOUR_GOOGLE_PROJECT_ID>

Choose a GOOGLE_REGION:

    export GOOGLE_REGION=

Submit a Blast Biocontainer as DRMAA2 JSON file:

Adapt the _blast.json_ file to point to a Google Cloud bucket you own (gs://).

    qsub -b googlebatch -j ./blast.json
    projects/<YOUR_GOOGLE_PROJECT_ID>/locations/us-central1/jobs/drmaa2-166957536-9685

The returned string is the Google Batch job ID. Check Google Console for the
newly created batch job.
