package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Block struct {
	index int
	url   string
	data  []byte
}

var wg = sync.WaitGroup{}
var connChan chan int

func work(b *Block) error {
	defer wg.Done()
	conn := <-connChan
	b.data = nil
	resp, err := http.Get(b.url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("[%02d] get %s %d\n", conn, b.url, resp.StatusCode)
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		b.data = body
	}
	connChan <- conn
	return err
}

func main() {
	beginTime := time.Now()
	var urlBase = "https://wolongzywcdn3.com:65/20220415/3f7cISA9/"
	const blockNum = 100 //3335
	blocks := make([]Block, blockNum)
	const maxConnections = 32
	connChan = make(chan int, maxConnections)
	for i := 0; i < maxConnections; i++ {
		connChan <- i
	}

	wg.Add(blockNum)
	for i := 0; i < blockNum; i++ {
		url := urlBase + fmt.Sprintf("0%d", i) + ".ts"
		blocks[i].index = i
		blocks[i].url = url
		go work(&blocks[i])
	}
	wg.Wait()

	file := "/tmp/test.ts"
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Writing data....")
	for i := 0; i < blockNum; i++ {
		if blocks[i].data != nil {
			_, err = f.Write(blocks[i].data)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	endTime := time.Now()
	fmt.Println(endTime.Sub(beginTime))
	fmt.Println("Done.")
}
