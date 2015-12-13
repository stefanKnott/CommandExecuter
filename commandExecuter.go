package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"encoding/csv"
	"io"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type output struct{
	filename, word string 
	count int
	dur time.Duration
}

func checkSum(filename string){
	//startTime := time.Now()
	defer wg.Done()
	fmt.Println("checksum")
	//dur := time.Since(startTime)
}

//NOTE: words connected by a hyphen are counted as one whole word
func wordCount(filename string, wcChan chan output){
	startTime := time.Now()
	defer wg.Done()
	file, err := os.Open(filename)
	if err != nil{
		log.Fatal(err)
	}	
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wc := 0
	for scanner.Scan(){
		words := strings.Fields(scanner.Text())
		wc += len(words)
	}
	dur := time.Since(startTime)
	wcChan <- output{filename:filename, count:wc, dur:dur, word:""}
	//fmt.Printf("WORDCOUNT, %s, %d, %v\n",filename, wc, dur)
}

//Checks the frequency of occurance of word within a file with a given filename
func wordFreq(filename string , word string, wfChan chan output){
	startTime := time.Now()
	defer wg.Done()
	file, err := os.Open(filename)
	if err != nil{
		log.Fatal(err)
	}	
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wc := 0
	for scanner.Scan(){
		words := strings.Fields(scanner.Text())
		for _, lword := range words{
			if lword == word || lword == word +","{
				wc += 1
			}
		}
	}

	dur := time.Since(startTime)
	wfChan <- output{filename:filename, word:word, count:wc, dur:dur}
	//fmt.Printf("WORDFREQ, %s, %s, %d, %v\n", filename, word, wc, dur)
}

//Spin up routines collecting output data from given channels, close channels when done extracting data
func collector(wcChan chan output, wfChan chan output, csChan chan output){
	go func(){
		defer close(wcChan)
		for res := range wcChan{
			fmt.Printf("WORDCOUNT, %s, %d, %v\n", res.filename, res.count, res.dur)
		}
	}()

	go func(){
		defer close(wfChan)
		for res := range wfChan{
			fmt.Printf("WORDFREQ, %s, %s, %d, %v\n", res.filename, res.word, res.count, res.dur)
		}
	}()

	/*go func(){
		defer close(csChan)
		for res := range csChan{
			fmt.Println(res)
			//fmt.Println("CHECKSUM, %s, %s, %d, %v\n", res.filename, res.word, res.count, res.dur)
		}
	}()*/
	csChan = nil
}

//Produce work request to be done using csv, trimming white space from values
func cmdReader(cmdFile string){
	defer wg.Done()
	cmds, _ := os.Open(cmdFile)
	defer cmds.Close()

	wcChan := make(chan output)
	wfChan := make(chan output)
	csChan := make(chan output)

	go collector(wcChan, wfChan, csChan)

	r := csv.NewReader(bufio.NewReader(cmds))
	for{
		record, err := r.Read()
		if err == io.EOF{
			break
		}
		cmd := strings.TrimSpace(record[0])
		arg1 := strings.TrimSpace(record[1])
		switch(cmd){
		case "CHECKSUM":
			wg.Add(1)
			go checkSum(arg1)
		case "WORDCOUNT":
			wg.Add(1)
			go wordCount(arg1, wcChan)
		case "WORDFREQ":
			wg.Add(1)
			arg2 := strings.TrimSpace(record[2])
			go wordFreq(arg1, arg2, wfChan)
		default:
			fmt.Println("Invalid command: ", cmd)
		}
	}
}

func main(){
	cmdFile := os.Args[1]
	wg.Add(1)
	go cmdReader(cmdFile)
	wg.Wait()
}