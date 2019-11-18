
include ../build/common.mk

lint:
	@echo LINT DISABLED FOUR SOURCES UNTIL THE GCP ERRORS ARE FIXED

check-env:
ifndef UUID
	$(error UUID must be defined)
endif
ifndef SOURCE
	$(error SOURCE must be defined)
endif
ifndef CONFIG_FILE
	$(error CONFIG_FILE must be defined)
endif

NAME=source-${UUID}

PHONY+= deploy
deploy: check-env
	@echo "🔘 Copying config to bucket ... (`date '+%H:%M:%S'`)"
	gsutil cp ${CONFIG_FILE} gs://${NAME}/config.json
	@echo "🔘 Deploying ${NAME} ... (`date '+%H:%M:%S'`)"
	@gcloud functions deploy \
		--runtime go111 \
		--trigger-topic ${NAME} \
		--set-env-vars MIMOSA_GCP_BUCKET=${NAME}, \
		--source ${SOURCE} \
		--entry-point=HandleMessage \
		${NAME}
	@echo "✅ Deployed system-reusabolt (`date '+%H:%M:%S'`)"

# Separating these dependencies allows parallel deployment
GENERATED_UUID := ${shell uuidgen| awk '{print tolower($0)}'}
create-source:
	@echo "🔘 Creating source-${GENERATED_UUID} ... (`date '+%H:%M:%S'`)"
	@gsutil mb -b on gs://source-${GENERATED_UUID}
	@gcloud pubsub topics create source-${GENERATED_UUID}
	@echo "✅ Created source-${GENERATED_UUID} (`date '+%H:%M:%S'`)"

delete-source:
ifndef UUID
	$(error UUID must be defined)
endif
	@echo "🔘 Deleting source-${UUID} ... (`date '+%H:%M:%S'`)"
	@-gsutil rm -r gs://source-${UUID}
	@-gcloud pubsub topics delete source-${UUID}
	@-gcloud functions delete --quiet source-${UUID}
	@-gcloud functions delete --quiet system-router-${UUID}
	@echo "✅ Deleted source-${UUID} (`date '+%H:%M:%S'`)"

# Separating these dependencies allows parallel deployment
deploy-source: --deploy-source-function --deploy-router

--deploy-source-function:
	$(MAKE) -C sources deploy UUID=${UUID}

--deploy-router:
	$(MAKE) -C system/router deploy UUID=${UUID}
