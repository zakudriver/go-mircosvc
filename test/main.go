package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	// ch := make(chan int, 1)
	fmt.Println("start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	time.Sleep(time.Second * 5) // 假装 5 秒没准备好接收
	s := <-c
	fmt.Println("Signal")
	log.Println(s)
}
