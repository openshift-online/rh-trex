#!/bin/bash
tmpdir="$(mktemp -d)"
echo "creating app $tmpdir"

go run cmd/trex/main.go clone --name myapp --destination "/tmp/${tmpdir}" --repo-base github.com/openshift-online

cd "/tmp/${tmpdir}"

go run ./scripts/generator.go --kind Pet --fields "name:string,color:string"

make generate

make binary

make test

# disable RYUK for podman
export TESTCONTAINERS_RYUK_DISABLED=true
make test-integration
