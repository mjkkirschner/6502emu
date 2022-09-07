package main

import (
	"fmt"
	"os"
)

func (sim *Simulator) printMemoryAt(address uint8, width uint8) {
	start := address - width/2
	end := address + width/2
	for i := start; i <= end; i++ {
		fmt.Println(i, ":", sim.Memory[i])
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (sim *Simulator) loadMemoryFromBinary(filePath string) {
	data, err := os.ReadFile(filePath)
	check(err)
	fmt.Println(len(data))
	//TODO is this safe? Should we instead iterate?
	sim.Memory = data
}
