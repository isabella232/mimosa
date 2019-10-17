#
# Trigger source
#

set -e

if [ -z "$MIMOSA_GCP_PROJECT" ]; then
    echo "MIMOSA_GCP_PROJECT must be defined";
    exit 1
fi

if [ -z "$1" ]; then
    echo "usage: trigger-source.sh <full-source-name> e.g. trigger-source.sh src-aws1-a24f";
    exit 1
fi

echo
echo "Triggering source ..."
gcloud pubsub topics publish projects/$MIMOSA_GCP_PROJECT/topics/$1 --message "go"
