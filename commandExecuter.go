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
	filename, word, count string
	dur time.Duration
}

func checkSum(filename string){
	//startTime := time.Now()
	defer wg.Done()
	fmt.Println("checksum")
	//dur := time.Since(startTime)
}

//NOTE: words connected by a hyphen are counted as one whole word
func wordCount(filename string){
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
	fmt.Printf("WORDCOUNT, %s, %d, %v\n",filename, wc, dur)
}

//Checks the frequency of occurance of word within a file with a given filename
func wordFreq(filename string , word string){
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
	fmt.Printf("WORDFREQ, %s, %s, %d, %v\n", filename, word, wc, dur)
}

//Produce work request to be done using csv, trimming white space from values
func cmdReader(filename string){
	defer wg.Done()
	cmds, _ := os.Open(filename)
	defer cmds.Close()

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
			go wordCount(arg1)
		case "WORDFREQ":
			wg.Add(1)
			arg2 := strings.TrimSpace(record[2])
			go wordFreq(arg1, arg2)
		default:
			fmt.Println("Invalid command: ", cmd)
		}
	}
}

func main(){
	filename := os.Args[1]
	wg.Add(1)
	go cmdReader(filename)
	wg.Wait()
}