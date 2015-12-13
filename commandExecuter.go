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

func checkSum(filename string) {
	//startTime := time.Now()
	//	defer wg.Done()
	fmt.Println("checksum")
	//dur := time.Since(startTime)
}

//NOTE: words connected by a hyphen are counted as one whole word
func wordCount(filename string, wcChan chan output) {
	startTime := time.Now()
	//	defer wg.Done()
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wc := 0
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		wc += len(words)
	}
	dur := time.Since(startTime)
	wcChan <- output{filename: filename, count: wc, dur: dur, word: ""}
}

//Checks the frequency of occurance of word within a file with a given filename
func wordFreq(filename string, word string, wfChan chan output) {
	startTime := time.Now()
	//defer wg.Done()
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
	wfChan <- output{filename: filename, word: word, count: wc, dur: dur}
}

//Spin up routines collecting output data from given channels, close channels when done extracting data
func outputCollector(wcChan chan output, wfChan chan output, csChan chan output) {
	go func(wcChan chan output) {
		for res := range wcChan {
			fmt.Printf("WORDCOUNT, %s, %d, %v\n", res.filename, res.count, res.dur)
		}
	}(wcChan)

	go func(wfChan chan output) {
		for res := range wfChan {
			fmt.Printf("WORDFREQ, %s, %s, %d, %v\n", res.filename, res.word, res.count, res.dur)
		}
	}(wfChan)

	/*go func(){
		defer close(csChan)
		for res := range csChan{
			fmt.Println(res)
			//fmt.Println("CHECKSUM, %s, %s, %d, %v\n", res.filename, res.word, res.count, res.dur)
		}
	}()*/
	csChan = nil
}

//Consume command requests, spin up 5 goroutines to read these requests via recordChan
func cmdConsumer(recordChan chan recWrap) {
	var lwg sync.WaitGroup //local waitgroup needed to ensure channels get closed properly

	wcChan := make(chan output)
	wfChan := make(chan output)
	csChan := make(chan output)
	defer close(wcChan)
	defer close(wfChan)
	defer close(csChan)

	go outputCollector(wcChan, wfChan, csChan)

	for i := 0; i < 5; i++ {

		lwg.Add(1)
		go func() {
			for record := range recordChan {

				cmd := strings.TrimSpace(record.input[0])
				filename := strings.TrimSpace(record.input[1])

				switch cmd {
				case "CHECKSUM":
					checkSum(filename)
				case "WORDCOUNT":
					wordCount(filename, wcChan)
				case "WORDFREQ":
					word := strings.TrimSpace(record.input[2])
					wordFreq(filename, word, wfChan)
				default:
					fmt.Println("Invalid command: ", cmd)
				}
			}
			lwg.Done()

		}()
	}
	lwg.Wait()
}

//Produce command requests to be done using csv package
func cmdProducer(cmdFile string, recordChan chan recWrap) {
	defer wg.Done()
	cmds, _ := os.Open(cmdFile)
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
	cmdFile := os.Args[1]
	wg.Add(1)
	go cmdProducer(cmdFile, recordChan)
	go cmdConsumer(recordChan)
	wg.Wait()
}
