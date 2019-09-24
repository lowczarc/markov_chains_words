package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const depth = 10

var randSeed = rand.NewSource(time.Now().UnixNano())
var randGen = rand.New(randSeed)

func chargeFilesStats(fileName string, markovProba map[string]int) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	var chain [depth]byte

	for i := 0; i < depth; i++ {
		chain[i] = '\n'
	}

	n, err := f.Read(chain[depth-1:])
	for ; err == nil && n == 1; n, err = f.Read(chain[depth-1:]) {
		markovProba[string(chain[0:])]++
			for i := 1; i < depth; i++ {
				chain[i-1] = chain[i]
			}
	}

	return nil
}

func nextLetter(markovProba map[string]int, precLetters [depth - 1]byte) byte {
	var max int
	var tmp [depth]byte

	for i := 0; i < depth-1; i++ {
		tmp[i] = precLetters[i]
	}

	for i := 0; i <= 255; i++ {
		tmp[depth-1] = byte(i)
		max += markovProba[string(tmp[0:])]
	}

	if max == 0 {
		return 0
	}
	r := randGen.Intn(max) + 1

	for i := 0; i <= 255; i++ {
		tmp[depth-1] = byte(i)
		r -= markovProba[string(tmp[0:])]
		if r <= 0 && markovProba[string(tmp[0:])] != 0 {
			return byte(i)
		}
	}

	return 0
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wordgen file_names ...")
		return
	}

	markovProba := make(map[string]int)

	for i := 1; i < len(os.Args); i++ {
		err := chargeFilesStats(os.Args[i], markovProba)
		if err != nil {
			panic(err)
		}
	}

	var letter [depth - 1]byte

	for i := 0; i < depth-1; i++ {
		letter[i] = '\n'
	}
	for {
		newLetter := nextLetter(markovProba, letter)
		if newLetter == 0 {
			break;
		}
		for i := 1; i < depth-1; i++ {
			letter[i-1] = letter[i]
		}

		letter[depth-2] = newLetter
		fmt.Printf("%c", newLetter)
	}
	fmt.Println("\n")
}
