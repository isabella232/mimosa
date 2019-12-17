
#
# The COMPONENTS and corresponding targets allow us to run targets in parallel
#
# e.g. "make deploy -j" deploys all components in parallel
#

COMPONENTS := \
	api \
	sources \
	docker \
	system/evaluator \
	system/reusabolt \
	system/router \
	system/usermgmt \
	actions/reaper \
	worldbuilders/inventory \
	worldbuilders/vulns

.PHONY: $(COMPONENTS)

all: $(COMPONENTS)
test: $(COMPONENTS)
deploy: $(COMPONENTS)

$(COMPONENTS):
	@echo
	@echo "🔘 ===> Component $@ ..."
	@$(MAKE) -C $@ $(MAKECMDGOALS)

UI-TARGETS := test-ui build-ui

deploy-ui: $(UI-TARGETS)
	cd ui/mimosa-ui && \
	test-ui && build-ui

ui-local:
	cd ui/mimosa-ui && \
	yarn start;

build-ui:
	cd ui/mimosa-ui && \
	yarn build;

test-ui:
	cd ui/mimosa-ui && \
	yarn test --watchAll=false;
