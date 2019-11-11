#
# Delete source
#

if [ -z "$UUID" ]; then
    echo "UUID must be defined";
    exit 1
fi

NAME=source-$1

echo
echo "Deleting bucket ..."
gsutil -m rm -r gs://$NAME

echo
echo "Deleting pub-sub topic ..."
gcloud pubsub topics delete $NAME

echo
echo "Deleting cloud functions ..."
gcloud functions delete --quiet $NAME
gcloud functions delete --quiet system-router-$UUID

echo
echo "Finished"
