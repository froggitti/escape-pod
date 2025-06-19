
# RPI instructions 
These instructions are for building your own, not for using a pre-built image.  These instructions are mostly for internal users, unless we choose to release the source for this.

These were tested on a raspberry pi 3 bv1.2, because its what I had lying around.

## Flash the image to the SD card

You'll want to install ubuntu server 20.10 64 bit.

##  Preconfigure your system

More in-depth instructions can be found [here](https://ubuntu.com/tutorials/how-to-install-ubuntu-on-your-raspberry-pi#1-overview)

1.  Configure wifi.  

Edit the network-config file to look something like this:

```sh
wifis:
  wlan0:
    dhcp4: true
    optional: true
    access-points:
      <wifi network name>:
        password: "<wifi password>"
```
2.  Add a file named "ssh" to the same directory to enable remote access

```sh
# touch ssh
```

## Configure mongodb

see pi3-installation.sh

# Install and run the escape pod software

1.  Run "build-pi-3" from the root of this repo, which will make an arm binary.  We'll probably need multiple targets here because not all PIs are GOARM=5.

2.  Copy the file to the pi
```sh
# scp escape-pod pi@[IP ADDRESS]:`~/.
```

3.  Ping Brett to get the environment configurations

4.  Run it!

```sh
# ./escape-pod
```

TODO:

* Automate all of this
* Make systemd stuff
* roll entire installable images