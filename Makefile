# Set up variables
EXECUTABLE_NAME := mtconvy

#do the job
include internal/goappbase/Makefile.inc.mk

# additional targets
.PHONY: run
run:
	clear
	go run main.go

# additional commands
build-all::
	sha256sum ${DIST_DIR}/*.7z > ${DIST_DIR}/checksums.sha256
