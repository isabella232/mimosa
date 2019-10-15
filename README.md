# mimosa

Experiments with Google Cloud

## Setup

Installed the Google Cloud SDK tools:

    brew cask install google-cloud-sdk

You'll need to authenticate with gcloud:

    gcloud auth login

To keep mimosa infrastructure separate from everything else, create a new Google Cloud project here: https://console.cloud.google.com/projectcreate e.g. "mimosa-255913"

Set your GOOGLE_CLOUD_PROJECT environment variable to match your project name:

    export GOOGLE_CLOUD_PROJECT=mimosa-255913

Configure gcloud to use the correct project:

    gcloud config set project mimosa-255913

## Source credentials

You will need credentials and other configuration for your source. We recommend following the principle of least privilege. You can create a dedicated account for the target cloud provider with read-only permissions and use the credentials for that account in mimosa. Do not upload high privilege creds to mimosa at this stage!


The following example assumes AWS. You will need your AWS access key, secret key and region. Create a file called "config.json" and put the values in there:

```
{
    "region": "eu-west-1",
    "accessKey": "AKIAIOSFODNN7EXAMPLE",
    "secretKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}
```

## Deploying a source

Deploying a source is currently a manual process.

1) Choose a salt value of 8 random characters e.g. 7870d1da. If the command fails, any random means of choosing the value is fine.

    xxd -l4 -ps /dev/urandom

2) Choose a source name e.g. "aws1"

3) Create a bucket:

    gsutil mb -b on gs://mimosa-source-SOURCE_NAME-SALT

e.g.

    gsutil mb -b on gs://mimosa-source-aws1-7870d1da

The salt is needed at this step to ensure the bucket name is globally unique.

4) Copy the "config.json" source configuration file that you created above into the bucket:

    gsutil cp config.json gs://mimosa-source-SOURCE_NAME-SALT

e.g.

    gsutil cp config.json gs://mimosa-source-aws1-7870d1da

5) Create a service account:

    gcloud iam service-accounts create mimosa-source-SOURCE_NAME --display-name "mimosa source: SOURCE_NAME"

e.g.

    gcloud iam service-accounts create mimosa-source-aws1 --display-name "mimosa source: aws1"

6) Optionally create a key file for the service account if you need to run the source outside of GCP. This isn't usually necessary.

    gcloud iam service-accounts keys create --iam-account=SOURCE_NAME@PROJECT_NAME.iam.gserviceaccount.com SOURCE_NAME.json

e.g.

    gcloud iam service-accounts keys create --iam-account=mimosa-source-aws1@mimosa-255913.iam.gserviceaccount.com mimosa-source-aws1.json

7) Give the service account permission to access to the bucket:

The command looks like this and writes the key into a file called "SOURCE_NAME.json":

    gsutil iam ch serviceAccount:SOURCE_NAME@PROJECT_NAME.iam.gserviceaccount.com:objectAdmin gs://BUCKET_NAME

e.g.

    gsutil iam ch serviceAccount:mimosa-source-aws1@mimosa-255913.iam.gserviceaccount.com:objectAdmin gs://mimosa-source-aws1-7870d1da

8) Create a pub-sub topic:

    gcloud pubsub topics create mimosa-source-SOURCE_NAME

e.g.

    gcloud pubsub topics create mimosa-source-aws1

9) Deploy the cloud function (don't substitute anything for MIMOSA_GCP_BUCKET, just use it verbatim!):

gcloud functions deploy --runtime go111 --trigger-topic mimosa-source-SOURCE_NAME --service-account=mimosa-source-SOURCE_NAME@PROJECT_NAME.iam.gserviceaccount.com --set-env-vars MIMOSA_GCP_BUCKET=mimosa-source-SOURCE_NAME-SALT --source sources/aws/ SourceSubscriber

e.g.

gcloud functions deploy --runtime go111 --trigger-topic mimosa-source-aws1 --service-account=mimosa-source-aws1@mimosa-255913.iam.gserviceaccount.com --set-env-vars MIMOSA_GCP_BUCKET=mimosa-source-aws1-7870d1da --source sources/aws/ SourceSubscriber


10) Configure a scheduler job so that the source is run periodically:

    gcloud scheduler jobs create FIXME

e.g.

    gcloud scheduler jobs create FIXME

11) Test your source by posting a message to the topic:

    gcloud pubsub topics publish projects/PROJECT_NAME/topics/mimosa-source-SOURCE_NAME --message "go"

e.g.

    gcloud pubsub topics publish projects/mimosa-255913/topics/mimosa-source-aws1 --message "go"
