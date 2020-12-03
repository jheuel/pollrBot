# pollrBot
[![Donate](https://camo.githubusercontent.com/11b2f47d7b4af17ef3a803f57c37de3ac82ac039/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f70617970616c2d646f6e6174652d79656c6c6f772e737667)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=Q43CKDHGRVSML)
[![Donate](https://camo.githubusercontent.com/c19db43a081a84a33e1bec7e4d454f801b6e2628/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f626974636f696e2d646f6e6174652d79656c6c6f772e737667)](https://commerce.coinbase.com/checkout/6bf4d01e-c638-41d9-9ac2-4a2aaf1beba9)

This is a telegram bot that helps by creating inline polls in telegram chats
without spamming multiple messages.

It is online and can be used at [@pollrBot](https://telegram.me/pollrBot).

The bot uses inline queries and feedback to inline queries, which have to be
enabled with the telegram [@BotFather](https://telegram.me/BotFather).

## Usage
The bot can be installed with
```
go get github.com/jheuel/pollrBot
```
if you have a working Go environment.


After that you can run the bot with
```
URL="https://pollr.yourdomain.com" DB="database.db" APITOKEN="euiaeouiaouiao" pollrBot
```

## Statistics of pollrBot
Here are some numbers about how many users interacted with my instance of the
pollrBot and how many polls were created. The instance is running next to a few
other things on the smallest droplet you can get from digitalocean.com.

![Graph of number of users and polls vs time](stats.png)
