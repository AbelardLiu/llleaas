
ifeq ($(BUILD_VERBOSE),1)
    Q =
else
    Q = @
endif

all: add-excutable unix-format binary image

build-binary:
    $(Q)cd builder/; ./compile.sh

build-image:
    $(Q)cd builder/; ./docker.sh

clean-binary:
    $(Q)cd builder/; ./uncompile.sh

clean-image:
    $(Q)cd builder/; ./undocker.sh

build-tests:
    $(Q)cd builder/; ./test.sh

.PHONY: clean
clean: clean-image clean-binary

image: clean-image build-image

binary: clean-binary build-binary

tests: build-tests

unix-format:
    $(Q)dos2unix builder/*

add-executable:
    $(Q)chmod +x builder/*.sh
    $(Q)chmod -R +x builder/compile
    $(Q)chmod -R +x builder/docker