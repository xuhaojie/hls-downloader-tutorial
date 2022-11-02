package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Block struct {
	index int
	url   string
	data  []byte
}

func work(b *Block) error {
	b.data = nil
	resp, err := http.Get(b.url)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("get %s %d\n", b.url, resp.StatusCode)
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		b.data = body
	}
	return err
}

func main() {
	beginTime := time.Now()

	urlBase := "https://wolongzywcdn3.com:65/20220415/3f7cISA9/"
	const blockNum = 3335
	blocks := make([]Block, blockNum)

	for i := 0; i < blockNum; i++ {
		url := urlBase + fmt.Sprintf("0%d", i) + ".ts"
		blocks[i].index = i
		blocks[i].url = url
		go work(&blocks[i])
	}

	endTime := time.Now()
	fmt.Println(endTime.Sub(beginTime))
	fmt.Println("Done.")
}
