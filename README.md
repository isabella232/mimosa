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

Deploy in Cloud Run:

    gcloud beta run deploy --image gcr.io/mimosa-256008/runner --platform managed --region europe-west1 --no-allow-unauthenticated runner

Test using this handy Google supplied alias:

    alias gcurl='curl --header "Authorization: Bearer $(gcloud auth print-identity-token)"'\n

Now test the deployment:

    gcurl  https://runner-xxxx-ew.a.run.app --data-binary "@payload.json"

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

Once the berglas bootstrap process has completed there'll be a new KMS key and storage bucket.

The service account running the "system-reusabolt" cloud function (likely "Compute Engine default service account") needs some permissions:

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
