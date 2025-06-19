#!/bin/bash -e

UBUNTUIMAGE=https://cdimage.ubuntu.com/releases/20.04.1/release/ubuntu-20.04.1-preinstalled-server-arm64+raspi.img.xz
DIRNAME="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

##################################################################################################
# Disk setup                                                                                     #
# This part mounts the downloaded pi image on a loopback device so the dockerized qemu           #
# stuff can read/write it                                                                        #
##################################################################################################

if [ ! -f $DIRNAME/factory.img ]; then
    curl ${UBUNTUIMAGE} -o $DIRNAME/factory.img.xz
    xz --decompress $DIRNAME/factory.img.xz
    qemu-img resize $DIRNAME/factory.img +2G
    export DEVICE=$(sudo losetup -P --find --show factory.img)
    sudo resize2fs -p /dev/$DEVICEp2
fi
cp $DIRNAME/factory.img $DIRNAME/pi-64.img

LOOPDEVICE=`losetup -f -P --show $DIRNAME/pi-64.img`
mount ${LOOPDEVICE}p2 $DIRNAME/img
mount ${LOOPDEVICE}p1 $DIRNAME/img/boot/firmware   
# mount -t bind /run $DIRNAME/img/run/
mount -t proc /proc $DIRNAME/img/proc/
mount -t sysfs /sys $DIRNAME/img/sys/
mount -o bind /dev $DIRNAME/img/dev/


##################################################################################################
# Pre-config                                                                                     #
# This part mounts copies local files into the disk image for use later                          #
##################################################################################################

# enable BLE and ssh
cp $DIRNAME/files/user-data $DIRNAME/img/boot/firmware/user-data
echo "dtparam=i2c_arm=on" >> $DIRNAME/img/boot/firmware/syscfg.txt
echo "dtparam=spi=on" >> $DIRNAME/img/boot/firmware/syscfg.txt
echo "include btcfg.txt" >> $DIRNAME/img/boot/firmware/usercfg.txt
touch $DIRNAME/img/boot/firmware/ssh

echo escapepod > $DIRNAME/img/etc/hostname
mkdir -p $DIRNAME/img/usr/local/escapepod/bin

# deepspeech requirements
cp $DIRNAME/files/lib/* $DIRNAME/img/usr/lib
cp $DIRNAME/files/models.tflite $DIRNAME/img/usr/local/escapepod/model.tflite
cp $DIRNAME/files/models.scorer $DIRNAME/img/usr/local/escapepod/model.scorer

# escape pod specific
cp $DIRNAME/files/escape-pod $DIRNAME/img/usr/local/escapepod/bin/escape-pod
chmod +x $DIRNAME/img/usr/local/escapepod/bin
cp -rpf $DIRNAME/files/escape-pod.conf $DIRNAME/img/etc/escape-pod.conf
cp $DIRNAME/files/escape_pod.service $DIRNAME/img/lib/systemd/system
cp -rpf $DIRNAME/files/default-intent-list.json $DIRNAME/img/default-intent-list.json
mkdir -p $DIRNAME/img/usr/local/escapepod/ui
mkdir -p $DIRNAME/img/usr/local/escapepod/ui/ota
cp $DIRNAME/firmwares/* $DIRNAME/img/usr/local/escapepod/ui/ota
mkdir -p $DIRNAME/img/usr/local/escapepod/ui/logs
cp -rpf $DIRNAME/files/escape-pod-ui/dist/* $DIRNAME/img/usr/local/escapepod/ui

# execute in-image configuration
docker run --privileged=true -it \
    --mount type=bind,source=$DIRNAME/img,target=/work/img \
    --mount type=bind,source=$DIRNAME/scripts,target=/work/scripts \
    pibuilder /work/scripts/run.sh

# clean up
umount $DIRNAME/img/proc
umount $DIRNAME/img/sys
umount $DIRNAME/img/dev

umount $DIRNAME/img/boot/firmware
umount $DIRNAME/img
losetup -d ${LOOPDEVICE}

# shrink

sudo bash $DIRNAME/scripts/pishrink.sh $DIRNAME/pi-64.img
