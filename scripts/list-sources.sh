#
# List sources
#

set -e

gsutil ls |grep gs://src-|cut -d/ -f3
