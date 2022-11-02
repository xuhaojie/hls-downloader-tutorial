package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	urlBase := "https://wolongzywcdn3.com:65/20220415/3f7cISA9/"
	blockNum := 3336

	f, err := os.OpenFile("/tmp/video.ts", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("Create file failed ", err)
		return
	}

	defer f.Close()

	const MaxRetry = 3
	for i := 3300; i < blockNum; i++ {
		url := urlBase + fmt.Sprintf("0%d", i) + ".ts"
		retry := 0
		for retry < MaxRetry {
			retry++
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Get block %d from %s %d\n", i, url, resp.StatusCode)
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
					break
				}
			}
		}
		if retry >= MaxRetry {
			fmt.Printf("Get block %d from %s failed for %d times.\n", i, url, retry)
			break
		}
	}
}
