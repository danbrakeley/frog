package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/danbrakeley/frog"
)

func main() {
	log := frog.New(frog.Auto)
	defer log.Close()

	log.Info("Spawning example threads...", frog.Int("count", 3))
	time.Sleep(time.Second)

	wg := new(sync.WaitGroup)
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go doWork(wg, frog.AddAnchor(log), i)
	}

	time.Sleep(time.Second)
	log.Info("waited for one second...")
	time.Sleep(time.Second)
	log.Warning("waited for two seconds...")
	time.Sleep(time.Second)
	log.Error("BORED OF WAITING")
	wg.Wait()

	log.Info("All threads done!")
}

func doWork(wg *sync.WaitGroup, log frog.Logger, n int) {
	defer wg.Done()
	defer frog.RemoveAnchor(log)

	log.Transient(" + starting...", frog.Int("thread", n))
	time.Sleep(time.Duration(400*n) * time.Millisecond)

	for j := 0; j <= 100; j++ {
		log.Transient(" + Status", frog.Int("thread", n), frog.Int("percent", j))
		time.Sleep(time.Duration(10+rand.Intn(50)) * time.Millisecond)
	}
}
