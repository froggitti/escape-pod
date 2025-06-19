#!/bin/bash -e

update-binfmts --enable qemu-aarch64

cp /work/scripts/qemu-run.sh ./img/qemu-run.sh
chroot img ./qemu-run.sh
rm ./img/qemu-run.sh
rm ./img/default-intent-list.json