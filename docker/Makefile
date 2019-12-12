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

PHONY+= mimosabuild
mimosabuild:
	@echo "🔘 Creating 'mimosabuild' container (`date '+%H:%M:%S'`)"
	@docker build build -t gcr.io/${MIMOSA_GCP_PROJECT}/mimosabuild
	@docker push gcr.io/${MIMOSA_GCP_PROJECT}/mimosabuild
	@echo "✅ Created 'mimosabuild' container (`date '+%H:%M:%S'`)"

PHONY+= enable-cloud-run-api
enable-cloud-run-api:
ifndef MIMOSA_GCP_PROJECT
	$(error MIMOSA_GCP_PROJECT must be defined)
endif
	@gcloud services enable --project ${MIMOSA_GCP_PROJECT} containerregistry.googleapis.com

PHONY+= build-container
build-container:
ifndef MIMOSA_GCP_PROJECT
	$(error MIMOSA_GCP_PROJECT must be defined)
endif
	${call build-container,runner}
# ${call build-container,reusabolt}

define build-container
	@echo "🔘 Building container ${1} (`date '+%H:%M:%S'`)"
	@docker build ${1} -t gcr.io/${MIMOSA_GCP_PROJECT}/${1}
	@echo "✅ Built container ${1} (`date '+%H:%M:%S'`)"
endef

PHONY+= deploy-container
deploy-container: enable-cloud-run-api build-container
ifndef MIMOSA_GCP_PROJECT
	$(error MIMOSA_GCP_PROJECT must be defined)
endif
	${call deploy-container,runner}
# ${call deploy-container,reusabolt}

define deploy-container
	@echo "🔘 Deploying container ${1} (`date '+%H:%M:%S'`)"
	@docker push gcr.io/${MIMOSA_GCP_PROJECT}/${1}
	@gcloud beta run deploy --image gcr.io/${MIMOSA_GCP_PROJECT}/${1} --platform managed --region europe-west1 --no-allow-unauthenticated ${1}
	@echo "✅ Deployed container ${1} (`date '+%H:%M:%S'`)"
endef
