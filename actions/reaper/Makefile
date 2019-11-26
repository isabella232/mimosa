
include ../../build/common.mk

PHONY+= deploy
deploy:
	@echo "🔘 Deploying actions-reaper ... (`date '+%H:%M:%S'`)"
	@if ! gcloud pubsub topics describe actions-reaper &> /dev/null; then gcloud pubsub topics create actions-reaper; fi
	@gcloud functions deploy \
		--runtime go111 \
		--no-allow-unauthenticated \
		--trigger-topic actions-reaper \
		--source . \
		--entry-point Reap \
		system-reaper
	@echo "✅ Deployed actions-reaper (`date '+%H:%M:%S'`)"
