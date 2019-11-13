
include ../../build/common.mk

check-env:
ifndef WORKSPACE
	$(error WORKSPACE must be defined)
endif
ifndef UUID
	$(error UUID must be defined)
endif

PHONY+= deploy
deploy: check-env
	@echo "🔘 Deploying system-router-${WORKSPACE}-${UUID} ... (`date '+%H:%M:%S'`)"
	@gcloud functions deploy \
		--runtime go111 \
		--trigger-resource source-${UUID} \
		--trigger-event google.storage.object.finalize \
		--set-env-vars WORKSPACE=${WORKSPACE}, \
		--source . \
		--entry-point Route \
		system-router-${WORKSPACE}-${UUID}
	@echo "✅ Deployed system-router-${WORKSPACE}-${UUID} (`date '+%H:%M:%S'`)"
