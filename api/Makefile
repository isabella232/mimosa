
include ../build/common.mk

check-env:
# No special environment needed

# Use "make deploy -j" for parallel deployment
PHONY+= deploy
deploy: deploy-api-v1-runtask
	@echo "✅ All deployments finished (`date '+%H:%M:%S'`)"

PHONY+= deploy-api-v1-runtask
deploy-api-v1-runtask: check-env
	@echo "🔘 Deploying api-v1-runtask ... (`date '+%H:%M:%S'`)"
	@gcloud functions deploy \
		--runtime go111 \
		--no-allow-unauthenticated \
		--trigger-http \
		--source . \
		--entry-point RunTask \
		api-v1-runtask
	@echo "✅ Deployed api-v1-runtask (`date '+%H:%M:%S'`)"
