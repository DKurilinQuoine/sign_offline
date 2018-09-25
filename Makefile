BUILD_FLAGS = -gcflags=all='-N -l'
PROJECT_NAME = sign_offline
SOURCE_LIST = ./src/sign_offline.go
DOCKER_IMAGE_NAME = tmp_docker_image

all: $(PROJECT_NAME)

$(PROJECT_NAME): deps $(SOURCE_LIST)
	go build $(BUILD_FLAGS) -o build/${PROJECT_NAME} $(SOURCE_LIST) 

docker:	Dockerfile $(SOURCE_LIST)
	docker build -t $(DOCKER_IMAGE_NAME) .
	docker run --rm 					\
			   -w /go/src/sign_offline	\
			   $(DOCKER_IMAGE_NAME)		\
			   make
	docker image rm -f $(DOCKER_IMAGE_NAME)

clear:
	rm -rf ./build

# dep
DEP = github.com/golang/dep/cmd/dep

DEP_CHECK := $(shell command -v dep 2> /dev/null)

deps:
ifndef DEP_CHECK
	go get -v $(DEP)
endif
	dep ensure -v
#==================

