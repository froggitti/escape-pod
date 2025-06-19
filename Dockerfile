FROM golang:buster as builder

RUN apt update
# libc6 is for make utils
# libopus* are for audio codex
# libsodium-dev ?? idk what it is for
# libsystemd-dev is to access journalctl
# git is for `go get`
RUN apt install libc6 libopus0 libopus-dev libopusfile-dev libopusfile0 libsodium-dev libsystemd-dev git -y

# COPY OVER FILES
ADD . /go/src/github.com/DDLbots/escape-pod
WORKDIR /go/src/github.com/DDLbots/escape-pod

ARG GITHUB_USER="not-set"
ARG GITHUB_TOKEN="not-set"

RUN echo "machine github.com login ${GITHUB_USER} password ${GITHUB_TOKEN}" > ~/.netrc

ENV GOPRIVATE github.com/anki,github.com/DDLbots

RUN make install-escapepod

FROM debian:buster
RUN apt-get update && apt-get install -y apt-transport-https

RUN apt-get install libc6 libopus0 libopus-dev libopusfile-dev libopusfile0 libsodium-dev libsystemd-dev git -y

ENV LD_LIBRARY_PATH=/usr/lib/
ENV ROOT_DIRECTORY="/usr/local/escapepod" 
ENV BLE_LOG_DIRECTORY="logs" 
ENV OTA_DIRECTORY="ota" 
ENV UI_DIRECTORY="/dist" 
ENV JDOCS_FILEPATH="jdocs.json" 
ENV INTENTS_FILEPATH="intents.json" 
ENV LICENSES_FILEPATH="licenses.json" 
ENV DEFAULT_INTENTS_FILEPATH="./default-intents/default-intent-list.json" 
ENV STT_MODEL="/usr/local/escapepod/model.tflite" 
ENV STT_SCORER="/usr/local/escapepod/large_vocabulary.scorer" 
ENV ESCAPEPOD_EXTENDER="" 
ENV ESCAPEPOD_EXTENDER_TARGET="" 
ENV ESCAPEPOD_EXTENDER_DISABLE_TLS="" 
ENV GRPC_HOST="" 
ENV GRPC_PORT=443 
ENV UI_PORT=80 
ENV GRPC_TLS=true 
ENV LOG_LEVEL=debug

ADD coqui/linux-tflite-amd64/* /usr/lib/
ADD modules/escape-pod-ui/dist/* /usr/local/escapepod/dist/

ADD coqui/large_vocabulary.scorer /usr/local/escapepod/large_vocabulary.scorer
ADD coqui/model.tflite /usr/local/escapepod/model.tflite

COPY --from=builder /go/src/github.com/DDLbots/escape-pod/bin/escape-pod /usr/local/escapepod/bin/escape-pod
ENV PATH="/usr/local/escapepod/bin:$PATH"

ENTRYPOINT [ "escape-pod" ]