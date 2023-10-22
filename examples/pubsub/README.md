# Sending DRMAA2 JobTemplates as CloudEvents to PubSub

With qsub you can also send a JobTemplate to a Google PubSub topic.
It will be wrapped as a CloudEvent. A backend can receive the 
JobTemplate and execute the job accordingly.

There are following requirements for submitting:

- googleProjectID extension needs to be set to the google project
- queueName must be set to the PubSub topic you want to send to
- your shell needs to be logged into your Google Cloud account and
  your account needs to have the priviledges to publish a message
  to the topic ('gcloud auth application-default login').

```
    qsub -b pubsub ./drmaa2jobtemplatescloudevent.json
```