
include ../../build/common.mk

# Use "make deploy -j" for parallel deployment
PHONY+= deploy
deploy: deploy-qualys
	@echo "✅ All deployments finished (`date '+%H:%M:%S'`)"

PHONY+= deploy-qualys
deploy-aws:
	${call deploy-vulns,type-qualys-asset,HandleMessage}

define deploy-vulns
	@echo "🔘 Deploying wb-vulns-${1} ... (`date '+%H:%M:%S'`)"
	@if ! gcloud pubsub topics describe ${1} &> /dev/null; then gcloud pubsub topics create ${1}; fi
	@gcloud functions deploy \
	--runtime go111 \
	--no-allow-unauthenticated \
	--trigger-topic ${1} \
	--source . \
	--entry-point ${2} \
	wb-vulns-${1}
	@echo "✅ Deployed wb-vulns-${1} (`date '+%H:%M:%S'`)"
endef
