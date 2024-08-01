.DEFAULT_GOAL := help

# CGO_ENABLED=0 is not FIPS compliant. large commercial vendors and FedRAMP require FIPS compliant crypto
CGO_ENABLED := 1

# Enable users to override the golang used to accomodate custom installations
GO ?= go

# Allow overriding `oc` command.
# Used by pr_check.py to ssh deploy inside private Hive cluster via bastion host.
oc:=oc

# The version needs to be different for each deployment because otherwise the
# cluster will not pull the new image from the internal registry:
version:=$(shell date +%s)
# Tag for the image:
image_tag:=$(version)

# The namespace and the environment are calculated from the name of the user to
# avoid clashes in shared infrastructure:
environment:=${USER}
namespace:=ocm-${USER}

# a tool for managing containers and images, etc. You can set it as docker
container_tool ?= podman

# In the development environment we are pushing the image directly to the image
# registry inside the development cluster. That registry has a different name
# when it is accessed from outside the cluster and when it is accessed from
# inside the cluster. We need the external name to push the image, and the
# internal name to pull it.
external_apps_domain?=apps-crc.testing
external_image_registry?=default-route-openshift-image-registry.$(external_apps_domain)
internal_image_registry?=image-registry.openshift-image-registry.svc:5000

# The name of the image repository needs to start with the name of an existing
# namespace because when the image is pushed to the internal registry of a
# cluster it will assume that that namespace exists and will try to create a
# corresponding image stream inside that namespace. If the namespace doesn't
# exist the push fails. This doesn't apply when the image is pushed to a public
# repository, like `docker.io` or `quay.io`.
image_repository:=$(namespace)/rh-trex

# Database connection details
db_name:=rhtrex
db_host=trex-db.$(namespace)
db_port=5432
db_user:=trex
db_password:=foobar-bizz-buzz
db_password_file=${PWD}/secrets/db.password
db_sslmode:=disable
db_image?=docker.io/library/postgres:14.2

# Log verbosity level
glog_v:=10

# Location of the JSON web key set used to verify tokens:
jwks_url:=https://sso.redhat.com/auth/realms/redhat-external/protocol/openid-connect/certs

# Test output files
unit_test_json_output ?= ${PWD}/unit-test-results.json
integration_test_json_output ?= ${PWD}/integration-test-results.json

# Prints a list of useful targets.
help:
	@echo ""
	@echo "OpenShift CLuster Manager Example Service"
	@echo ""
	@echo "make verify               verify source code"
	@echo "make lint                 run golangci-lint"
	@echo "make binary               compile binaries"
	@echo "make install              compile binaries and install in GOPATH bin"
	@echo "make run                  run the application"
	@echo "make run/docs             run swagger and host the api spec"
	@echo "make test                 run unit tests"
	@echo "make test-integration     run integration tests"
	@echo "make generate             generate openapi modules"
	@echo "make image                build docker image"
	@echo "make push                 push docker image"
	@echo "make deploy               deploy via templates to local openshift instance"
	@echo "make undeploy             undeploy from local openshift instance"
	@echo "make project              create and use an Example project"
	@echo "make clean                delete temporary generated files"
	@echo "$(fake)"
.PHONY: help

# Encourage consistent tool versions
OPENAPI_GENERATOR_VERSION:=5.4.0
GO_VERSION:=go1.21.

### Constants:
version:=$(shell date +%s)
GOLANGCI_LINT_BIN:=$(shell go env GOPATH)/bin/golangci-lint

### Envrionment-sourced variables with defaults
# Can be overriden by setting environment var before running
# Example:
#   OCM_ENV=testing make run
#   export OCM_ENV=testing; make run
# Set the environment to development by default
ifndef OCM_ENV
	OCM_ENV:=development
endif

ifndef TEST_SUMMARY_FORMAT
	TEST_SUMMARY_FORMAT=short-verbose
endif

ifndef OCM_BASE_URL
	OCM_BASE_URL:="https://api.integration.openshift.com"
endif

# Checks if a GOPATH is set, or emits an error message
check-gopath:
ifndef GOPATH
	$(error GOPATH is not set)
endif
.PHONY: check-gopath

# Verifies that source passes standard checks.
verify: check-gopath
	${GO} vet \
		./cmd/... \
		./pkg/...
	! gofmt -l cmd pkg test |\
		sed 's/^/Unformatted file: /' |\
		grep .
	@ ${GO} version | grep -q "$(GO_VERSION)" || \
		( \
			printf '\033[41m\033[97m\n'; \
			echo "* Your go version is not the expected $(GO_VERSION) *" | sed 's/./*/g'; \
			echo "* Your go version is not the expected $(GO_VERSION) *"; \
			echo "* Your go version is not the expected $(GO_VERSION) *" | sed 's/./*/g'; \
			printf '\033[0m'; \
		)
.PHONY: verify

# Runs our linter to verify that everything is following best practices
# Requires golangci-lint to be installed @ $(go env GOPATH)/bin/golangci-lint
# Linter is set to ignore `unused` stuff due to example being incomplete by definition
lint:
	$(GOLANGCI_LINT_BIN) run -e unused \
		./cmd/... \
		./pkg/...
.PHONY: lint

# Build binaries
# NOTE it may be necessary to use CGO_ENABLED=0 for backwards compatibility with centos7 if not using centos7
binary: check-gopath
	${GO} build ./cmd/trex
.PHONY: binary

# Install
install: check-gopath
	CGO_ENABLED=$(CGO_ENABLED) GOEXPERIMENT=boringcrypto ${GO} install -ldflags="$(ldflags)" ./cmd/trex
	@ ${GO} version | grep -q "$(GO_VERSION)" || \
		( \
			printf '\033[41m\033[97m\n'; \
			echo "* Your go version is not the expected $(GO_VERSION) *" | sed 's/./*/g'; \
			echo "* Your go version is not the expected $(GO_VERSION) *"; \
			echo "* Your go version is not the expected $(GO_VERSION) *" | sed 's/./*/g'; \
			printf '\033[0m'; \
		)
.PHONY: install

# Runs the unit tests.
#
# Args:
#   TESTFLAGS: Flags to pass to `go test`. The `-v` argument is always passed.
#
# Examples:
#   make test TESTFLAGS="-run TestSomething"
test: install
	OCM_ENV=testing gotestsum --format short-verbose -- -p 1 -v $(TESTFLAGS) \
		./pkg/... \
		./cmd/...
.PHONY: test

# Runs the unit tests with json output
#
# Args:
#   TESTFLAGS: Flags to pass to `go test`. The `-v` argument is always passed.
#
# Examples:
#   make test-unit-json TESTFLAGS="-run TestSomething"
ci-test-unit: install
	@echo $(db_password) > ${PWD}/secrets/db.password
	OCM_ENV=testing gotestsum --jsonfile-timing-events=$(unit_test_json_output) --format short-verbose -- -p 1 -v $(TESTFLAGS) \
		./pkg/... \
		./cmd/...
.PHONY: ci-test-unit

# Runs the integration tests.
#
# Args:
#   TESTFLAGS: Flags to pass to `go test`. The `-v` argument is always passed.
#
# Example:
#   make test-integration
#   make test-integration TESTFLAGS="-run TestAccounts"     acts as TestAccounts* and run TestAccountsGet, TestAccountsPost, etc.
#   make test-integration TESTFLAGS="-run TestAccountsGet"  runs TestAccountsGet
#   make test-integration TESTFLAGS="-short"                skips long-run tests
ci-test-integration: install
	@echo $(db_password) > ${PWD}/secrets/db.password
	OCM_ENV=testing gotestsum --jsonfile-timing-events=$(integration_test_json_output) --format $(TEST_SUMMARY_FORMAT) -- -p 1 -ldflags -s -v -timeout 1h $(TESTFLAGS) \
			./test/integration
.PHONY: ci-test-integration

# Runs the integration tests.
#
# Args:
#   TESTFLAGS: Flags to pass to `go test`. The `-v` argument is always passed.
#
# Example:
#   make test-integration
#   make test-integration TESTFLAGS="-run TestAccounts"     acts as TestAccounts* and run TestAccountsGet, TestAccountsPost, etc.
#   make test-integration TESTFLAGS="-run TestAccountsGet"  runs TestAccountsGet
#   make test-integration TESTFLAGS="-short"                skips long-run tests
test-integration: install
	@echo $(db_password) > ${PWD}/secrets/db.password
	OCM_ENV=testing gotestsum --format $(TEST_SUMMARY_FORMAT) -- -p 1 -ldflags -s -v -timeout 1h $(TESTFLAGS) \
			./test/integration
.PHONY: test-integration

# Regenerate openapi client and models
generate:
	rm -rf pkg/api/openapi
	$(container_tool) build -t ams-openapi -f Dockerfile.openapi .
	$(eval OPENAPI_IMAGE_ID=`$(container_tool) create -t ams-openapi -f Dockerfile.openapi .`)
	$(container_tool) cp $(OPENAPI_IMAGE_ID):/local/pkg/api/openapi ./pkg/api/openapi
	$(container_tool) cp $(OPENAPI_IMAGE_ID):/local/data/generated/openapi/openapi.go ./data/generated/openapi/openapi.go
.PHONY: generate

run: install
	trex migrate
	trex serve
.PHONY: run

# Run Swagger and host the api docs
run/docs:
	@echo "Please open http://localhost/"
	docker run -d -p 80:8080 -e SWAGGER_JSON=/trex.yaml -v $(PWD)/openapi/rh-trex.yaml:/trex.yaml swaggerapi/swagger-ui
.PHONY: run/docs

# Delete temporary files
clean:
	rm -rf \
		$(binary) \
		templates/*-template.json \
		data/generated/openapi/*.json \
.PHONY: clean

.PHONY: cmds
cmds:
	for cmd in $$(ls cmd); do \
		${GO} build \
			-ldflags="$(ldflags)" \
			-o "$${cmd}" \
			"./cmd/$${cmd}" \
			|| exit 1; \
	done


# NOTE multiline variables are a PITA in Make. To use them in `oc process` later on, we need to first
# export them as environment variables, then use the environment variable in `oc process`
%-template:
	oc process \
		--filename="templates/$*-template.yml" \
		--local="true" \
		--ignore-unknown-parameters="true" \
		--param="ENVIRONMENT=$(OCM_ENV)" \
		--param="GLOG_V=$(glog_v)" \
		--param="DATABASE_HOST=$(db_host)" \
		--param="DATABASE_NAME=$(db_name)" \
		--param="DATABASE_PASSWORD=$(db_password)" \
		--param="DATABASE_PORT=$(db_port)" \
		--param="DATABASE_USER=$(db_user)" \
		--param="DATABASE_SSLMODE=$(db_sslmode)" \
		--param="IMAGE_REGISTRY=$(internal_image_registry)" \
		--param="IMAGE_REPOSITORY=$(image_repository)" \
		--param="IMAGE_TAG=$(image_tag)" \
		--param="VERSION=$(version)" \
		--param="AUTHZ_RULES=$$AUTHZ_RULES" \
		--param="ENABLE_SENTRY"=false \
		--param="SENTRY_KEY"=TODO \
		--param="JWKS_URL=$(jwks_url)" \
		--param="OCM_SERVICE_CLIENT_ID=$(CLIENT_ID)" \
		--param="OCM_SERVICE_CLIENT_SECRET=$(CLIENT_SECRET)" \
		--param="TOKEN=$(token)" \
		--param="OCM_BASE_URL=$(OCM_BASE_URL)" \
		--param="ENVOY_IMAGE=$(envoy_image)" \
		--param="ENABLE_JQS="false \
	> "templates/$*-template.json"


.PHONY: project
project:
	$(oc) new-project "$(namespace)" || $(oc) project "$(namespace)" || true

.PHONY: image
image: cmds
	$(container_tool) build -t "$(external_image_registry)/$(image_repository):$(image_tag)" .

.PHONY: push
push:	\
	image \
	project
	$(container_tool) push "$(external_image_registry)/$(image_repository):$(image_tag)" --tls-verify=false

deploy-%: project %-template
	$(oc) apply --filename="templates/$*-template.json" | egrep --color=auto 'configured|$$'

undeploy-%: project %-template
	$(oc) delete --filename="templates/$*-template.json" | egrep --color=auto 'deleted|$$'


.PHONY: template
template: \
	secrets-template \
	db-template \
	service-template \
	route-template \
	$(NULL)

# Depending on `template` first helps clustering the "foo configured", "bar unchanged",
# "baz deleted" messages at the end, after all the noisy templating.
.PHONY: deploy
deploy: \
	push \
	template \
	deploy-secrets \
	deploy-db \
	deploy-service \
	deploy-route \
	$(NULL)

.PHONY: undeploy
undeploy: \
	template \
	undeploy-secrets \
	undeploy-db \
	undeploy-service \
	undeploy-route \
	$(NULL)

.PHONY: db/setup
db/setup:
	@echo $(db_password) > $(db_password_file)
	$(container_tool) run --name psql-rhtrex -e POSTGRES_DB=$(db_name) -e POSTGRES_USER=$(db_user) -e POSTGRES_PASSWORD=$(db_password) -p $(db_port):5432 -d $(db_image)

.PHONY: db/login
db/login:
	$(container_tool) exec -it psql-rhtrex bash -c "psql -h localhost -U $(db_user) $(db_name)"

.PHONY: db/teardown
db/teardown:
	$(container_tool) stop psql-rhtrex
	$(container_tool) rm psql-rhtrex

crc/login:
	@echo "Logging into CRC"
	@crc console --credentials -ojson | jq -r .clusterConfig.adminCredentials.password | oc login --username kubeadmin --insecure-skip-tls-verify=true https://api.crc.testing:6443
	@oc whoami --show-token | $(container_tool) login --username kubeadmin --password-stdin "$(external_image_registry)" --tls-verify=false
.PHONY: crc/login
