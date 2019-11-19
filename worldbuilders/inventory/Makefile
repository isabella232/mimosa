
include ../../build/common.mk

# Use "make deploy -j" for parallel deployment
PHONY+= deploy
deploy: deploy-aws deploy-gcp deploy-vmpooler
	@echo "✅ All deployments finished (`date '+%H:%M:%S'`)"

PHONY+= deploy-aws
deploy-aws:
	${call deploy-inventory,type-aws-instance,AWS}

PHONY+= deploy-gcp
deploy-gcp:
	${call deploy-inventory,type-gcp-instance,GCP}

PHONY+= deploy-vmpooler
deploy-vmpooler:
	${call deploy-inventory,type-vmpooler-instance,VMPooler}

define deploy-inventory
	@echo "🔘 Deploying wb-inventory-${1} ... (`date '+%H:%M:%S'`)"
	@if ! gcloud pubsub topics describe ${1} &> /dev/null; then gcloud pubsub topics create ${1}; fi
	@gcloud functions deploy \
	--runtime go111 \
	--no-allow-unauthenticated \
	--trigger-topic ${1} \
	--source . \
	--entry-point ${2} \
	wb-inventory-${1}
	@echo "✅ Deployed wb-inventory-${1} (`date '+%H:%M:%S'`)"
endef
