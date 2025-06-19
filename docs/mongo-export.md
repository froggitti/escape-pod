# exporting a database for an image build

As you're most likely testing mongo stuff in a container, here's the easiest way to export it

1.  Get the container ID

```sh
# docker ps
```

You'll see something like this:

```
CONTAINER ID        IMAGE                COMMAND                  CREATED             STATUS              PORTS                      NAMES
a2870e41d028        mongo:4.4.1-bionic   "docker-entrypoint.s…"   3 minutes ago       Up 3 minutes        0.0.0.0:27017->27017/tcp   escapepod-mongo
```

The container ID is, in this case, a2870e41d028

2.  Dump the contents

```sh
# docker exec -it a2870e41d028 mongoexport \
  --collection intents -d database -u root -p password \
  --authenticationDatabase admin \
  --pretty --jsonArray > filename.json
```

Just remember to chop off the first and last lines in the file, or it'll be invalid!