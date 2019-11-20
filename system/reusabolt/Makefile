
include ../../build/common.mk

SERVICE_URL:=${shell gcloud beta run services list --platform managed | grep runner | grep -o -E "https[^ ]+"}
GENERATED_SECRETS_BUCKET:=mimosa-secrets-${shell uuidgen | awk '{print tolower($0)}'}

check-env:
ifndef SERVICE_URL
	$(error SERVICE_URL must be defined)
endif

PHONY+= ensure-secrets-bucket
ensure-secrets-bucket:
ifndef MIMOSA_GCP_PROJECT
	$(error MIMOSA_GCP_PROJECT must be defined)
endif
	@if ! gsutil ls | grep mimosa-secrets &> /dev/null; then \
		berglas bootstrap --project ${MIMOSA_GCP_PROJECT} --bucket ${GENERATED_SECRETS_BUCKET}; \
	fi

PHONY+= deploy
deploy: ensure-secrets-bucket check-env
	@echo "🔘 Deploying system-reusabolt ... (`date '+%H:%M:%S'`)"
	@gcloud functions deploy \
		--runtime go111 \
		--no-allow-unauthenticated \
		--trigger-topic reusabolt \
		--set-env-vars MIMOSA_SERVICE_URL=${SERVICE_URL},MIMOSA_SECRETS_BUCKET=${shell gsutil ls | grep mimosa-secrets | grep -o -E "gs:\/\/[^\/]+"} \
		--source . \
		--entry-point TriggerReusabolt \
		system-reusabolt
	@echo "✅ Deployed system-reusabolt (`date '+%H:%M:%S'`)"
