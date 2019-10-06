package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/danbrakeley/frog"
)

func main() {
	log := frog.New()
	defer log.Close()

	log.SetMinLevel(frog.Progress)
	log.Progressf("progress line")
	log.Verbosef("verbose line")
	log.Infof("info line")
	log.Warningf("warning line")
	log.Errorf("error line")

	threads := 5
	log.Infof("Spawning %d threads...", threads)
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		n := i
		fl := log.AddFixedLine()
		go func() {
			runProcess(fl, n)
			fl.Close()
			wg.Done()
		}()
	}

	time.Sleep(time.Second)
	log.Infof("still running...")
	time.Sleep(time.Duration(500) * time.Millisecond)
	log.Infof("yup, still running...")
	time.Sleep(time.Duration(100) * time.Millisecond)
	log.Warningf("something happened on the main thread")
	time.Sleep(time.Duration(500) * time.Millisecond)
	log.Infof("the main thread again")
	time.Sleep(time.Duration(5000) * time.Millisecond)
	log.Errorf("the main thread had an error?")

	wg.Wait()
	log.Infof("done!")
}

func runProcess(log frog.Logger, n int) {
	log.Progressf("starting...")
	time.Sleep(time.Duration(400*n) * time.Millisecond)
	for j := 0; j <= 100; j++ {
		log.Progressf("Loop %d Status: %0d%%", n, j)
		if j == 90 {
			log.Verbosef("Loop %d download complete, opening file for write", n)
		} else if j == 100 {
			log.Infof("Download of %d complete", n)
		}
		time.Sleep(time.Duration(50-(10*n)+rand.Intn(50)) * time.Millisecond)

		if j == 50 && rand.Intn(3) == 0 {
			log.Warningf("Loop %d encountered a problem at 50%%, retrying", n)
			time.Sleep(time.Duration(n) * time.Second)
		}
	}
}
