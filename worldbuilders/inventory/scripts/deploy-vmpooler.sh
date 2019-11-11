set -e

# Check for the presence of the .git dir to determine if we're in the root of the repo
if [ ! -d ".git" ]; then
    echo "script must be run in the root of the mimosa repo";
    exit 1
fi

echo "Deploying inventory for aws-instance..."
gcloud functions deploy \
 --runtime go111 \
 --trigger-topic vmpooler-instance \
 --source worldbuilders/inventory \
 --entry-point VMPooler \
 wb-inventory-vmpooler-instance
