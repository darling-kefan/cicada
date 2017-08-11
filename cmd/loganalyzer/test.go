package main

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	f, err := os.Create("profile.out")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer func() {
		pprof.StopCPUProfile()
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	for i := 0; i < 10000; i++ {
		doIt(i)
	}
}

func doIt(i int) {
	out := ""
	for j := 0; j < i; j++ {
		if i%2 == 0 {
			out += "X"
		}
		if j%5 == 0 {
			out += "Y"
		}
	}
	fmt.Println(i)
}
