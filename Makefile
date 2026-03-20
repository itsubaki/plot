SHELL := /bin/bash

cos:
	go run main.go testdata/cos.csv

bloch:
	go run main.go --scatter --swap testdata/bloch16.csv
	go run main.go --scatter --swap testdata/bloch24.csv
	go run main.go --scatter --swap testdata/bloch32.csv
