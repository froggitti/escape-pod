#!/bin/sh
This won't actually run yet.

curl -s https://www.mongodb.org/static/pgp/server-4.4.asc | sudo apt-key add -
echo "deb [ arch=arm64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.4 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.2.list
sudo apt update
sudo apt install mongodb-org -y
sudo systemctl enable mongod
sudo systemctl start mongod

# generate random password, something like date +%s | sha256sum | base64 | head -c 16 ; echo

# run this
mongo
use admin
db.createUser(
  {
    user: "myUserAdmin",
    pwd: "MzBmMWFmY2NhYzE0",
    roles: [ { role: "userAdminAnyDatabase", db: "admin" }, "readWriteAnyDatabase" ]
  }
)

add this to the end of /etc/mongod.conf
security:
    authorization: enabled


sudo systemctl restart mongod
