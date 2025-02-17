## Setup

Installed the Google Cloud SDK tools:

    brew cask install google-cloud-sdk

Install jq

    brew install jq

Install berglas

    brew install berglas

You'll need to authenticate with gcloud:

    gcloud auth login

To keep mimosa infrastructure separate from everything else, create a new Google Cloud project:

    gcloud projects list

Set the MIMOSA_GCP_PROJECT environment variable to match your project ID (not project name):

    export MIMOSA_GCP_PROJECT=mimosa-123456

Configure gcloud to use the correct project ID (not project name):

    gcloud config set project $MIMOSA_GCP_PROJECT

Enable Firestore in your new project by going to this URL:

    https://console.cloud.google.com/firestore

When prompted select "Native Mode" then choose a region. Choose one close to you if you have no other preference.

Configure docker to work with GCP for cloud run:

    gcloud auth configure-docker

You may find that GCP asks you to enable particular APIs or to enable billing as you deploy mimosa e.g. if you have not yet enabled Cloud Functions you may see a message like this. Choose "y".

    API [cloudfunctions.googleapis.com] not enabled on project [123456]. Would you like to enable and retry (this will take a few minutes)? (y/N)?  y

Enable Identity Platform in your new project:

    https://console.cloud.google.com/customer-identity

Enable 'machine' auth:

    gcloud auth application-default login

## Test

Build, test and lint Mimosa like this:

    make

There are also per-module Makefiles.

## Deployment

Deploy Mimosa like this:

    make deploy

There are also per-module Makefiles to allow deployment of individual Cloud Functions.

## Continuous Integration with Cloud Build

Cloud Build is supported.

As a one-off setup step you need to create the "mimosabuild" container that Cloud Run requires:

    cd docker
    make mimosabuild

This builds the container locally and pushes it to your project, where Cloud Run can find it.

You need to enable triggers in Cloud Run: https://console.cloud.google.com/cloud-build/triggers

Set up triggers for push and/or pull request.

## Creating a user and workspace

In order to use Mimosa, you need to create an applcation user.

Once Mimosa is deployed as described above, you need to create the user in the Identity Platform:

https://console.cloud.google.com/customer-identity

The simplest approach is to specify an example email with password.

The "usermgmt" world builder will create a workspace for the new user in Firestore automatically.

You'll need the ID of the new workspace later and you can find it by browsing Firestore directly here:

https://console.cloud.google.com/firestore/data

All data for a workspace is stored in "/ws/<id>" where ID is a five character string like "mywm3".

## Sources

Sources must be created individually by hand.

First you need to know the ID of the workspace where you're going to deploy your source - find this as described above.

Second you need a config file that contains credentials and other configuration for your source. We recommend following the principle of least privilege. You can create a dedicated account for the target cloud provider with read-only permissions and use the credentials for that account in mimosa. Do not upload high privilege creds to mimosa at this stage!

The following example assumes AWS. You will need your AWS access key, secret key and region. Create a file called "aws.json" and put your values in there:

```
{
    "region": "eu-west-1",
    "accessKey": "xxxx",
    "secretKey": "xxxx"
}
```

Finally create the source with a command like this:

    WORKSPACE=xxxx CONFIG_FILE=xxxx.json make -C sources create

e.g.

    WORKSPACE=mywm3 CONFIG_FILE=aws.json make -C sources create

## Extensible Service Proxy (ESP)

Deploy ESP to handle auth and CORS for all API calls.

Before deployment make sure these services are enabled:

    gcloud services enable servicemanagement.googleapis.com
    gcloud services enable servicecontrol.googleapis.com
    gcloud services enable endpoints.googleapis.com

Update the 'openapi/openapi-mimosa.yaml' file to contain your project details

* Update 'host' param to your cloud run instance url
* Update 'x-google-issuer' to contain your project id
* Update 'x-google-audiences' to your project id
* Update the 'x-google-backend' address to your cloud function url

Deploy the ESP to Cloud Run:

    gcloud beta run deploy mimosa-esp \
    --image="gcr.io/endpoints-release/endpoints-runtime-serverless:1" \
    --allow-unauthenticated \

Create the endpoint service:

    gcloud endpoints services deploy openapi/openapi-mimosa.yaml

When you see this error:

    Serverless ESP expects ENDPOINTS_SERVICE_NAME in environment variables.

The endpoints are deployed, but we need to deploy the cloud run container again but this time we are specifying the ENDPOINTS_SERVICE_NAME env var for it. Replace the "xxxx" with your endpoint service name:

    gcloud beta run deploy mimosa-esp \
    --image="gcr.io/endpoints-release/endpoints-runtime-serverless:1" \
    --allow-unauthenticated \
    --set-env-vars ENDPOINTS_SERVICE_NAME=mimosa-esp-xxxx.a.run.app,ESP_ARGS=--cors_preset=basic,--cors_allow_origin=localhost

Test your endpont is authenticating calls by making an unauthenticated curl request:

    curl -X POST https://mimosa-esp-tfmdd2vwoq-uc.a.run.app/api/v1/runtask

You should see a 401 error.

Now try an authenticated call using a Firebase token (Remember to change URL to match your project):

    export FIREBASE_TOKEN=`curl 'https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=xxxx' -H 'Content-Type: application/json' --data-binary '{"email":"xxxx@example.com","password":"xxxx","returnSecureToken":true}'| jq -r .idToken`

    curl -d "{}" -X POST --header "Authorization: Bearer $FIREBASE_TOKEN" https://mimosa-esp-xxxx-uc.a.run.app/api/v1/runtask

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

    mimosa-berglas/a8e1e136ae5ea7c143a345e99aae843f22d6xxxx

If none is found, we fall back to the "default" secret:

    mimosa-berglas/default

This means you can associate a key with each host but you can also simplify development by using a single key for all hosts and uploading it as the default secret.

Upload a secret, like the pem file you got from AWS, as follows. Using "edit" helps avoid problems with your shell mangling your private key:

    EDITOR="code -w" berglas edit mimosa-berglas/default
