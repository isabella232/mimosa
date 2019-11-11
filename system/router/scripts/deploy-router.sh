set -e

# Check for the presence of the .git dir to determine if we're in the root of the repo
if [ ! -d ".git" ]; then
    echo "script must be run in the root of the mimosa repo";
    exit 1
fi

if [ -z "$MIMOSA_WORKSPACE" ]; then
    echo "MIMOSA_WORKSPACE must be defined";
    exit 1
fi

if [ -z "$UUID" ]; then
    echo "UUID must be defined";
    exit 1
fi

echo "Deploying router for $MIMOSA_WORKSPACE - $UUID ..."
gcloud functions deploy \
 --runtime go111 \
 --trigger-resource source-$UUID \
 --trigger-event google.storage.object.finalize \
 --set-env-vars MIMOSA_WORKSPACE=$MIMOSA_WORKSPACE, \
 --source system/router \
 --entry-point Route \
 system-router-$UUID
