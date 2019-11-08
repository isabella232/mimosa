#
# Deploy source
#

set -e

# Check for the presence of the .git dir to determine if we're in the root of the repo
if [ ! -d ".git" ]; then
    echo "script must be run in the root of the mimosa repo";
    exit 1
fi

if [ -z "$MIMOSA_GCP_PROJECT" ]; then
    echo "MIMOSA_GCP_PROJECT must be defined";
    exit 1
fi

if [ -z "$1" ]; then
    echo "usage: deploy-source.sh <uuid> <source-dir> <config-file> e.g. deploy-source.sh -11e15a54-72ca-44a0-9cd6-29d2a385d46b	sources/aws config.json";
    exit 1
fi

if [ ! -d "$2" ]; then
    echo "source dir does not exist: $2";
    exit 1
fi

if [ ! -f "$3" ]; then
    echo "config file does not exist: $3";
    exit 1
fi

NAME=source-$1
CLOUD_FUNCTION_SOURCE=$2
CONFIG_FILE=$3

echo "Name        : $NAME"
echo "Code Dir    : $CLOUD_FUNCTION_SOURCE"
echo "Config File : $CONFIG_FILE"

echo
echo "Copying config to bucket ..."
gsutil cp $CONFIG_FILE gs://$NAME/config.json

echo
echo "Deploying source cloud function ..."
gcloud functions deploy \
 --runtime go111 \
 --trigger-topic $NAME \
 --set-env-vars MIMOSA_GCP_BUCKET=$NAME, \
 --source $CLOUD_FUNCTION_SOURCE \
 --entry-point=HandleMessage \
 $NAME

echo
echo "Finished"
