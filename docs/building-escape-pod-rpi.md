

## Naming Images
```
${distro}-${distro_version}-${service}-${service_version}.${rc|other metadta}?.img
ubuntu-20.04.3-escape-pod-1.8.rc4.img
```


## Building Manually Linux
```sh
# NOTE: you don't need root for this.
# You can do this part from the working directory of escape-pod
# sudo apt install upx

sudo su -

source /home/joseph/go/src/github.com/DDLbots/escape-pod/modules/escape-pod-arm-image-builder/scripts/builder.sh

unpack_factory_image img_dir=/home/joseph/go/src/github.com/DDLbots/escape-pod/modules/escape-pod-arm-image-builder/os_images

exit

cd $GOPATH/src/github.com/DDLbots/escape-pod

# make build-pi

sudo mkdir -p `pwd`/modules/escape-pod-arm-image-builder/mnt

cd modules/escape-pod-arm-image-builder

sudo ./scripts/image-mount ./os_images/copy_factory_arm.img mnt linux

cd -

# If you don't have podman; replace it with docker

sudo docker build -t escape-pod-arm-image-builder -f `pwd`/modules/escape-pod-arm-image-builder .

sudo docker run -it --privileged \
    -e OS_FACTORY_IMAGE_UNPACKED=copy_factory_arm.img \
    -e OS_IMAGES_DIR=os_images \
    -e LOOP=/dev/loop0 \
    -e TMP=/mnt \
    -v `pwd`/modules/escape-pod-arm-image-builder/mnt:/mnt \
    -v `pwd`/modules/escape-pod-arm-image-builder/scripts:/builder/scripts \
    -v `pwd`/image-builder/files:/builder/service_files \
    -v `pwd`/modules/escape-pod-arm-image-builder/os_images:/builder/os_images \
    -v `pwd`/default-intents:/builder/default_intents \
    -v `pwd`/bin:/builder/escape_pod \
    -v `pwd`/modules/escape-pod-ui:/builder/escape_pod_ui \
    -v `pwd`/coqui:/builder/deepspeech \
    -v `pwd`/coqui/linux-tflite-aarch64:/builder/deepspeech/lib \
    -v `pwd`/image-builder/firmwares:/builder/firmwares \
    escape-pod-arm-image-builder \
    ./scripts/build_arm.sh


# please unmount before compressing you might corrupt the partitions
sudo umount `pwd`/modules/escape-pod-arm-image-builder/mnt

# write the image or compress it; i like to test it before i compress

```