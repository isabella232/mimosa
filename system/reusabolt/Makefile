
include ../../build/common.mk

check-env:
ifndef SERVICE_URL
	$(error SERVICE_URL and SECRETS_BUCKET must be defined)
endif
ifndef SECRETS_BUCKET
	$(error SERVICE_URL and SECRETS_BUCKET must be defined)
endif

PHONY+= deploy
deploy: check-env
	@echo "🔘 Deploying system-reusabolt ... (`date '+%H:%M:%S'`)"
	@gcloud functions deploy \
		--runtime go111 \
		--trigger-topic reusabolt \
		--set-env-vars MIMOSA_SERVICE_URL=${SERVICE_URL},MIMOSA_SECRETS_BUCKET=${SECRETS_BUCKET} \
		--source . \
		--entry-point TriggerReusabolt \
		system-reusabolt
	@echo "✅ Deployed system-reusabolt (`date '+%H:%M:%S'`)"
