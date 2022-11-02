package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const MaxRetry = 3

type Block struct {
	index int
	conn  int
	url   string
	retry int
	err   error
	data  []byte
}

type Task struct {
	urlBase        string
	blockNum       int
	maxConnections int
	file           string
	connChan       chan (int)
	blocks         []*Block
	saveChan       chan (*Block)
	endChan        chan (error)
	finishedBlocks int
	client         *http.Client
}

func NewTask(urlBase string, blockNum int, file string, maxConnections int) *Task {
	t := new(Task)
	t.urlBase = urlBase
	t.blockNum = blockNum
	t.maxConnections = maxConnections
	t.blocks = make([]*Block, blockNum)
	t.saveChan = make(chan (*Block), blockNum)
	t.connChan = make(chan int, maxConnections)
	t.endChan = make(chan error, 1)
	t.file = file
	t.client = &http.Client{
		Timeout: 30 * time.Second,
	}
	return t
}

func (t *Task) fetch(b *Block) {
	b.data = nil
	resp, err := t.client.Get(b.url)

	if err != nil {
		b.err = err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			body, err := ioutil.ReadAll(resp.Body)
			b.err = err
			b.data = body
		}
	}
	t.saveChan <- b
}

func (t *Task) Run() error {
	f, err := os.OpenFile(t.file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	for i := 0; i < t.maxConnections; i++ {
		t.connChan <- i
	}

	fetchIndex := 0
	saveIndex := 0
	totalSize := 0
	for {
		select {
		case c := <-t.connChan:
			if fetchIndex < t.blockNum {
				url := t.urlBase + fmt.Sprintf("0%d", fetchIndex) + ".ts"
				b := new(Block)
				b.conn = c
				b.err = nil
				b.index = fetchIndex
				b.url = url
				b.retry = 1
				fetchIndex++
				go t.fetch(b)
			}

		case b := <-t.saveChan:
			if b.err != nil {
				fmt.Println(b.err)
				b.retry++
				if b.retry <= MaxRetry {
					fmt.Printf("Block %d retry %d/3\n", b.index, b.retry)
					b.err = nil
					go t.fetch(b)
				} else {
					t.endChan <- errors.New("get block failed for 3 times")
				}
			} else {
				fmt.Printf("Block %d got %d bytes.\n", b.index, len(b.data))
				t.connChan <- b.conn
				t.blocks[b.index] = b
				for i := saveIndex; i < t.blockNum; i++ {
					b = t.blocks[i]
					if b != nil {
						if b.index == saveIndex {
							if b.data != nil {
								blockSize := len(b.data)
								fmt.Printf("Save block index %d size=%d\n", i, blockSize)
								_, err := f.Write(b.data)
								if err != nil {
									fmt.Println(err.Error())
									return err
								} else {
									totalSize += blockSize
								}
							}
							t.blocks[i] = nil
							saveIndex++
						}
					}
				}
				if saveIndex >= t.blockNum {
					t.endChan <- nil
				}
			}

		case err := <-t.endChan:
			return err
		}
	}
	return nil
}

func main() {

	beginTime := time.Now()

	task := NewTask("https://wolongzywcdn3.com:65/20220415/3f7cISA9/", 100, "/tmp/test1.ts", 16)
	err := task.Run()

	endTime := time.Now()
	elapsedTime := endTime.Sub(beginTime).Seconds()
	if err != nil {
		fmt.Printf("Task failed! %s", err)
	} else {
		fmt.Printf("Task finished in %0.2f second(s).\n", elapsedTime)
	}
	fmt.Println("Done.")
}
