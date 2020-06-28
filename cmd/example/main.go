package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/danbrakeley/frog"
)

var verbose = flag.Bool("verbose", false, "drop min level from info to verbose")
var json = flag.Bool("json", false, "output structured JSON")
var noTime = flag.Bool("notime", false, "do not include timestamps (ignored if using -json)")
var noLevel = flag.Bool("nolevel", false, "do not include level (ignored if using -json)")

func main() {
	flag.Parse()

	style := frog.Auto
	if *json {
		style = frog.JSON
	}

	var opts []frog.Option
	if *noTime {
		opts = append(opts, frog.HideTimestamps)
	}
	if *noLevel {
		opts = append(opts, frog.HideLevel)
	}

	log := frog.New(style, opts...)
	defer log.Close()

	log.Info("Frog Example App")
	flag.VisitAll(func(f *flag.Flag) {
		log.Info(fmt.Sprintf("  --%s :: %s", f.Name, f.Usage))
	})
	var fields []frog.Fielder
	for i, v := range os.Args {
		arg := v
		if i == 0 {
			arg = filepath.Base(arg)
		}
		fields = append(fields, frog.String(fmt.Sprintf("arg%d", i), arg))
	}
	log.Info("os.Args", fields...)

	log.SetMinLevel(frog.Transient)
	log.Transient("transient line")
	log.Verbose("verbose line")
	log.Info("info line")
	log.Warning("warning line")
	log.Error("error line")

	if *verbose {
		log.SetMinLevel(frog.Verbose)
	} else {
		log.SetMinLevel(frog.Info)
	}

	threads := 5
	log.Info("Spawning threads...", frog.Int("count", threads))
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		n := i
		fl := frog.AddFixedLine(log)
		go func() {
			fl.Verbose("thread spawned", frog.Int("thread", n))
			runProcess(fl, n)
			fl.Verbose("thread closing", frog.Int("thread", n))
			frog.RemoveFixedLine(fl)
			wg.Done()
		}()
	}

	time.Sleep(time.Second)
	log.Info("still running...")
	time.Sleep(time.Duration(500) * time.Millisecond)
	log.Info("yup, still running...")
	time.Sleep(time.Duration(100) * time.Millisecond)
	log.Warning("something happened on the main thread")
	time.Sleep(time.Duration(500) * time.Millisecond)
	log.Info("the main thread again")
	time.Sleep(time.Duration(5000) * time.Millisecond)
	log.Error("the main thread had an error?")

	wg.Wait()
	log.Info("done!")
}

func runProcess(log frog.Logger, n int) {
	log.Transient(" + starting...", frog.Int("thread", n))
	time.Sleep(time.Duration(400*n) * time.Millisecond)
	for j := 0; j <= 100; j++ {
		if j == 90 {
			log.Verbose("transitioning from downloading to writing", frog.Int("thread", n))
		} else if j == 100 {
			log.Info("finished downloading", frog.Int("thread", n))
		}
		log.Transient(" + Status", frog.Int("thread", n), frog.Int("percent", j))
		time.Sleep(time.Duration(50-(10*n)+rand.Intn(50)) * time.Millisecond)

		if j == 50 && rand.Intn(3) == 0 {
			log.Warning("encountered a problem, retrying", frog.Int("thread", n), frog.Int("percent", 50))
			time.Sleep(time.Duration(n+1) * time.Second)
		}
	}
}
