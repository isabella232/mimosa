# mimosa

Experiments with Google Cloud

## Setup

Installed the Google Cloud SDK tools:

    brew cask install google-cloud-sdk

You'll need to authenticate with gcloud:

    gcloud auth login

To keep mimosa infrastructure separate from everything else, create a new Google Cloud project:

    gcloud projects list

Set the MIMOSA_GCP_PROJECT environment variable to match your project ID (not project name):

    export MIMOSA_GCP_PROJECT=mimosa-255913

Configure gcloud to use the correct project ID (not project name):

    gcloud config set project mimosa-255913

Enable Firestore in Native Mode in your new project: [TODO Which region should we choose? For now, anything is fine]:

    https://console.cloud.google.com/firestore

You may find that GCP asks you to enable particular APIs or to enable billing as you deploy mimosa.

## User Interface

Deploy the UI:

    gcloud functions deploy --entry-point HandleHTTPRequest --runtime go111  --trigger-http --source ui ui

## Cloud Run

To deploy the runner container to Cloud Run you need to perform a one-time setup step to configure docker to authenticate with GCP:

    gcloud auth configure-docker

To build the container (substitute your own GCP project ID):

    cd docker
    docker build . -t gcr.io/PROJECT_ID/runner

To run locally:

    docker run -a STDOUT -a STDERR -it --env PORT=8080 -p 8080:8080 gcr.io/PROJECT_ID/runner

To test, create a customized payload that contains your AWS instance details and PEM file contents and run:

    curl localhost:8080 --data-binary "@payload.json"

Push to GCR like this:

    docker push gcr.io/mimosa-256008/runner

Deploy in Cloud Run (currently by hand, gcloud commands to follow). Be sure to enable authentication.

Test using this handy Google supplied alias:

    alias gcurl='curl --header "Authorization: Bearer $(gcloud auth print-identity-token)"'\n

Now test the deployment:

    gcurl  https://runner-xxxx-ew.a.run.app --data-binary "@payload.json"

## Sources

### Source credentials

You will need credentials and other configuration for your source. We recommend following the principle of least privilege. You can create a dedicated account for the target cloud provider with read-only permissions and use the credentials for that account in mimosa. Do not upload high privilege creds to mimosa at this stage! [TODO We need additional docs here explaining how to achieve this.]

The following example assumes AWS. You will need your AWS access key, secret key and region. Create a file called "config.json" and put the values in there:

```
{
    "region": "eu-west-1",
    "accessKey": "AKIAIOSFODNN7EXAMPLE",
    "secretKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}
```

### Deploying a source

Run `scripts/create-source.sh` specifying the name, source dir and config file for your source e.g.

    sh scripts/create-source.sh aws1 sources/aws awsconfig.json

If you have not yet enabled Cloud Functions you may see a message like this. Choose "y".

    API [cloudfunctions.googleapis.com] not enabled on project
    [870066425029]. Would you like to enable and retry (this will take a few minutes)? (y/N)?  y
