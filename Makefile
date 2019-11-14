
#
# The COMPONENTS and corresponding targets allow us to run targets in parallel
#
# e.g. "make deploy -j" deploys all components in parallel
#

COMPONENTS := \
	api \
	sources/aws \
	sources/gcp \
	sources/vmpooler \
	system/reusabolt \
	system/router \
	system/usermgmt \
	worldbuilders/inventory

DEPLOYABLE := \
	api \
	system/reusabolt \
	system/usermgmt \
	worldbuilders/inventory

DEPLOYABLE_COMPONENTS = $(addprefix deploy-, $(DEPLOYABLE))

.PHONY: $(COMPONENTS)

all: $(COMPONENTS)
test: $(COMPONENTS)
deploy: $(DEPLOYABLE_COMPONENTS)

$(COMPONENTS):
	@echo
	@echo "🔘 ===> Component $@ ..."
	@$(MAKE) -C $@ $(MAKECMDGOALS)

$(DEPLOYABLE_COMPONENTS):
	@echo
	@echo "🔘 ===> Deploying $(@:deploy-%=%) ..."
	@$(MAKE) -C $(@:deploy-%=%) $(MAKECMDGOALS)
