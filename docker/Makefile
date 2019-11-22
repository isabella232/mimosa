



all:
	@echo "✅ Nothing to be done for docker (`date '+%H:%M:%S'`)"

PHONY+= test
test:
	@echo "✅ Nothing to be done for docker (`date '+%H:%M:%S'`)"

PHONY+= tidy
tidy:
	@echo "✅ Nothing to be done for docker (`date '+%H:%M:%S'`)"

PHONY+= format
format:
	@echo "✅ Nothing to be done for docker (`date '+%H:%M:%S'`)"

PHONY+= lint
lint:
	@echo "✅ Nothing to be done for docker (`date '+%H:%M:%S'`)"

PHONY+= deploy
deploy:
ifndef MIMOSA_GCP_PROJECT
	$(error MIMOSA_GCP_PROJECT must be defined)
endif
	@docker build . -t gcr.io/${MIMOSA_GCP_PROJECT}/runner
	@gcloud services enable --project ${MIMOSA_GCP_PROJECT} containerregistry.googleapis.com
	@docker push gcr.io/${MIMOSA_GCP_PROJECT}/runner
	@gcloud beta run deploy --image gcr.io/${MIMOSA_GCP_PROJECT}/runner --platform managed --region europe-west1 --no-allow-unauthenticated runner
