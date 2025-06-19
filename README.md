# ESCAPE POD!!

You can...
* [run the bot via docker](docs/running-via-docker.md)
* [run the bot locally](docs/running-locally.md)
* [run the bot from a PI](docs/running-via-pi.md)

# Building and running

## Pre-requisites

Aside from having docker installed, you'll need to do the following
1.   Install the deepspeech model scorers
	```sh
	$ make get-models-scorers
	```
	
2.  Download any firmware images you'll want distributed with the image to the image-builder/firmwares directory.

3.  Download an [ubuntu 64-bit factory image](https://ubuntu.com/download/raspberry-pi), and with xz-utils installed, run "unxz" against the file, and place it in image-builder/factory.img

# Running locally via docker on linux - Linux only!!!!

```sh
sudo apt install libopus0 libopus-dev libopusfile-dev libsystemd-dev libsodium-dev
# if building for coqui
export CGO_LDFLAGS="-L`pwd`/coqui/linux-tflite-amd64/" 
export CGO_CXXFLAGS="-I`pwd`/coqui/linux-tflite-amd64/" 
export LD_LIBRARY_PATH=`pwd`/coqui/linux-tflite-amd64/:$LD_LIBRARY_PATH

make make get-coqui-tflite-linux-amd64

# DEV SETUP optional if you don't have this setup
git submodule update --init --recursive
# echo "machine github.com login ${GITHUB_USER} password ${GITHUB_TOKEN}" > ~/.netrc
# export GOPRIVATE=github.com/anki,github.com/DDLbots
# You WILL need to add a production_key.go file in both the /internal/license/issuer and /internal/license/validator

# IMPORTANT!!!
# NOTE: if running avahi for mdns set host-name to escapepod in /etc/avahi/avahi-daemon.conf    
# Example
# [server]
# host-name=escapepod


mkdir -p image-builder/files/escape-pod-ui/ota
mkdir -p image-builder/files/escape-pod-ui/logs

make build-escapepod-linux && sudo make run 
```

## Running locally (not via docker)

1.  Make a local environment file named ".env". This should suffice:

	```yaml
	export DDL_RPC_PORT=8084
	export DDL_HTTP_PORT=8085
	export DDL_OTA_PORT=8086
	export DDL_HTTPS_PORT=8443
	export DDL_RPC_CLIENT_AUTHENTICATION="RequestClientCert"
	export BLE_LOG_DIRECTORY="."
	export DDL_DB_NAME=database
	export DDL_DB_HOST=localhost
	export DDL_DB_PASSWORD=password
	export DDL_DB_PORT=27017
	export DDL_DB_USERNAME=root
	export DDL_DB_DIRECT=true
	export DDL_SAYWHATNOW_STT_MODEL="./deepspeech/deepspeech-0.9.1-models.tflite"
	export DDL_SAYWHATNOW_STT_SCORER="./deepspeech/deepspeech-0.9.1-models.scorer"
	export DDL_SAYWHATNOW_VAD_TIMEOUT=10
	export CGO_LDFLAGS="-L$(pwd)/deepspeech/linux-amd64/"
	export CGO_CXXFLAGS="-I$(pwd)/deepspeech/include/"
	export LD_LIBRARY_PATH=$(pwd)/deepspeech/linux-amd64/:$LD_LIBRARY_PATH
	export ENABLE_EXTENSIONS=true
	export ESCAPEPOD_EXTENDER_TARGET=localhost:8089
	export ESCAPEPOD_EXTENDER_DISABLE_TLS=true
	```

2.  Start mongo however you wish.  You'll need to turn on clustering for this to function.  Full configuration options can be found [here](https://github.com/DDLbots/escape-pod/blob/8d923df84b80f3e26235e188c20df002dcd43d86/image-builder/scripts/qemu-run.sh#L33)

3.  Run the go binary
```sh
go run -tags production cmd/escapepod/main.go
```

// TODO: this needs to update!
# Building

All building is handled via docker, so this part should be pretty easy.  You will be prompted for a root password, as some of the steps (mounting virtual images, etc) are done via privileged calls.

```sh
$ make release
```
All release files will be built and the output will be to the "release" directory.


#To down load the UI as a submodule - escapepod-ui
```sh
git submodule update --init --recursive # --force if you have old file
```