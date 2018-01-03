# Smartest Telegram Bot

Telegram bot built in go.

A lot of code pulled from [here](https://github.com/go-telegram-bot-api/telegram-bot-api/blob/13c54dc548f7ca692fe434d4b7cac072b0de0e0b/types.go#L129).

## Development + Deployment

To spin up a container with hot reloading, make your own file *my.env* and place it the root of this project. Put in it these variables:
```
TELE_KEY=<key given from fatherbot>
REDDIT_CLIENT_ID=<id of app generated on reddit>
REDDIT_CLIENT_SECRET=<secret of the app>
REDDIT_USERNAME=<your reddit username>
REDDIT_PASSWORD=<your reddit password>
```

To actually run the app:
```
docker-compose up
```