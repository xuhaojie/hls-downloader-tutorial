package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	urlBase := "https://wolongzywcdn3.com:65/20220415/3f7cISA9/"
	blockNum := 100

	f, err := os.OpenFile("/tmp/video.ts", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	for i := 0; i < blockNum; i++ {
		url := urlBase + fmt.Sprintf("0%d", i) + ".ts"
		resp, err := http.Get(url)
		fmt.Printf("Get block %d form %s %d\n", i, url, resp.StatusCode)
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
					_, err := f.Write(body)
					if err != nil {
						fmt.Println(err.Error())
						return
					}
				}
			}
		}
	}
}
