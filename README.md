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

    gcloud config set project $MIMOSA_GCP_PROJECT

Enable Firestore in Native Mode in your new project: [TODO Which region should we choose? For now, anything is fine]:

    https://console.cloud.google.com/firestore

Configure docker to work with GCP for cloud run:

    gcloud auth configure-docker

You may find that GCP asks you to enable particular APIs or to enable billing as you deploy mimosa e.g. if you have not yet enabled Cloud Functions you may see a message like this. Choose "y".

    API [cloudfunctions.googleapis.com] not enabled on project
    [870066425029]. Would you like to enable and retry (this will take a few minutes)? (y/N)?  y

## Test

Build, test and lint Mimosa like this:

    make

There are also per-module Makefiles.

## Deployment

Deploy Mimosa like this:

    make deploy

There are also per-module Makefiles to allow deployment of individual Cloud Functions.

### Sources 

Sources must be created individually:

    WORKSPACE=xxxxx CONFIG_FILE=xxx.json make -C sources create 

The workspace id can be obtained from firestore.

The config file contains credentials and other configuration for your source. We recommend following the principle of least privilege. You can create a dedicated account for the target cloud provider with read-only permissions and use the credentials for that account in mimosa. Do not upload high privilege creds to mimosa at this stage! [TODO We need additional docs here explaining how to achieve this.]

The following example assumes AWS. You will need your AWS access key, secret key and region. Create a file called "config.json" and put the values in there:

```
{
    "region": "eu-west-1",
    "accessKey": "AKIAIOSFODNN7EXAMPLE",
    "secretKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}
```

## Extensible Service Proxy (ESP)

Deploy ESP to handle auth and CORS for all API calls.

Before deployment make sure these services are enabled:

    gcloud services enable servicemanagement.googleapis.com\t
    gcloud services enable servicecontrol.googleapis.com\t
    gcloud services enable endpoints.googleapis.com\t

Deploy the ESP to Cloud Run:

    gcloud beta run deploy mimosa-esp \
    --image="gcr.io/endpoints-release/endpoints-runtime-serverless:1" \
    --allow-unauthenticated \

Create the endpoint service:

    gcloud endpoints services deploy openapi/openapi-mimosa.yaml

If you see this error:

    Serverless ESP expects ENDPOINTS_SERVICE_NAME in environment variables.

Then deploy again but this time specifying the ENDPOINTS_SERVICE_NAME env var:

    gcloud beta run deploy mimosa-esp \
    --image="gcr.io/endpoints-release/endpoints-runtime-serverless:1" \
    --allow-unauthenticated \
    --set-env-vars ENDPOINTS_SERVICE_NAME=mimosa-esp-tfmdd2vwoq-uc.a.run.app,ESP_ARGS=--cors_preset=basic,--cors_allow_origin=localhost

Test your endpont is authenticating calls by making an unauthenticated curl request:

    curl https://mimosa-esp-tfmdd2vwoq-uc.a.run.app/hello

You should see a 401 error.

Now try an authenticated call using a Firebase token:

    export FIREBASE_TOKEN=`curl 'https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=AIzaSyCQieKOS6B36ut_o5n0loeW8rXetEqXnb0' -H 'Content-Type: application/json' --data-binary '{"email":"xxxx@example.com","password":"xxxx","returnSecureToken":true}'| jq -r .idToken`

    curl --header "Authorization: Bearer $FIREBASE_TOKEN" https://mimosa-esp-tfmdd2vwoq-uc.a.run.app/hello

## Secrets

We use [berglas](https://github.com/GoogleCloudPlatform/berglas) for secrets which should be installed locally as described in their README to allow secrets to be added to Mimosa.

Berglas is set up as part of mimosa deployment.

The service account running the "system-reusabolt" cloud function (likely "App Engine default service account") needs some additional permissions:

* Storage Viewer on the berglas bucket: `gsutil acl ch -u <service-account-email>:R gs://<bucket-name>`
* Cloud KMS CryptoKey Decryptor on the berglas key (https://console.cloud.google.com/security/kms)

Berlgas secrets are referenced like this:

    <secrets-bucket>/<secret-name>

e.g.

    mimosa-berglas/foo

When running bolt, we check for a secret named like this:

    <secrets-bucket>/<host-firestore-id>

e.g.

    mimosa-berglas/a8e1e136ae5ea7c143a345e99aae843f22d6e5b1

If none is found, we fall back to the "default" secret:

    mimosa-berglas/default

This means you can associate a key with each host but you can also simplify development by using a single key for all hosts and uploading it as the default secret.

Upload a secret, like the pem file you got from AWS, as follows. Using "edit" helps avoid problems with your shell mangling your private key:

    EDITOR="code -w" berglas edit mimosa-berglas/default
