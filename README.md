# mimosa

Experiments with Google Cloud

## Setup

Installed the Google Cloud SDK tools:

    brew cask install google-cloud-sdk

You'll need to authenticate with gcloud:

    gcloud auth login

To keep mimosa infrastructure separate from everything else, create a new Google Cloud project here: https://console.cloud.google.com/projectcreate e.g. "mimosa-255913"

Set the MIMOSA_GCP_PROJECT environment variable to match your project name:

    export MIMOSA_GCP_PROJECT=mimosa-255913

Configure gcloud to use the correct project:

    gcloud config set project mimosa-255913

Enable Firestore in your new project here: https://console.cloud.google.com/firestore

## User Interface

Deploy the UI:

    gcloud functions deploy HandleHTTPRequest --runtime go111  --trigger-http --source ui

## Sources

### Source credentials

You will need credentials and other configuration for your source. We recommend following the principle of least privilege. You can create a dedicated account for the target cloud provider with read-only permissions and use the credentials for that account in mimosa. Do not upload high privilege creds to mimosa at this stage!


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

    sh create-source.sh aws1 sources/aws awsconfig.json
