#!/usr/bin/env just

set shell := ["zsh", "-cu"]
set dotenv-load := true

DOCKER_REPOSITORY := "wot"
TAG := `echo ${TAG:=latest}`
NAME := "logwild"
GIT_COMMIT := `git describe --dirty --always`
VERSION := `grep 'VERSION' pkg/version/version.go | awk '{ print $4 }' | tr -d '"'`
# default to no extra run args, but allow them to be supplied in env
EXTRA_RUN_ARGS := `echo ${EXTRA_RUN_ARGS:=""}`

help :
  @echo "TAG: {{ TAG }}"
  @echo "NAME: {{ NAME }}"
  @echo "DOCKER_REPOSITORY: {{ DOCKER_REPOSITORY }}"
  @echo "GIT_COMMIT: {{ GIT_COMMIT }}"
  @echo "VERSION: {{ VERSION }}"
  @just --list

build :
  CGO_ENABLED=0 go build -ldflags "-s -w -X mcgaunn.com/logwild/pkg/version.REVISION={{ GIT_COMMIT }}" \
    -a -o ./bin/logwild ./cmd/logwild/*

build-container :
  @echo "this should build docker container"

build-charts :
  @echo "this should build helm charts"
  helm lint charts/*
  helm package charts/*

push-container :
  @echo "tagging and pushing image"

version-set :
  #!/bin/zsh -exu
  next="{{ TAG }}"
  current="{{ VERSION }}"
  /usr/bin/sed -i "s/$current/$next/g" pkg/version/version.go
  echo "Version $next set in code"

release:
  git tag -s -m {{ VERSION }} {{ VERSION }}
  git push alert {{ VERSION }}

run :
  go run -ldflags "-s -w -X mcgaunn.com/logwild/pkg/version.REVISION={{ GIT_COMMIT }}" \
    cmd/logwild/* --debug run {{ EXTRA_RUN_ARGS }}

fmt :
  gofmt -l -s -w ./

tidy :
  rm -f go.sum; go mod tidy -compat=1.21

vet :
  go vet ./...

test :
  go test -v ./... -coverprofile cover.out

# run the app with appropriate arguments to forward traces, metrics, logs
# to collector
run-app :
  OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:44317 \
    bin/logwild \
    --debug \
    --port-metrics 8889 \
    --port 8888 \
    --otel-service-name 'logwild-imcg' \
    --log-rate 10000 \
    --log-bytes 1024 \
    --log-burst-duration 5 \
    --log-out-file /tmp/logwild.log \
    run {{ EXTRA_RUN_ARGS }}

# run standalone openobserve service to receive traces, metrics, logs from collector
run-observe-backend :
  scripts/run_standalone_observe_backend.sh

# run standalone otel collector listening on grpc ports for traces/logs
run-otel-collector :
  scripts/run_standalone_collector.sh

run-in-panels :
  scripts/run_everything_in_panels.sh
