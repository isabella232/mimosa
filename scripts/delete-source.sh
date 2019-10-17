#
# Delete source
#

if [ -z "$MIMOSA_GCP_PROJECT" ]; then
    echo "MIMOSA_GCP_PROJECT must be defined";
    exit 1
fi

if [ -z "$1" ]; then
    echo "usage: create-source.sh <name> e.g. create-source.sh aws1";
    exit 1
fi

NAME=$1

echo
echo "Deleting bucket ..."
gsutil -m rm -r gs://$NAME

echo
echo "Deleting service account ..."
gcloud iam service-accounts delete --quiet $NAME@$MIMOSA_GCP_PROJECT.iam.gserviceaccount.com

echo
echo "Deleting pub-sub topic ..."
gcloud pubsub topics delete $NAME

echo
echo "Deleting cloud functions ..."
gcloud functions delete --quiet $NAME
gcloud functions delete --quiet WorldBuilder-$NAME

echo
echo "Finished"
