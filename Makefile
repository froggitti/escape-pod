.PHONY: build-escapepod build-pi build-image release get-models-scorers

DEEPSPEECH_VERSION=0.9.1

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git rev-parse --short HEAD)
CWD := $(shell cd -P -- '$(shell dirname -- "$0")' && pwd -P)
SSH_PRIVATE_KEY=`cat ~/.ssh/id_rsa`

VERSION ?= $(shell ./scripts/git-version)
COMMIT ?= $(shell ./scripts/git-commit)

UI_VERSION ?= $(shell ./scripts/ui-version)
UI_COMMIT ?= $(shell ./scripts/ui-commit)

REPO_PATH=github.com/DDLbots/escape-pod/internal

LD_FLAGS= -X $(REPO_PATH)/version.Version=$(VERSION) \
-X $(REPO_PATH)/version.Build=$(COMMIT) \
-X $(REPO_PATH)/version.UIVersion=$(UI_VERSION) \
-X $(REPO_PATH)/version.UIBuild=$(UI_COMMIT)

GOARCH=$(shell go env GOARCH)

get-models-scorers:	deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.pbmm deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.scorer deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.tflite
	
deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.pbmm:
	curl -L https://github.com/mozilla/DeepSpeech/releases/download/v$(DEEPSPEECH_VERSION)/deepspeech-$(DEEPSPEECH_VERSION)-models.pbmm -o deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.pbmm

deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.scorer:
	curl -L https://github.com/mozilla/DeepSpeech/releases/download/v$(DEEPSPEECH_VERSION)/deepspeech-$(DEEPSPEECH_VERSION)-models.scorer -o deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.scorer

deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.tflite:
	curl -L https://github.com/mozilla/DeepSpeech/releases/download/v$(DEEPSPEECH_VERSION)/deepspeech-$(DEEPSPEECH_VERSION)-models.tflite -o deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.tflite

get-coqui-tflite-linux-aarch64:
	mkdir -p ./coqui/linux-tflite-aarch64
	wget https://github.com/coqui-ai/STT/releases/download/v1.2.0/native_client.tflite.linux.aarch64.tar.xz -O ./coqui/linux-tflite-aarch64/native_client.tflite.linux.aarch64.tar.xz
	tar -xf ./coqui/linux-tflite-aarch64/native_client.tflite.linux.aarch64.tar.xz -C ./coqui/linux-tflite-aarch64

get-coqui-tflite-linux-amd64:
	mkdir -p ./coqui/linux-tflite-amd64
	wget https://github.com/coqui-ai/STT/releases/download/v1.2.0/native_client.tflite.Linux.tar.xz -O ./coqui/linux-tflite-amd64/native_client.tflite.Linux.tar.xz
	tar -xf ./coqui/linux-tflite-amd64/native_client.tflite.Linux.tar.xz -C ./coqui/linux-tflite-amd64

build-escapepod: get-models-scorers
	echo $(LD_LIBRARY_PATH)

	CGO_LDFLAGS="-L$(CWD)/coqui/linux-tflite-amd64/" \
	CGO_CXXFLAGS="-I$(CWD)/coqui/linux-tflite-amd64/" \
	LD_LIBRARY_PATH=$(CWD)/coqui/linux-tflite-amd64/:$(LD_LIBRARY_PATH) \
	CGO_ENABLED=1 \
	go build \
	-tags production \
	-trimpath \
	-a -installsuffix cgo \
	-ldflags '-w -s -r $(CWD)/deepspeech/linux-amd64/ $(LD_FLAGS)' \
	-o ./bin/escape-pod-linux-amd64 cmd/escapepod/main.go

build-escapepod-linux: build-escapepod-prod
	sudo setcap 'cap_net_raw,cap_net_admin+eip' ./bin/escape-pod-linux-amd64

install-escapepod:
	CGO_LDFLAGS="-L$(CWD)/coqui/linux-tflite-amd64/" \
	CGO_CXXFLAGS="-I$(CWD)/coqui/linux-tflite-amd64/" \
	LD_LIBRARY_PATH=$(CWD)/coqui/linux-tflite-amd64/:$(LD_LIBRARY_PATH) \
	CGO_ENABLED=1 \
	go build \
	-tags image,nolibopusfile,production \
	-ldflags "-w -s $(LD_FLAGS)" \
	-trimpath \
	-o bin/escape-pod-linux-arm64 cmd/escapepod/*.go
	
	# upx ./bin/escape-pod-linux-arm64


build-escapepod-prod: 
	CGO_LDFLAGS="-L$(CWD)/deepspeech/linux-amd64/" \
	CGO_CXXFLAGS="-I$(CWD)/deepspeech/include/" \
	LD_LIBRARY_PATH=$(CWD)/deepspeech/linux-amd64/:$(LD_LIBRARY_PATH) \
	go build \
	-trimpath \
	-tags image,nolibopusfile,production \
	-a -installsuffix cgo \
	-ldflags '-w -s -r $(CWD)/deepspeech/linux-amd64/ $(LD_FLAGS)' \
	-o ./bin/escape-pod-linux-amd64 cmd/escapepod/*.go
	# upx escape-pod

build-docker:
	docker build \
	--build-arg GITHUB_USER="$(GITHUB_TOKEN)" \
	--build-arg GITHUB_TOKEN="$(GITHUB_TOKEN)" \
	-t 367813035514.dkr.ecr.us-east-2.amazonaws.com/escape-pod:test . \
	-f dockerfiles/Dockerfile.escapepod

build-pi:
	docker build \
	--build-arg GITHUB_USER="$(GITHUB_TOKEN)" \
	--build-arg GITHUB_TOKEN="$(GITHUB_TOKEN)" \
	-t armbuilder-es:latest . \
	-f dockerfiles/Dockerfile.armbuilder

	docker container run -it --rm \
	-v "$(PWD)":/go/src/digital-dream-labs/escape-pod \
	-v $(GOPATH)/pkg/mod:/go/pkg/mod \
	-w /go/src/digital-dream-labs/escape-pod \
	--user $(UID):$(GID) \
	armbuilder-es:latest \
	go build \
	-tags image,nolibopusfile \
	-ldflags "-w -s $(LD_FLAGS)" \
	-trimpath \
	-o bin/escape-pod-linux-arm64 cmd/escapepod/*.go
	
	upx ./bin/escape-pod-linux-arm64

build-pi-release:
	docker build \
	--build-arg GITHUB_USER="$(GITHUB_TOKEN)" \
	--build-arg GITHUB_TOKEN="$(GITHUB_TOKEN)" \
	-t armbuilder-es:latest . \
	-f dockerfiles/Dockerfile.armbuilder

	docker container run -it --rm \
	-v "$(PWD)":/go/src/digital-dream-labs/escape-pod \
	-v $(GOPATH)/pkg/mod:/go/pkg/mod \
	-w /go/src/digital-dream-labs/escape-pod \
	--user $(UID):$(GID) \
	armbuilder-es:latest \
	go build \
	-tags image,nolibopusfile,production \
	-ldflags "-w -s $(LD_FLAGS)" \
	-trimpath \
	-o bin/escape-pod-linux-arm64 cmd/escapepod/*.go
	
	upx ./bin/escape-pod-linux-arm64

build-image: 
	docker build \
	-t pibuilder:latest . \
	-f dockerfiles/Dockerfile.pibuilder

	cp deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.tflite image-builder/files/models.tflite
	cp deepspeech/deepspeech-$(DEEPSPEECH_VERSION)-models.scorer image-builder/files/models.scorer
	cp -rpf deepspeech/include image-builder/files/include
	mkdir -p image-builder/files/lib/
	mkdir -p image-builder/img/
	cp -rpf deepspeech/linux-arm64/* image-builder/files/lib/.

	sudo ./image-builder/build.sh
	mv image-builder/pi-64.img escape-pod-$(COMMIT).img
	sudo xz -T0 -z escape-pod-$(COMMIT).img

build-license-generator:
	CGO_ENABLED=0 \
	go build \
	-ldflags "-w -s -extldflags "-static"" \
	-trimpath \
	-o license-generator cmd/license/main.go
	upx license-generator

build-license-generator-release:
	CGO_ENABLED=0 \
	go build \
	-ldflags "-w -s -extldflags "-static"" \
	-trimpath \
	-tags production \
	-o license-generator cmd/license/main.go
	upx license-generator

release: 
	mkdir -p release
	GOOS=windows GOARCH=amd64 make build-license-generator-release
	mv license-generator release/license-generator-win-amd64-$(COMMIT)
	GOOS=darwin GOARCH=amd64 make build-license-generator-release
	mv license-generator release/license-generator-darwin-amd64-$(COMMIT)
	GOOS=linux GOARCH=amd64 make build-license-generator-release
	mv license-generator release/license-generator-linux-amd64-$(COMMIT)
	make build-pi-release
	make build-image
	mv escape-pod-$(COMMIT).img.xz release

clean:
	rm -fr image-builder/files/node
	rm -fr image-builder/files/lib/*
	rm -fr image-builder/files/include/*
	rm -f image-builder/files/escape-pod
	rm -f image-builder/pi-64.img
	rm -fr release/*
	rm -f *.upx
	sudo umount image-builder/img/dev
	sudo umount image-builder/img/sys
	sudo umount image-builder/img/proc
	sudo umount image-builder/img/boot/firmware
	sudo umount image-builder/img

version:
	echo $(VERSION)

run:
	LD_LIBRARY_PATH=$(CWD)/coqui/linux-tflite-amd64/:$(LD_LIBRARY_PATH) \
	ROOT_DIRECTORY="." \
	BLE_LOG_DIRECTORY="/home/joseph/go/src/github.com/DDLbots/escape-pod/logs" \
	OTA_DIRECTORY="ota" \
	UI_DIRECTORY="modules/escape-pod-ui/dist" \
	JDOCS_FILEPATH="jdocs.json" \
	INTENTS_FILEPATH="intents.json" \
	LICENSES_FILEPATH="licenses.json" \
	DEFAULT_INTENTS_FILEPATH="./default-intents/default-intent-list.json" \
	XSTT_MODEL="$(CWD)/coqui/model.tflite" \
	XSTT_SCORER="$(CWD)/coqui/large_vocabulary.scorer" \
	STT_SCORER="$(CWD)/deepspeech/deepspeech-0.9.1-models.scorer" \
	STT_MODEL="$(CWD)/deepspeech/deepspeech-0.9.1-models.tflite" \
	ESCAPEPOD_EXTENDER="" \
	ESCAPEPOD_EXTENDER_TARGET="" \
	ESCAPEPOD_EXTENDER_DISABLE_TLS="" \
	GRPC_HOST="" \
	GRPC_PORT=443 \
	UI_PORT=80 \
	GRPC_TLS=true \
	LOG_LEVEL=debug \
	./bin/escape-pod-linux-amd64

license:
	@CGO_ENABLED=0 \
	go run \
	-ldflags "-w -s -extldflags "-static"" \
	-trimpath \
	-tags production \
	./cmd/license -email=stephanotter@gmail.com -robot=vic:0dd1dcab > 0dd1dcab.key



licenses:
