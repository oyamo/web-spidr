package downloader

import (
	"bufio"
	"bytes"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	folder = "kenya_law_org"
)

const (
	gmailAppPassword = ""
	smtpHost         = ""
	smtpPort         = 587
	smtpEmail        = ""
	senderName       = "Spidr Systems"
	toEmail          = ""
)

type Config struct {
	UrlTXTPath string
	Database   *mongo.Database
}

type Downloader struct {
	Config
	NumCores int
	client   http.Client
}

type scan struct {
	wg           *sync.WaitGroup
	db           *mongo.Database
	res          rawResponse
	docsFrmNet   *[]rawResponse
	maxLen       int
	counter      *int
	changedPages *[]Response
	docsFrmDB    *[]Response
	filter       *[]bson.M
}

func New(config Config) *Downloader {
	d := Downloader{}
	d.Database = config.Database
	d.UrlTXTPath = config.UrlTXTPath
	d.client = http.Client{
		Timeout: time.Minute,
	}

	return &d
}

type dFUnctions interface {
	Download()
}

func (d *Downloader) Download() (*scan, error) {

	var wg sync.WaitGroup

	f, err := os.Open(d.UrlTXTPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	urlsLinkedList := new(Node)
	tmpNode := urlsLinkedList
	counter := 0

	for s.Scan() {
		counter++
		tmpNode.data = s.Text()
		tmpNode.next = new(Node)
		tmpNode = tmpNode.next
	}

	var jobs = make(chan *string, counter)
	var results = make(chan *rawResponse, counter)

	for i := 0; i < 200; i++ {
		go download(i, jobs, results)
	}

	tmpNode = urlsLinkedList
	for tmpNode.data != nil {
		str := tmpNode.data.(string)
		jobs <- &str
		tmpNode = tmpNode.next
	}

	close(jobs)

	unsavedResponses := make([]rawResponse, 0)
	upsertedPages := make([]Response, 0)
	changeLogValues := make([]Response, 0)
	filterByUrl := make([]bson.M, 0)
	checkProgressCounter := 0

	gscan := scan{
		wg:           &wg,
		db:           d.Database,
		docsFrmNet:   &unsavedResponses,
		maxLen:       counter,
		counter:      &checkProgressCounter,
		changedPages: &upsertedPages,
		docsFrmDB:    &changeLogValues,
		filter:       &filterByUrl,
	}

	for i := 0; i < counter; i++ {
		var page = *<-results
		if page.err == nil {
			wg.Add(1)
			gscan.res = page
			go checkChanges(&gscan)
		}
	}

	wg.Wait()

	vChangedMap := make(map[string]*Response)
	vfrmNetMap := make(map[string]*rawResponse)

	for i := 0; i < len(*gscan.changedPages); i++ {
		vChangedMap[(*gscan.changedPages)[i].Url] = &(*gscan.changedPages)[i]
	}

	for i := 0; i < len(*gscan.docsFrmNet); i++ {
		vfrmNetMap[*((*gscan.docsFrmNet)[i].url)] = &(*gscan.docsFrmNet)[i]
	}

	for _, response := range *gscan.docsFrmNet {
		if _, found := vChangedMap[*response.url]; found {
			wg.Add(1)
			go dumpHtml(&wg, response)
		}
	}

	wg.Wait()

	if 0 != len(*gscan.changedPages) {
		zipFolder("./"+folder, folder+".zip")

		var writer bytes.Buffer
		tmpl, err := template.ParseFiles("./template/mail.html")
		if err != nil {
			panic(err)
		}

		err = tmpl.Execute(&writer, *gscan.changedPages)
		if err != nil {
			panic(err)
		}

		mail := NewMail(MailConf{
			FromEmail:  smtpEmail,
			FromName:   senderName,
			ToEmail:    toEmail,
			Subject:    fmt.Sprintf("Here you go! %d Changes", len(*gscan.changedPages)),
			Message:    writer.String(),
			Attachment: "." + string(os.PathSeparator) + folder + ".zip",
			Password:   gmailAppPassword,
			SmtpHost:   smtpHost,
			SmtpPort:   smtpPort,
			AdminEmail: smtpEmail,
		})

		mail.SendMail()
	}

	return &gscan, nil

}
