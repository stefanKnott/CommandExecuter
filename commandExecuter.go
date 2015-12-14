package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type output struct {
	filename, word string
	count          int
	dur            time.Duration
}

type recWrap struct {
	input []string //command input
}

func checkSum(filename string)(int, time.Duration){
	startTime := time.Now()
	dur := time.Since(startTime)
	return 0, dur
}

//NOTE: words connected by a hyphen are counted as one whole word
func wordCount(filename string)(int, time.Duration){
	startTime := time.Now()
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	totWc := 0
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		totWc += len(words)
	}
	dur := time.Since(startTime)
	return totWc, dur
}

//Checks the frequency of occurance of word within a file with a given filename
func wordFreq(filename string, word string)(int, time.Duration){
	startTime := time.Now()
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wc := 0
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for _, lword := range words {
			if lword == word || lword == word+"," {
				wc += 1
			}
		}
	}

	dur := time.Since(startTime)
	return wc, dur
}

//Consume command requests, spin up 5 goroutines to read these requests via recordChan
func cmdConsumer(recordChan chan recWrap) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		go func() {
			for record := range recordChan {
				cmd := strings.TrimSpace(record.input[0])
				arg1 := strings.TrimSpace(record.input[1])
				switch (strings.ToUpper(cmd)) {
				case "CHECKSUM":
					cs, dur := checkSum(arg1)
					fmt.Printf("%s,%s, %d, %v\n", record.input[0], record.input[1], cs, dur)
				case "WORDCOUNT":
					totWc, dur := wordCount(arg1)
					fmt.Printf("%s,%s, %d, %v\n", record.input[0], record.input[1], totWc, dur)
				case "WORDFREQ":
					arg2 := strings.TrimSpace(record.input[2])
					wc, dur := wordFreq(arg1, arg2)
					fmt.Printf("%s,%s, %d, %v\n", record.input[0], record.input[1], wc, dur)
				default:
					fmt.Println("Invalid line: ", record.input)
				}
			}
		}()
	}
}

//Produce command requests to be done using csv package
func cmdProducer(cmdFile string, recordChan chan recWrap) {
	defer wg.Done()
	cmds, err := os.Open(cmdFile)
	if err != nil{
		log.Fatal(err)
	}
	defer cmds.Close()

	r := csv.NewReader(bufio.NewReader(cmds))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		recordChan <- recWrap{input: record}
	}
}

func main() {
	recordChan := make(chan recWrap)
	if len(os.Args) != 2{
		fmt.Println("Use format: ./commandExecuter <command_file.txt>")
		log.Fatal()
	}
	cmdFile := os.Args[1]
	wg.Add(2)
	go cmdProducer(cmdFile, recordChan)
	go cmdConsumer(recordChan)
	wg.Wait()
}
