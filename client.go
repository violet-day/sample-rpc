package main

import (
	"github.com/gjvnq/go-logger"
	"github.com/violet-day/simple-rpc/rpc"
	"os"
	"strconv"
	"sync"
	"time"
)

var log, _ = logger.New("client", 1, os.Stdout)

func main() {
	var wg sync.WaitGroup

	c := rpc.NewClient()

	time.Sleep(time.Second * 1)

	c.Connect(":8080")
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(v int) {
			resp, err := c.Do(strconv.Itoa(v))

			if err != nil {
				log.Error(err.Error())
			} else {
				log.InfoF("%d double is %s", v, resp)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
