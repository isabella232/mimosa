set -e

# Check for the presence of the .git dir to determine if we're in the root of the repo
if [ ! -d ".git" ]; then
    echo "script must be run in the root of the mimosa repo";
    exit 1
fi

echo "Deploying api-v1-runtask ..."
gcloud functions deploy \
 --runtime go111 \
 --no-allow-unauthenticated \
 --trigger-http \
 --source api \
 --entry-point RunTask \
 api-v1-runtask
