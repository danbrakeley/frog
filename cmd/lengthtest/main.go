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
	log.Info("Spawning threads...", frog.Int("count", threads))
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		n := i
		anchored := frog.AddAnchor(log)
		dark := frog.WithOptions(anchored, frog.POPalette(frog.Palette{
			{frog.DarkCyan, frog.DarkGray}, // Transient
			{frog.DarkCyan, frog.DarkGray}, // Verbose
			{frog.DarkCyan, frog.DarkGray}, // Info
			{frog.DarkCyan, frog.DarkGray}, // Warning
			{frog.DarkCyan, frog.DarkGray}, // Error
		}))
		go func() {
			runProcess(dark, n)
			dark.Info("thread finished", frog.Int("thread", n))
			frog.RemoveAnchor(anchored)
			wg.Done()
		}()
	}

	sleep(800)
	log.Info("main thread reporting in with a really, really long line that just goes on and on and on and on and never seems to end because it keeps going and going and going and oh wow it is still going right up until it stops")
	sleep(400)
	log.Warning("main thread again, this time with a warning line that is also unusually long and full of clauses that just keep appending themselves without really ever stopping for maybe a punctuation mark or whatever but no this line just keeps going and goin and going and did you notice I missed a g back there crazy right where is my head at")
	sleep(500)
	log.Info("more from main thread, but this line is a lot shorter")

	wg.Wait()
	log.Info("done!")
}

func runProcess(log frog.Logger, n int) {
	log.Info("thread started", frog.Int("thread", n))
	for j := 0; j <= 100; j++ {
		log.Transient(
			" + Status that is really, really long; so long that it will probably fall off the end of the terminal edge... which is exactly what I want to check...",
			frog.Int("thread", n),
			frog.Int("percent", j),
		)
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	}
}
