package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"mcgaunn.com/logwild/pkg/cmd"
)

func main() {
	cpuProfile := os.Getenv("LOGWILD_CPU_PROFILE")
	memProfile := os.Getenv("LOGWILD_MEM_PROFILE")

	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal("could not start cpu profile and it was requested: ", err)
		}
		defer f.Close() // this can throw an error but whatev
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	if err := cmd.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
	// after running, if requested, generate a memory profile
	if memProfile != "" {
		f, err := os.Create(memProfile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // this can also throw error
		runtime.GC()    // ask runtime to perform garbage collection
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile :(): ", err)
		}
	}
}
