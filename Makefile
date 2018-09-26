BUILD_FLAGS = -gcflags=all='-N -l'
PROJECT_NAME = sign_offline
SOURCE_LIST = ./src/main.go
DOCKER_IMAGE_NAME = tmp_docker_image

all: $(PROJECT_NAME)

$(PROJECT_NAME): deps $(SOURCE_LIST)
	go build $(BUILD_FLAGS) -o build/${PROJECT_NAME} $(SOURCE_LIST) 

docker:	Dockerfile $(SOURCE_LIST)
	docker build -t $(DOCKER_IMAGE_NAME) .
	mkdir -p build
	docker run \
			   -w /go/src/sign_offline	\
			   -v $(shell pwd)/build:/go/src/sign_offline/build \
			   $(DOCKER_IMAGE_NAME)		\
			   make test
#docker image rm -f $(DOCKER_IMAGE_NAME)

clear:
	rm -rf ./build

test: $(PROJECT_NAME)
	touch build/test_file
	go test -v ./src 2>&1 | go-junit-report > ./build/test_report.xml

# dep
DEP = github.com/golang/dep/cmd/dep

DEP_CHECK := $(shell command -v dep 2> /dev/null)

deps:
ifndef DEP_CHECK
	go get -v $(DEP)
endif
	dep ensure -v
#==================

