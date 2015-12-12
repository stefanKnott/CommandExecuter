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

type cmdWrap struct{
	cmd []string
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
	totWc := 0
	wc := 0
	for scanner.Scan(){
		words := strings.Fields(scanner.Text())
		totWc += len(words)
		for _, lword := range words{
			if lword == word || lword == word +","{
				wc += 1
			}
		}
	}

	dur := time.Since(startTime)
	fmt.Printf("WORDFREQ, %s, %s, %.3f, %v\n", filename, word, float64(wc)/float64(totWc), dur)
}

//Consume work to be done, spinning up new goroutine for each command
func consume(cmdChan chan cmdWrap){
	defer wg.Done()
	for res := range recChan{
		fmt.Println(res)
		switch(res.cmd[0]){
		case "CHECKSUM":
			wg.Add(1)
			//go checkSum(record[1])
		case "WORDCOUNT":
			wg.Add(1)
			go wordCount(res.cmd[1])
		case "WORDFREQ":
			wg.Add(1)
			go wordFreq(res.cmd[1], res.cmd[2])
		default:
			break
		}
	}
}

//Produce work to be done using csv package to read comma seperated command values in commandl=_file.txt
func produce(cmdChan chan cmdWrap){
	defer wg.Done()
	cmds, _ := os.Open("command_file.txt")
	defer cmds.Close()

	r := csv.NewReader(bufio.NewReader(cmds))
	for{
		record, err := r.Read()
		if err == io.EOF{
			break
		}

		cmdWrapped := cmdWrap{cmd:record}
		cmdChan <- cmdWrapped
	}
}

func main(){
	cmdChan := make(chan cmdWrap)
	wg.Add(2)
	go produce(cmdChan)
	go consume(cmdChan)
	wg.Wait()
}