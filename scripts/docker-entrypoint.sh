#!/bin/bash

set -e

arelo -d 2s -p 'configs/*.yml' -p '**/*.go' -i '**/*_test.go' --  bash -c "go run -tags dynamic cmd/mpj-apiserver/*.go $@"

