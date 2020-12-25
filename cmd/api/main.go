package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	_"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/keepcalmist/workwithElastic/pkg/Server"
	"github.com/keepcalmist/workwithElastic/pkg/config"
	"github.com/keepcalmist/workwithElastic/pkg/storage"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	_ "strconv"
	"strings"
	"sync"
)

var (
	LOG_FILE = "logs.log"
)

func init(){
	_ = config.New()
}

func main(){
	//ctx, cancel  := context.WithCancel(context.Background())
	//defer cancel()
	ctx := context.Background()
	var wg sync.WaitGroup
	//creating new file to write logs
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		//TODO add handle this error
		panic("Cant open or create logFile")
		return
	}
	//init logger
	logs := initLogger(file)

	status := make(chan int,1 )
	//server with chan to graceful shutdown
	go Server.Run(status,logs)


	cl, err := storage.InitElastic()
	fmt.Println(err,cl)
	if err != nil {
		logs.Println(err)
	}

	var r map[string]interface{}
	res, err := cl.Info()
	if err != nil {
		logs.Println("Error getting info")
		fmt.Println("Error getting info", err)
	}
	defer res.Body.Close()


	//fmt.Println(res)
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logs.Fatalf("Error parsing the response body: %s", err)
	}
	//Show info about ElasticSearch server and client
	{
		fmt.Println(strings.Repeat("~", 37))
		fmt.Printf("Client: %s\n", elasticsearch.Version)
		fmt.Printf("Server: %s\n", r["version"].(map[string]interface{})["number"])
		fmt.Println(strings.Repeat("~", 37))
	}

	for i, title := range []string{"Test One", "Test Two"} {
		wg.Add(1)
		fmt.Println("add waitgroup")
		go func(i int, title string) {
			fmt.Println("in func")
			var b strings.Builder
			b.WriteString(`{"title" : "`)
			b.WriteString(title)
			b.WriteString(`"}`)

			req := esapi.IndexRequest{
				Index:      "test",
				DocumentID: strconv.Itoa(i),
				Body:       strings.NewReader(b.String()),
				Refresh:    "true",
			}
			fmt.Println("before do")
			res, err := req.Do(context.Background(), cl)
			fmt.Println("after do")
			if err != nil{
				fmt.Println(err)
				logs.Println(err)
			}
			fmt.Println(err)
			defer res.Body.Close()

			if res.IsError() {
				logs.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					logs.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					fmt.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
			wg.Done()
		}(i, title)

	}
	wg.Wait()

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "one",
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		logs.Fatalf("Error encoding query: %s", err)
	}

	res, err = cl.Search(
		cl.Search.WithContext(ctx),
		cl.Search.WithIndex("test"),
		cl.Search.WithBody(&buf),
		cl.Search.WithTrackTotalHits(true),
		cl.Search.WithPretty(),
	)
	if err != nil {
		logs.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			logs.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			logs.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logs.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	fmt.Printf(
		"\n[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		logs.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
		fmt.Printf("\n * ID=%s, %s\n", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	fmt.Println(strings.Repeat("=", 37))

	fmt.Println("Exit from main with status: ")
}


func initLogger(file *os.File) *log.Logger {
	logger := log.New()
//Set output for logs
	logger.SetOutput(file)
	logger.SetFormatter(&log.TextFormatter{})
	log.RegisterExitHandler(ShutDown)
	return logger
}


func ShutDown () {
	log.Exit(0)
}