package main

import (
	"testing"
)

func TestCheckURL(t *testing.T) {
	pool := make(chan Worker, 1)
	var curWorker Worker
	curWorker.CheckURL("http://somerandomurl.abc", pool)
S1:
	select {
	case <-pool:
		break S1
	default:
		t.Error("Worker wasn't added to pool after checking had finished")
	}

	curWorker.CheckURL("http://github.com/randomstring/randomstring", pool)
S2:
	select {
	case <-pool:
		break S2
	default:
		t.Error("Worker wasn't added to pool after checking had finished")
	}

	curWorker.CheckURL("http://google.com", pool)
S3:
	select {
	case <-pool:
		break S3
	default:
		t.Error("Worker wasn't added to pool after checking had finished")
	}
}
