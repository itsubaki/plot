SHELL := /bin/bash

cos:
	go run main.go testdata/cos.csv

bloch:
	go run main.go --sphere --swap testdata/bloch.csv
