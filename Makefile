SHELL := /bin/bash

cos:
	go run main.go testdata/cos.csv

bloch:
	go run main.go --scatter --swap-xy testdata/htcover16.csv
	go run main.go --scatter --swap-xy testdata/htcover24.csv
	go run main.go --scatter --swap-xy testdata/htcover32.csv
