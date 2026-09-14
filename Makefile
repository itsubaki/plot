SHELL := /bin/bash

update:
	go get -u
	go mod tidy

cos:
	go run main.go testdata/cos.csv
