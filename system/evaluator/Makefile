
include ../../build/common.mk

PHONY+= deploy
deploy:
	@echo "🔘 Deploying system-evaluator ... (`date '+%H:%M:%S'`)"
	@if ! gcloud pubsub topics describe system-evaluator &> /dev/null; then gcloud pubsub topics create system-evaluator; fi
	@gcloud functions deploy \
		--runtime go111 \
		--no-allow-unauthenticated \
		--trigger-topic system-evaluator \
		--source . \
		--entry-point Evaluate \
		system-evaluator
	@echo "✅ Deployed system-evaluator (`date '+%H:%M:%S'`)"
