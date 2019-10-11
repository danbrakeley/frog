package main

import (
	"flag"
	"math/rand"
	"sync"
	"time"

	"github.com/danbrakeley/frog"
)

func sleep(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func main() {
	flag.Parse()

	log := frog.New(frog.Auto)
	defer log.Close()

	threads := 3
	log.Infof("Spawning %d threads...", threads)
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		n := i
		fixed := frog.AddFixedLine(log)
		go func() {
			runProcess(fixed, n)
			fixed.Infof("thread %d finished", n)
			frog.RemoveFixedLine(fixed)
			wg.Done()
		}()
	}

	sleep(800)
	log.Infof("main thread reporting in")
	sleep(400)
	log.Warningf("main thread warning")
	sleep(500)
	log.Infof("main thread reporting in")
	sleep(1000)
	log.Infof("main thread reporting in")

	wg.Wait()
	log.Infof("done!")
}

func runProcess(log frog.Logger, n int) {
	for j := 0; j <= 100; j++ {
		log.Transientf(" + [%d] Status: %0d%%", n, j)
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	}
}
