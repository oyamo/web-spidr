package downloader

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hash/fnv"
	"io"
	"net/http"
	"regexp"
	"time"
)

// Reduces io time and stack usage
type rawResponse struct {
	body  *string
	title *string
	url   *string
	err   *error
	hash  *uint32
}
type Response struct {
	ID           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	Url          string             `bson:"url"`
	Hash         uint32             `bson:"hash"`
	DateModified int64              `bson:"modified"`
}


func createPageHash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32() - 10
}

//download Downloads the Page from the link in the link.txt file
func download(i int, jobs <-chan *string, results chan<- *rawResponse) {
	for url := range jobs {
		fmt.Println("Job", i, "downloading", *url)
		time.Sleep(time.Second)
		res, err := http.Get(*url)

		if err != nil {

			results <- &rawResponse{
				err: &err,
			}
		} else {
			data, _ := io.ReadAll(res.Body)
			body := string(data)
			hash := createPageHash(body)

			title := regexp.MustCompile("<title>([\\d\\w\\D\\W\\S\\s]*)</title>")
			titleStr := title.FindStringSubmatch(string(data))
			if len(titleStr) >= 2 {
				titleStr = titleStr[1:]
			} else {
				titleStr = append(titleStr, "Empty")
			}

			results <- &rawResponse{
				body:  &body,
				title: &titleStr[0],
				err:   nil,
				url:   url,
				hash:  &hash,
			}
		}
		fmt.Println("Job", i, "completed")
	}

}
