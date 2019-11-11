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

if [ -z "$MIMOSA_WORKSPACE" ]; then
    echo "MIMOSA_WORKSPACE must be defined";
    exit 1
fi

UUID=`uuidgen| awk '{print tolower($0)}'`
NAME="source-$UUID"

echo "UUID      : $UUID"
echo "Workspace : $MIMOSA_WORKSPACE"

echo
echo "Creating storage bucket ..."
gsutil mb -b on gs://$NAME

echo
echo "Creating pub-sub topic ..."
gcloud pubsub topics create $NAME

echo
MIMOSA_WORKSPACE=$MIMOSA_WORKSPACE UUID=$UUID sh system/router/scripts/deploy-router.sh

echo
echo "Finished"
