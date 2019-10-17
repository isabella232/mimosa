#
# Create source
#

set -e

if [ -z "$MIMOSA_GCP_PROJECT" ]; then
    echo "MIMOSA_GCP_PROJECT must be defined";
    exit 1
fi

if [ -z "$1" ]; then
    echo "usage: create-source.sh <name> e.g. create-source.sh aws1";
    exit 1
fi

# Name cannot be more than 30 chars to be compatible with service account name requirements
SALT=`xxd -l2 -ps /dev/urandom`
NAME="src-$1-$SALT"

echo "Name        : $NAME"

echo
echo "Creating bucket ..."
gsutil mb -b on gs://$NAME

echo
echo "Creating service account ..."
gcloud iam service-accounts create $NAME --display-name "Source: $1"

echo
echo "Setting service account permissions ..."
gsutil iam ch serviceAccount:$NAME@$MIMOSA_GCP_PROJECT.iam.gserviceaccount.com:objectAdmin gs://$NAME
echo "Permisions set."

echo
echo "Creating pub-sub topic ..."
gcloud pubsub topics create $NAME

echo
echo "Finished"
