# Local onboarding

## Bot config

## Install the required binaries. 

1.  SSH into your bot and run this:

```sh
$ mount -o remount rw /
```

2.  Copy the two files from the [latest vector-cloud release](https://github.com/DDLbots/vector-cloud/releases/) into /anki/bin.  the permissions should be as followed:

```sh
-rwxr-xr-x    1 root     root       5368336 Oct 28 14:28 vic-cloud
-rwxr-xr-x    1 root     root       4376548 Oct 28 14:28 vic-gateway
```

Or just run this after scping the files over

```sh
mv /anki/bin/vic-cloud /anki/bin/vic-cloud.orig
mv /anki/bin/vic-gateway /anki/bin/vic-gateway.orig
mv vic-cloud /anki/bin
mv vic-gateway /anki/bin
chown cloud:anki /anki/bin/vic-cloud
chown net:anki /anki/bin/vic-gateway
chmod 755 /anki/bin/vic-gateway
chmod 755 /anki/bin/vic-cloud
```

3.  Edit your server_config.json

```sh
$ vi /anki/data/assets/cozmo_resources/config/server_config.json
```

The contents should be as followed:

```json
{
        "jdocs": "escapepod.local:8084",
        "tms": "escapepod.local:8084",
        "chipper": "escapepod.local:8084",
        "check": "escapepod.local:8080/ok",
        "logfiles": "s3://anki-device-logs-prod/victor",
        "appkey": "oDoa0quieSeir6goowai7f"
}
```

4.  Install the certificate

Add a file to /etc/ssl/certs/local/root.crt with the following contents:

```sh
-----BEGIN CERTIFICATE-----
MIIDzzCCAregAwIBAgIUfZNgw6BEADyNnpsskl0/e1yTFtEwDQYJKoZIhvcNAQEL
BQAwdzELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAlBBMRMwEQYDVQQHDApQaXR0c2J1
cmdoMRswGQYDVQQKDBJEaWdpdGFsIERyZWFtIExhYnMxKTAnBgkqhkiG9w0BCQEW
GmJyZXR0QGRpZ2l0YWxkcmVhbWxhYnMuY29tMB4XDTIwMTAyODE4MDcxN1oXDTI1
MTAyNzE4MDcxN1owdzELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAlBBMRMwEQYDVQQH
DApQaXR0c2J1cmdoMRswGQYDVQQKDBJEaWdpdGFsIERyZWFtIExhYnMxKTAnBgkq
hkiG9w0BCQEWGmJyZXR0QGRpZ2l0YWxkcmVhbWxhYnMuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAlf+BSaO02Di7IVfudT5PAiUeUF/Pe/6HUTL6
UKTvJA+dMTfdxeF5CjP0ivm6yb7k1vG7YJOt71zdtJeacLbVcPGy3O5Y+JCO07K1
BxIFysDbeT5a/S5Xksk2fjwFdxWDpRmR4Zux4wnH9kZH8Of6fvaLwCMM/AV9lL4D
ZEareSlj+ipawX0RI7nXZZBQKZpOyVpygrLln5XfLyJ6ksZofKYsiMw+equtpzv5
ONXWD2D/45zI7U56Kh7/Oj2FDc6VEJr0qnXGUROQua43GiTGfma9Vt6I9dsZGTlO
PCtPQv0/ZcZo1vtvVWQBrFDAJdRl9GeHtRioeZ60MK/DVysttwIDAQABo1MwUTAd
BgNVHQ4EFgQUNa2OFIgy6SpnACeH8BH1S99U/n4wHwYDVR0jBBgwFoAUNa2OFIgy
6SpnACeH8BH1S99U/n4wDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC
AQEARVxs5RAyrzMNMTnYeionpWXmkiVyDb2m77mMahPcq1ByB3g3V/+XC+lKopFS
l30pPV3K3npWY1oShgdP5w9+lkxL7nCu5dTdEIyyIx8ZocoaneRQEqsSA8+mvSCr
fvQI8+4ZRerWRhii4FYQPohHithhUE+1+s7LiPP6IgQ7ZdgbAbiYR/eXRnSbxenj
864LmPWJqq383Hr40quGLjHNPu2KA/k+qGqm3tA+LfHW0odyu/H0vbmLD5ih7tQC
x+e8c7UBe2sN++yqQnQ0Kbg441ze9+4dWYlPfzBOtnjj0+DjJjawUwQa7rFE0ixH
valKKA0LkLqYJpvK0970NY8X3Q==
-----END CERTIFICATE-----
```

5. Force the hostname

echo YOUR_IP >> /etc/hosts

6. Clear user settings and reboot


## License generation

[Follow these instructions](licensing.md)

Make sure the file is in the same directory that you're calling the binary from

## Onboarding

1.  cd node
2.  npm install
3.  node vector-web-setup serve
4.  Open your browser to localhost:8000