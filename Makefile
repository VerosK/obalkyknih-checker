.PHONY: image push

TARGET_IMAGE_NAME ?= localhost/obalkyknih-checker
TARGET_IMAGE_TAG  ?= latest
TARGET_IMAGE_VERSION_ID ?= "unknown"

CONTAINER_ENGINE := $(shell command -v podman 2> /dev/null | echo docker)


IMAGE_TAG = latest
IMAGE_NAME = $(TARGET_IMAGE_NAME):$(TARGET_IMAGE_TAG)

SRC_ROOT_PATH = $(CURDIR)

image:
	echo "Running build in $(SRC_ROOT_PATH)"
	$(CONTAINER_ENGINE) build -f Containerfile \
			--build-arg VERSION_ID=$(TARGET_IMAGE_VERSION_ID) \
			-t $(IMAGE_NAME) $(SRC_ROOT_PATH)

push: image
	$(CONTAINER_ENGINE) push $(IMAGE_NAME)
