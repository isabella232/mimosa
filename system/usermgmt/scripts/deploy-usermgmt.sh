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

echo "Deploying usercreation ..."
gcloud functions deploy \
 --runtime go111 \
 --trigger-event providers/firebase.auth/eventTypes/user.create \
 --trigger-resource $MIMOSA_GCP_PROJECT \
 --source system/usermgmt \
 --entry-point UserCreated \
 system-usercreation
