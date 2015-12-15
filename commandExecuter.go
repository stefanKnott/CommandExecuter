//Written by Stefan Knott Fall 2015
/*Note: Error code -1000 is used in this program to denote an error opening a file.  -1000 was chosen
as to never potentially interfere with 8bit checksum information
*/

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

/*
Wrapper needed to pass command string array over channel cleanly
*/
type recWrap struct {
	input []string 
}

/*
Performs checksum on the word count of a file
I was fuzzy as what to do here without further clarification

Assumptions:
	--8 bits restricts to values 0-255
	--check sum found via inverting sum of data (wc) and adding one

Params: 
	filename: name of file to be used in wordCount call
Return:
	int: a check sum for the word count
	time.Duration: how long the checkSum operation took
*/
func checkSum(filename string)(int, time.Duration){
	startTime := time.Now()
	wc, _ := wordCount(filename)
	if wc == -1000{	//error opening file
		return -1000, time.Since(startTime)
	}
	
	//limit to 8bit representation
	wc %= 256

	//two's comp
	ch := ^wc + 1

	dur := time.Since(startTime)
	return ch, dur
}

/*
Counts number of words within a file
Uses strings.Fields to seperate words by whitespace..much faster than regexing \S+\s+ repeatedly and incrementing counter

Assumptions:
	--all hyphenated words coun't as one word
	--words are seperated by whitespace

Params:
	filename: name of file to perform word count upon
Return:
	int: word count
	time.Duration: how long the word count operation took 
*/
func wordCount(filename string)(int, time.Duration){
	startTime := time.Now()
	file, err := os.Open(filename)
	if err != nil {
		return -1000, time.Since(startTime)
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

/*
Checks the frequency of occurance of word within a file with a given filename
Uses strings.Fields to seperate words by whitespace..much faster than regexing \S+\s+ repeatedly and incrementing counter

Assumptions for search "report":
	--only "report", "report," or "report." count as an occurance of "report"
	--instances such as "Report" and "Report-..." would not count as an occurance as they most likely have different meanings than the desired word

Params:
	filename: name of file to open and measure word frequency on
	word: search for repetitions of this word
Return:
	int: number of occurances of word
	time.Duration: how long the word frequency operation took
*/
func wordFreq(filename string, word string)(int, time.Duration){
	startTime := time.Now()
	file, err := os.Open(filename)
	if err != nil {
		return -1000, time.Since(startTime)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wc := 0
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for _, lword := range words {
			if lword == word || lword == word+"," || lword == word+"."{
				wc += 1
			}
		}
	}

	dur := time.Since(startTime)
	return wc, dur
}

/*
Consume command requests, spin up 5 goroutines to read these requests via channel from producer

Assumptions:
	--variations of command input possible in spacing, capitalization variations allowed for command names
	--system limits not infinite (limit open files, and concurrent ops)

Params:
	recordChan: channel by which consumer reads command input data from cmdProducer goroutine
*/
func cmdConsumer(recordChan chan recWrap) {
	defer wg.Done()

	for i := 0; i < 5; i++ {
		go func() {
			for record := range recordChan {
				cmd := strings.TrimSpace(record.input[0])
				switch (strings.ToUpper(cmd)){
				case "CHECKSUM":
					cs, dur := checkSum(strings.TrimSpace(record.input[1]))
					if cs == -1000{
						fmt.Println("Invalid line:", record.input)
						continue
					}
					fmt.Printf("%s,%s, %d, %v\n", record.input[0], record.input[1], cs, dur)
				case "WORDCOUNT":
					totWc, dur := wordCount(strings.TrimSpace(record.input[1]))

					if totWc == -1000{
						fmt.Println("Invalid line:", record.input)
						continue
					}
					fmt.Printf("%s,%s, %d, %v\n", record.input[0], record.input[1], totWc, dur)
				case "WORDFREQ":
					wc, dur := wordFreq(strings.TrimSpace(record.input[1]), 
										strings.TrimSpace(record.input[2]))
					if wc == -1000{
						fmt.Println("Invalid line:", record.input)
						continue
					}
					fmt.Printf("%s,%s,%s, %d, %v\n", record.input[0], record.input[1], record.input[2], wc, dur)
				default:
					fmt.Println("Invalid line: ", record.input)
				}
			}
		}()
	}
}

/*
Produce command requests to be done using csv package, send to consumer over channel

Params:
	cmdFile: name of file which holds command list to run
	recordChan: channel over which command lines are read and passed to cmdConsumer goroutine
*/
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

/*
Spin up producer and consumers to handle work queries
*/
func main() {
	recordChan := make(chan recWrap)
	defer close(recordChan)
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
