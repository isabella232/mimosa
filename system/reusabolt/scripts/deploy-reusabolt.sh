set -e

# Check for the presence of the .git dir to determine if we're in the root of the repo
if [ ! -d ".git" ]; then
    echo "script must be run in the root of the mimosa repo";
    exit 1
fi

if [ -z "$MIMOSA_SECRETS_BUCKET" ]; then
    echo "MIMOSA_SECRETS_BUCKET must be defined";
    exit 1
fi

if [ -z "$MIMOSA_SERVICE_URL" ]; then
    echo "MIMOSA_SERVICE_URL must be defined";
    exit 1
fi

echo "Deploying reusabolt cloud function ..."
gcloud functions deploy \
 --runtime go111 \
 --trigger-topic reusabolt \
 --set-env-vars MIMOSA_SERVICE_URL=$MIMOSA_SERVICE_URL,MIMOSA_SECRETS_BUCKET=$MIMOSA_SECRETS_BUCKET \
 --source system/reusabolt \
 --entry-point TriggerReusabolt \
 system-reusabolt
