
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
	worldbuilders/inventory

.PHONY: $(COMPONENTS)

all: $(COMPONENTS)
test: $(COMPONENTS)
deploy: $(COMPONENTS)

$(COMPONENTS):
	@echo
	@echo "🔘 ===> Component $@ ..."
	@$(MAKE) -C $@ $(MAKECMDGOALS)
