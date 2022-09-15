package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rizky-putra19/go-queue-system/queue"
)

func main() {
	queue.Init("amqp://localhost:15672/#/")

	if os.Args[1] == "worker" {
		worker()
	} else {
		publisher()
	}
}

func publisher() {
	for {
		if err := queue.Publish("add_q", []byte("1,1")); err != nil {
			panic(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func worker() {
	msgs, close, err := queue.Subcribe("add_q")
	if err != nil {
		panic(err)
	}
	defer close()
	stop := make(chan bool)
	go func() {
		for d := range msgs {
			i1, i2 := toNums(d.Body)
			fmt.Println(time.Now().Format("01-02-2006 15:04:05"), "::", i1+i2)
			d.Ack(false)
		}
	}()
	log.Printf("To exit press CTRL+C")
	<-stop
}

func toNums(b []byte) (int, int) {
	s := string(b)
	ss := strings.Split(s, ",")
	i1, _ := strconv.Atoi(ss[0])
	i2, _ := strconv.Atoi(ss[1])
	return i1, i2
}
