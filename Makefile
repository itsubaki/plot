SHELL := /bin/bash

cos:
	go run main.go testdata/cos.csv

bloch:
	go run main.go --scatter --swap testdata/bloch.csv
