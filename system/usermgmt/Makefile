
include ../../build/common.mk

check-env:
ifndef MIMOSA_GCP_PROJECT
	$(error MIMOSA_GCP_PROJECT must be defined)
endif

PHONY+= deploy
deploy: check-env
	@echo "🔘 Deploying system-usercreation ... (`date '+%H:%M:%S'`)"
	@gcloud functions deploy \
		--runtime go111 \
		--no-allow-unauthenticated \
		--trigger-event providers/firebase.auth/eventTypes/user.create \
		--trigger-resource ${MIMOSA_GCP_PROJECT} \
		--source . \
		--entry-point UserCreated \
		system-usercreation
	@echo "✅ Deployed system-usercreation (`date '+%H:%M:%S'`)"
