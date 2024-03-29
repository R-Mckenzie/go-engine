package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/R-Mckenzie/go-engine/engine"
)

const width, height = 32.0 * 30, 32.0 * 20

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	// Profiling -- only happens if we supply the data
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ==================================
	// Code for actually running the game
	game := engine.CreateGame(width, height)
	s2 := newScene2(game)
	s := newScene(game, s2)
	game.SetScene(s)
	game.Run()

	defer game.Quit()
	// ==================================

	// profiling
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

}
