# licensing

## explanation

The licenses are now stored in the db vs on the filesystem.  This allows for webpage interactions vs direct filesystem interactions.

Before carrying on, generate a license as followed:

```sh
# license-generator -email brett@digitaldreamlabs.com -robots vic:00507778
```

Note:  previous versions of the license generator are no longer compatible, and you can no longer generate a license for multiple bots.  One license per bot!

### defintions

The proto for the GRPC interface can be found [here](https://github.com/DDLbots/escape-pod/blob/master/internal/license/interceptor/proto/interceptor.proto)

### rest endpoints

* http://escapepod.local:8085/v1/license/add

This will add a license for a bot.  The format is as followed

```json
{
	"license" : "base64stringthatslong"
}
```

* http://escapepod.local:8085/v1/license/list

This provides a list of all authorized bots


* http://escapepod.local:8085/v1/license/delete

This will delete a license for a bot.  The format is as followed:
```json
{
	"bot" : "vic:1234abcd"
}
```

### Extra info

If an end user tampers with the license info (which have signatures, etc) in the database, we shut down ALL THE THINGS.  No bots will function, not just the bot with the tampered license.