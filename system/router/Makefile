
include ../../build/common.mk

PHONY+= deploy
deploy:
	@echo "🔘 Deploying system-router ... (`date '+%H:%M:%S'`)"
	@gcloud functions deploy \
		--runtime go111 \
		--no-allow-unauthenticated \
		--trigger-topic system-router \
		--source . \
		--entry-point Route \
		system-router
	@echo "✅ Deployed system-router (`date '+%H:%M:%S'`)"
