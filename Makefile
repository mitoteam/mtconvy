# Set up variables
DIST_DIR := dist
EXECUTABLE_NAME := mtconvy
ARCH_FILES := ./${DIST_DIR}/.mtconvy.yml

#do the job
include internal/goappbase/Makefile.inc.mk

# extend targets
before-build::
	cp ./.mtconvy.DEFAULT.yml ./${DIST_DIR}/.mtconvy.yml

after-build::
	rm -f ./${DIST_DIR}/.mtconvy.yml

# additional targets
.PHONY: run
run:
	clear
	go run main.go
