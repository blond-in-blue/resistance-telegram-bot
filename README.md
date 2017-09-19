# Smartest Telegram Bot

Telegram bot built in go.

A lot of code pulled from [here](https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129).

## Development

To spin up a debug container with hot reloading, make a duplicate of the `sample.Dockerfile.debug` and name it simply `Dockerfile.debug`, adding in the apprioriate keys to the environment variable.


```
docker build --rm -t smartest-telegram-bot -f Dockerfile.debug .
docker run -p 3000:3000 -v C:/dev/projects/StartNode/smartest-reddits/app:/go/src/github.com/user/myProject/app --name tele-bot smartest-telegram-bot 
```

## Production:

```
docker build --rm -t smartest-telegram-bot-prod .
docker run -p 80:80 --name tele-bot-prod smartest-telegram-bot-prod 
```