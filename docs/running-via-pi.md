# Running via a pi

1.  Bot configuration

Either [configure your bot](bot-configuration.md) or flash the escapepod firmware.

2.  Image flashing

If you don't have (or want) a pre-built image, run the following command:

```sh
# make build-image
```

Download something like [balenaEtcher](https://www.balena.io/etcher/) and simply flash the image to an SD card.

3.  WIFI configuration

Mount the "system-boot" partition on the SD card and edit the network-settings file.  For a fairly standard WIFI config, uncomment and edit this section:

```sh
wifis:
  wlan0:
    dhcp4: true
    optional: true
    access-points:
      <wifi network name>:
        password: "<wifi password>"
```

If you're using a LAN, you don't have to do much here.

4.  Boot!

Put the card into a pi and turn it on.

5.  Onboard your bot

[Follow these instructions](onboarding.md)