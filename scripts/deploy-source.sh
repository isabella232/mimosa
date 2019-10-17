#
# Deploy source
#

set -e

if [ -z "$MIMOSA_GCP_PROJECT" ]; then
    echo "MIMOSA_GCP_PROJECT must be defined";
    exit 1
fi

if [ -z "$1" ]; then
    echo "usage: deploy-source.sh <full-source-name> <source-dir> <world-builder-source-dir> <config-file> e.g. deploy-source.sh src-aws1-a24f sources/aws worldbuilders/awsfinalize config.json";
    exit 1
fi

if [ -z "$2" ]; then
    echo "usage: deploy-source.sh <full-source-name> <source-dir> <world-builder-source-dir> <config-file> e.g. deploy-source.sh src-aws1-a24f sources/aws worldbuilders/awsfinalize config.json";
    exit 1
fi

if [ ! -d "$2" ]; then
    echo "source dir does not exist: $2";
    exit 1
fi

if [ ! -d "$3" ]; then
    echo "world builder source dir does not exist: $3";
    exit 1
fi

if [ ! -f "$4" ]; then
    echo "config file does not exist: $4";
    exit 1
fi

NAME=$1
CLOUD_FUNCTION_SOURCE=$2
WORLD_BUILDER_CLOUD_FUNCTION_SOURCE=$3
CONFIG_FILE=$4

echo "Name                   : $NAME"
echo "Code Dir               : $CLOUD_FUNCTION_SOURCE"
echo "World Builder Code Dir : $WORLD_BUILDER_CLOUD_FUNCTION_SOURCE"
echo "Config File            : $CONFIG_FILE"

echo
echo "Copying config to bucket ..."
gsutil cp $CONFIG_FILE gs://$NAME/config.json

echo
echo "Deploying source cloud function ..."
gcloud functions deploy \
 --runtime go111 \
 --trigger-topic $NAME \
 --service-account=$NAME@$MIMOSA_GCP_PROJECT.iam.gserviceaccount.com \
 --set-env-vars MIMOSA_GCP_BUCKET=$NAME \
 --source $CLOUD_FUNCTION_SOURCE \
 --entry-point=SourceSubscriber \
 $NAME

echo
echo "Deploying world-builder cloud function ..."
gcloud functions deploy \
 --runtime go111 \
 --trigger-resource $NAME \
 --trigger-event google.storage.object.finalize \
 --source $WORLD_BUILDER_CLOUD_FUNCTION_SOURCE \
 --entry-point HandleInstance \
 WorldBuilder-$NAME

echo
echo "Finished"
