
include ../build/common.mk

# FIXME remove this once the GCP lint errors are fixed
lint:
	@echo "No linting until GCP errors are fixed ..."

SOURCE_NAME:=source-${shell uuidgen | awk '{print tolower($0)}'}

PHONY+= create
create:
ifndef WORKSPACE
	$(error WORKSPACE must be defined)
endif
ifndef CONFIG_FILE
	$(error CONFIG_FILE must be defined)
endif
	@echo "🔘 Creating ${SOURCE_NAME} ... (`date '+%H:%M:%S'`)"
	@echo "🔘 Creating bucket ${SOURCE_NAME} ... (`date '+%H:%M:%S'`)"
	@gsutil mb -b on gs://${SOURCE_NAME}
	@echo "🔘 Labelling bucket ... (`date '+%H:%M:%S'`)"
	@gsutil label ch -l ws:${WORKSPACE} gs://${SOURCE_NAME}
	@echo "🔘 Copying config to bucket ... (`date '+%H:%M:%S'`)"
	@gsutil cp ${CONFIG_FILE} gs://${SOURCE_NAME}/config.json
	@echo "🔘 Enabling notifications ... (`date '+%H:%M:%S'`)"
	@gsutil notification create -f json -t system-router gs://${SOURCE_NAME}
	@echo "✅ Created ${SOURCE_NAME} (`date '+%H:%M:%S'`)"

delete:
ifndef NAME
	$(error NAME must be defined)
endif
	@echo "🔘 Deleting bucket ${NAME} ... (`date '+%H:%M:%S'`)"
	@-gsutil rm -r gs://${NAME}
	@echo "✅ Deleted bucket ${NAME} (`date '+%H:%M:%S'`)"

# Use "make deploy -j" for parallel deployment
PHONY+= deploy
deploy: deploy-aws deploy-gcp deploy-vmpooler deploy-qualys
	@echo "✅ All deployments finished (`date '+%H:%M:%S'`)"

PHONY+= deploy-aws
deploy-aws:
	${call deploy-source,source-aws,AWS}

PHONY+= deploy-gcp
deploy-gcp:
	${call deploy-source,source-gcp,GCP}

PHONY+= deploy-qualys
deploy-qualys:
ifndef MIMOSA_GCP_PROJECT
	$(error MIMOSA_GCP_PROJECT must be defined)
endif
	${call deploy-source,source-qualys,Qualys}
	@echo "🔘 Creating bucket ${MIMOSA_GCP_PROJECT}-blueray (if it doesn't exist)... (`date '+%H:%M:%S'`)"
	@gsutil ls -b gs://${MIMOSA_GCP_PROJECT}-blueray || gsutil mb -b on gs://${MIMOSA_GCP_PROJECT}-blueray
	@echo "🔘 Copying qualys.xml to bucket ${MIMOSA_GCP_PROJECT}-blueray (`date '+%H:%M:%S'`)"
	@gsutil cp ../blueray/qualys.xml gs://${MIMOSA_GCP_PROJECT}-blueray

PHONY+= deploy-vmpooler
deploy-vmpooler:
	${call deploy-source,source-vmpooler,VMPooler}

define deploy-source
	@echo "🔘 Deploying ${1} ... (`date '+%H:%M:%S'`)"
	@if ! gcloud pubsub topics describe ${1} &> /dev/null; then gcloud pubsub topics create ${1}; fi
	@gcloud functions deploy \
		--runtime go111 \
		--no-allow-unauthenticated \
		--trigger-topic ${1} \
		--source . \
		--entry-point=${2} \
		${1}
	@echo "✅ Deployed ${1} (`date '+%H:%M:%S'`)"
endef
