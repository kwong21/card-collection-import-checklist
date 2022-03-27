#!/bin/bash
set -x 

go test -v integration_test.go function.go
