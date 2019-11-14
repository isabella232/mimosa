
include ../../build/common.mk

check-env:
# No special environment needed

# Use "make deploy -j" for parallel deployment
PHONY+= deploy
deploy: deploy-aws deploy-gcp deploy-vmpooler
	@echo "✅ All deployments finished (`date '+%H:%M:%S'`)"

PHONY+= deploy-aws
deploy-aws: check-env
	@echo "🔘 Deploying wb-inventory-aws-instance ... (`date '+%H:%M:%S'`)"
	@if ! gcloud pubsub topics describe aws-instance > /dev/null; then gcloud pubsub topics create aws-instance; fi
	@gcloud functions deploy \
	--runtime go111 \
	--trigger-topic aws-instance \
	--source . \
	--entry-point AWS \
	wb-inventory-aws-instance
	@echo "✅ Deployed wb-inventory-aws-instance (`date '+%H:%M:%S'`)"

PHONY+= deploy-gcp
deploy-gcp: check-env
	@echo "🔘 Deploying wb-inventory-gcp-instance ... (`date '+%H:%M:%S'`)"
	@if ! gcloud pubsub topics describe gcp-instance > /dev/null; then gcloud pubsub topics create gcp-instance; fi
	@gcloud functions deploy \
	--runtime go111 \
	--trigger-topic gcp-instance \
	--source . \
	--entry-point GCP \
	wb-inventory-gcp-instance
	@echo "✅ Deployed wb-inventory-gcp-instance (`date '+%H:%M:%S'`)"

PHONY+= deploy-vmpooler
deploy-vmpooler: check-env
	@echo "🔘 Deploying wb-inventory-vmpooler-instance ... (`date '+%H:%M:%S'`)"
	@if ! gcloud pubsub topics describe vmpooler-instance > /dev/null; then gcloud pubsub topics create vmpooler-instance; fi
	@gcloud functions deploy \
	--runtime go111 \
	--trigger-topic vmpooler-instance \
	--source . \
	--entry-point VMPooler \
	wb-inventory-vmpooler-instance
	@echo "✅ Deployed wb-inventory-vmpooler-instance (`date '+%H:%M:%S'`)"
