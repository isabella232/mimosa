#
# Deploy source
#

set -e

if [ -z "$GOOGLE_CLOUD_PROJECT" ]; then
    echo "GOOGLE_CLOUD_PROJECT must be defined";
    exit 1
fi

if [ -z "$1" ]; then
    echo "usage: create-source.sh <name> <source-dir> <config-file> e.g. create-source.sh aws1 sources/aws config.json";
    exit 1
fi

if [ -z "$2" ]; then
    echo "usage: create-source.sh <name> <source-dir> <config-file> e.g. create-source.sh aws1 sources/aws config.json";
    exit 1
fi

if [ ! -f "$3" ]; then
    echo "config file does not exist: $3";
    exit 1
fi

SOURCE_NAME=mimosa-source-$1
CLOUD_FUNCTION_SOURCE=$2
CONFIG_FILE=$3
SALT=`xxd -l4 -ps /dev/urandom`
BUCKET="$SOURCE_NAME-$SALT"

echo "Name        : $SOURCE_NAME"
echo "Src Dir     : $CLOUD_FUNCTION_SOURCE"
echo "Config File : $CONFIG_FILE"
echo "Salt        : $SALT"
echo "Bucket      : $BUCKET"

echo
echo "Creating bucket ..."
gsutil mb -b on gs://$BUCKET

echo
echo "Copying config to bucket ..."
gsutil cp $CONFIG_FILE gs://$BUCKET/config.json

echo
echo "Creating service account ..."
gcloud iam service-accounts create $SOURCE_NAME

echo
echo "Setting service account permissions ..."
gsutil iam ch serviceAccount:$SOURCE_NAME@$GOOGLE_CLOUD_PROJECT.iam.gserviceaccount.com:objectAdmin gs://$BUCKET
echo "Permisions set."

echo
echo "Creating pub-sub topic ..."
gcloud pubsub topics create $SOURCE_NAME

echo
echo "Deploying cloud function ..."
gcloud functions deploy \
 --runtime go111 \
 --trigger-topic $SOURCE_NAME \
 --service-account=$SOURCE_NAME@$GOOGLE_CLOUD_PROJECT.iam.gserviceaccount.com \
 --set-env-vars MIMOSA_GCP_BUCKET=$BUCKET \
 --source $CLOUD_FUNCTION_SOURCE \
 --entry-point=SourceSubscriber \
 $SOURCE_NAME

echo "Test your source by sending a message to the topic:"
echo
echo "gcloud pubsub topics publish projects/$GOOGLE_CLOUD_PROJECT/topics/$SOURCE_NAME --message \"go\""

echo
echo "Finished"
