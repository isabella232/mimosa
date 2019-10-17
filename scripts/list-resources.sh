#
# List source resources
#

echo
echo Buckets
echo
gsutil ls |grep gs://src-

echo
echo Service Accounts
echo
gcloud iam service-accounts list|grep src-

echo
echo Pub-Sub Topics
echo
gcloud pubsub topics list|grep src-

echo
echo Functions
echo
gcloud functions list|grep src-

