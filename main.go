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
	BotToken         string        `json:"telegramBotToken"`
	NotifiableChatID int64         `json:"notifiableChatID"`
	CheckFrequency   time.Duration `json:"checkFrequency"`
	URLs             []string      `json:"urls"`
}

var config Config
var parallelWorkers = 4
var pool chan Worker
var botMutex *sync.Mutex
var bot *tgbotapi.BotAPI

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
	msg := tgbotapi.NewMessage(config.NotifiableChatID, message)
	bot.Send(msg)
}
