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
	count := 3
	wg.Add(count)
	for i := 0; i < count; i++ {
		thread := i
		go func() {
			defer wg.Done()
			doWork(log, thread)
			doWork(log, thread)
		}()
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

func doWork(parent frog.Logger, n int) {
	log := frog.AddAnchor(parent)
	defer frog.RemoveAnchor(log)

	log.Transient(" + starting...", frog.Int("thread", n))
	time.Sleep(time.Duration(400*n) * time.Millisecond)

	for j := 0; j <= 100; j++ {
		log.Transient(" + Status", frog.Int("thread", n), frog.Int("percent", j))
		time.Sleep(time.Duration(5+rand.Intn(30)) * time.Millisecond)
	}
}
