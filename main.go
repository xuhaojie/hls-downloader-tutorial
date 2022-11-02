package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	urlBase := "https://wolongzywcdn3.com:65/20220415/3f7cISA9/"
	blockNum := 100
	for i := 0; i < blockNum; i++ {
		url := urlBase + fmt.Sprintf("0%d", i) + ".ts"
		fmt.Println(url)
		resp, err := http.Get(url)
		fmt.Printf("Get block %d from %s %d\n", i, url, resp.StatusCode)
		if err != nil {
			fmt.Println(err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
					return
				} else {
					fmt.Println(body)
				}
			}
		}
	}
}
