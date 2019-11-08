set -e

# Check for the presence of the .git dir to determine if we're in the root of the repo
if [ ! -d ".git" ]; then
    echo "script must be run in the root of the mimosa repo";
    exit 1
fi

echo
echo "Creating world builder pub-sub topics ..."
if ! gcloud pubsub topics describe aws-instance ; then
    gcloud pubsub topics create aws-instance
fi
if ! gcloud pubsub topics describe gcp-instance ; then
    gcloud pubsub topics create gcp-instance
fi
if ! gcloud pubsub topics describe vmpooler-instance ; then
    gcloud pubsub topics create vmpooler-instance
fi
echo

sh worldbuilders/inventory/scripts/deploy-aws.sh
sh worldbuilders/inventory/scripts/deploy-gcp.sh
sh worldbuilders/inventory/scripts/deploy-vmpooler.sh

echo
echo "Finished"
