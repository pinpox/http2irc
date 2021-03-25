# http2irc

## Configuration

Configuration is done with environment variables. The following options are
available:

| Variable      | Example               | Description                    |
|---------------|-----------------------|--------------------------------|
| IRC_NICK      | mybot                 | Nickname                       |
| IRC_CHANNEL   | #mychannel            | Channel to connect to          |
| IRC_SERVER    | irc.freenode.net:7000 | IRC server and port            |
| IRC_SASL_USER | myuser                | SASL user                      |
| IRC_SASL_PASS | verysecret            | SASL password                  |
| IRC_LISTEN    | localhost:8080        | Listening address              |
| IRC_DEBUG     | false                 | Verbose output                 |
| IRC_NOTICE    | true                  | Use notice instead of messages |
| IRC_TEMPLATE  | ./example.tmpl        | Path to the template           |
| IRC_BOT_TOKEN | myverysecrettoken     | Token for the bot              |

## Write a message

```bash
curl localhost:8989/webhook 
--header "Content-Type: application/json" \
--header  "Token: myverysecrettoken" \
--request POST \
--data '{"data":"test"}' 
```
