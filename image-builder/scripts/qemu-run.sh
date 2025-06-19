#!/bin/bash -e

##################################################################################################
# Mongo configuration                                                                            #
# Pre-configure the mongo database so the user doesn't have to                                   #
#                                                                                                #
# Perhaps eventually this can be replaced with a per-instance password generation that gets      #
# triggered on reboot.  It'd definitely be more secure.  It'd have to dynamically update         #
# /etc/escape-pod.conf                                                                           #
##################################################################################################

rm /etc/resolv.conf
echo 'nameserver 8.8.4.4' | sudo tee -a /etc/resolv.conf

apt update
dpkg -r snapd

apt install curl gnupg avahi-daemon libopus0 libsystemd-dev libsodium-dev pi-bluetooth -y

curl -s https://www.mongodb.org/static/pgp/server-4.4.asc | apt-key add -
echo "deb [ arch=arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/4.4 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-4.4.list
apt update
apt upgrade -y
apt install mongodb-org -y

echo -e "replication:
  replSetName: rs0\n\n" >> /etc/mongod.conf

mongod --config /etc/mongod.conf &

sleep 5

mongo << EOF
rsconf = {
  _id: "rs0",
  members: [
    {
     _id: 0,
     host: "localhost:27017"
    },
   ]
}
rs.initiate(rsconf)
EOF

sleep 5

mongo << EOF
use admin
db.createUser(
  {
    user: "myUserAdmin",
    pwd: "MzBmMWFmY2NhYzE0",
    roles: [ { role: "userAdminAnyDatabase", db: "admin" }, "readWriteAnyDatabase" ]
  }
)
EOF

sleep 2

mongoimport --uri "mongodb://myUserAdmin:MzBmMWFmY2NhYzE0@127.0.0.1/database?authsource=admin" \
  --collection intents \
  --file /default-intent-list.json \
  --jsonArray

mongo -u myUserAdmin -p MzBmMWFmY2NhYzE0 --authenticationDatabase admin database << EOF
db.licenses.createIndex( { "license.bot": 1.0 }, { unique: true } )
EOF

chown -R mongodb:mongodb /var/lib/mongodb
chown -R mongodb:mongodb /var/log/mongodb
rm /var/lib/mongodb/*.lock

openssl rand -base64 666 > /etc/mongod.key
chmod 0600 /etc/mongod.key
chown mongodb:mongodb /etc/mongod.key

echo -e "security:
  authorization: enabled
  keyFile: /etc/mongod.key
" >> /etc/mongod.conf

#### daemonize things
systemctl enable mongod
systemctl enable avahi-daemon.service
systemctl enable escape_pod.service

chown -R 1000:1000 /usr/local/escapepod

cat > /etc/profile.d/00-aliases.sh << EOF
alias showlogs='journalctl -u escape_pod.service'
alias showtext='journalctl -u escape_pod.service | grep incoming_text'
alias taillogs='journalctl -u escape_pod.service -f'
alias tailtext='journalctl -u escape_pod.service -f | grep incoming_text'
EOF

echo "127.0.0.1    escapepod.local escapepod" >> /etc/hosts
echo "::1    escapepod.local escapepod" >> /etc/hosts

#### various customizations
echo -e "APT::Periodic::Update-Package-Lists \"1\";\nAPT::Periodic::Unattended-Upgrade \"0\";" > /etc/apt/apt.conf.d/20auto-upgrades
echo SystemMaxUse=500M >> /etc/systemd/journald.conf

rm /etc/resolv.conf
ln -s ../run/systemd/resolve/stub-resolv.conf /etc/resolv.conf