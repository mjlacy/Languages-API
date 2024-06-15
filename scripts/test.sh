#!/usr/bin/env bash

if [[ "${1}" == "func" || "${1}" == "html" ]]; then
  go test ./... -coverprofile coverage.out
  go tool cover -${1}=coverage.out
  rm coverage.out
else
  go test ./...
fi