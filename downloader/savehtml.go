package downloader

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
)

func dumpHtml(wg *sync.WaitGroup, data rawResponse)  {
	defer wg.Done()
	var fileName =	*data.title + strconv.Itoa(int(*data.hash)) + ".html"
	fileName = "./" + folder + "/" + fileName
	fileName  = string(regexp.MustCompile("[-\\s:#]").ReplaceAll([]byte(fileName), []byte("_")))
	_, err := os.ReadDir(folder)
	if err != nil {
		_ = os.Mkdir(folder, 0755)
	}

	fmt.Println(*data.url)
	fmt.Println(*data.title)
	fmt.Println()
	err = os.WriteFile(fileName, []byte(*data.body), 0666)
	if err != nil {
		panic(err)
	}
}

