package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type Config struct {
	BotToken       string        `json:"telegramBotToken"`
	CheckFrequency time.Duration `json:"checkFrequency"`
	URLs           []string      `json:"urls"`
}

var (
	config           Config
	parallelWorkers  = 4
	pool             chan Worker
	botMutex         *sync.Mutex
	bot              *tgbotapi.BotAPI
	messageReceived        = false
	notifiableChatID int64 = 0
)

func init() {
	runtime.GOMAXPROCS(parallelWorkers)
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Println(err)
		return
	}
	body, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(body, &config)
	if err != nil {
		log.Println(err)
		return
	}
	config.CheckFrequency = time.Duration(
		int(config.CheckFrequency * time.Second))
	pool = make(chan Worker, parallelWorkers)
	for i := 0; i < parallelWorkers; i++ {
		pool <- Worker{}
	}
	botMutex = &sync.Mutex{}
	bot, err = tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		log.Println(err)
		return
	}
	bot.Debug = false
	fmt.Printf("Authorized as %s\n", bot.Self.UserName)
	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates, err := bot.GetUpdatesChan(u)

		if err != nil {
			log.Println(err)
			return
		}

		update := <-updates

		botMutex.Lock()
		messageReceived = true
		notifiableChatID = update.Message.Chat.ID
		botMutex.Unlock()

		log.Printf(
			"Sending reports to chat %d\n",
			notifiableChatID,
		)

		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID, "Chat registered")
		bot.Send(msg)
	}()
}

func main() {
	for {
		for _, url := range config.URLs {
			curWorker := <-pool
			go curWorker.CheckURL(url, pool)
		}
		time.Sleep(config.CheckFrequency)
	}
}

//sendMsg sends message by bot to id specified in config
func sendMsg(message string) {
	botMutex.Lock()
	defer botMutex.Unlock()
	if !messageReceived {
		log.Println(
			"Write a message to the bot for specifying notifiable chat ID",
		)
		return
	}
	msg := tgbotapi.NewMessage(notifiableChatID, message)
	bot.Send(msg)
}
