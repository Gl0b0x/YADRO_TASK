package main

import (
	"YADRO/pkg"
	"bufio"
	"fmt"
	"os"
)

func main() {
	var (
		computerClub *pkg.ComputerClub
		ok           bool
	)
	if len(os.Args) != 2 {
		fmt.Println("error: command line arguments")
		return
	}
	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error: file %s not found\n", filename)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	computerClub, ok = pkg.InitClub(scanner)
	if !ok {
		if scanner.Text() == "" {
			fmt.Println("error: not enough input data")
		} else {
			fmt.Println(scanner.Text())
		}
		return
	}
	ok = computerClub.ParseEvents(scanner, computerClub.CountComputers)
	if !ok {
		fmt.Println(scanner.Text())
		return
	}
	computerClub.DoWork()
}
