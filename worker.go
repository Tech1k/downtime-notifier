package main

import (
	"fmt"
	"net/http"
)

type Worker struct {
	URLChecker
}

type URLChecker interface {
	CheckURL()
}

/*
	CheckURL tries to get content of specified url.
	If http code is not 200, it sends notification with
	telegram bot
*/
func (w Worker) CheckURL(url string, pool chan<- Worker) {
	defer func() {
		pool <- w
	}()
	resp, err := http.Get(url)
	if err != nil {
		msg := fmt.Sprintf("URL %s returns error:\n%s", url, err.Error())
		go SendMsg(msg)
		return
	}
	if resp.StatusCode != 200 {
		msg := fmt.Sprintf("URL %s returns http code %d", url, resp.StatusCode)
		go SendMsg(msg)
	}
}
