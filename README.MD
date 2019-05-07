# Website downtime notifier

---

Utility for checking availability of specified websites and notifying about
errors via telegram bot.

## Building

###### Depencies:
* [Telegram-bot-api](https://github.com/go-telegram-bot-api/telegram-bot-api)
* [Golang 1.10+](https://golang.org/)

###### Building process
* Clone repository `git clone https://github.com/corelof/downtime-notifier`
* Inside of `downtme-notifier` directory run `go build *.go`

## Launching

* Rename `config_example.json` to `config.json`
* Fill all `config.json` fields with correct info(telegram api token, etc.)
* Specify URL's to check in config file
* Run previously compiled binary file(`main`)

Note: `checkFrequency` config parameter must be specified in seconds.

## Team

* [corelof](https://github.com/corelof)
