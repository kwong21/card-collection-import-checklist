#!/bin/bash
set -x 

export GOOGLE_PROJECT_ID=test
export FIRESTORE_EMULATOR_HOST=localhost:8080

go test -v integration_test.go function.go
