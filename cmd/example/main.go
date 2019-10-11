package main

import (
	"flag"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/danbrakeley/frog"
)

var verbose = flag.Bool("verbose", false, "drop min level from info to verbose")
var json = flag.Bool("json", false, "drop min level from info to verbose")

func main() {
	flag.Parse()

	style := frog.Auto
	if *json {
		style = frog.JSON
	}
	log := frog.New(style)
	defer log.Close()

	log.Infof("Frog Example App")
	log.Infof("Flags:")
	log.Infof("  --verbose   : enable Verbose level logging")
	log.Infof("  --json      : output structured JSON")
	log.Infof("os.Args:")
	log.Infof("  %v", os.Args)

	log.SetMinLevel(frog.Transient)
	log.Transientf("transient line")
	log.Verbosef("verbose line")
	log.Infof("info line")
	log.Warningf("warning line")
	log.Errorf("error line")

	if *verbose {
		log.SetMinLevel(frog.Verbose)
	} else {
		log.SetMinLevel(frog.Info)
	}

	threads := 5
	log.Infof("Spawning %d threads...", threads)
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		n := i
		fl := frog.AddFixedLine(log)
		go func() {
			fl.Verbosef("spawned thread %d", n)
			runProcess(fl, n)
			fl.Verbosef("closing thread %d", n)
			frog.RemoveFixedLine(fl)
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
	log.Transientf(" + [%d] starting...", n)
	time.Sleep(time.Duration(400*n) * time.Millisecond)
	for j := 0; j <= 100; j++ {
		if j == 90 {
			log.Verbosef("thread %d transitioning from downloading to writing", n)
		} else if j == 100 {
			log.Infof("thread %d finished downloading", n)
		}
		log.Transientf(" + [%d] Status: %0d%%", n, j)
		time.Sleep(time.Duration(50-(10*n)+rand.Intn(50)) * time.Millisecond)

		if j == 50 && rand.Intn(3) == 0 {
			log.Warningf("thread %d encountered a problem at 50%%, retrying", n)
			time.Sleep(time.Duration(n+1) * time.Second)
		}
	}
}
