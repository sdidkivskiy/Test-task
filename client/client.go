package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const url = "http://localhost:8080/"

var (
	MaxWorkers int
	MaxQueue   int
	JobQueue   chan Job
)

type Foo struct {
	Key   int `json:"key"`
	Value int `json:"value"`
}

func (foo Foo) do(client http.Client) {
	var jsonStr, err = json.Marshal(foo)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	defer resp.Body.Close()
}

type Job struct {
	Foo Foo
}

func main() {

	var rateLimit int
	var countOfRequests int

	fmt.Print("Enter count of worker: ")
	fmt.Scanf("%d", &MaxWorkers)

	fmt.Print("Enter size of queue: ")
	fmt.Scanf("%d", &MaxQueue)

	JobQueue = make(chan Job, MaxQueue)

	fmt.Print("Enter count of requests per minute: ")
	fmt.Scanf("%d", &rateLimit)

	fmt.Print("Enter count of total requests: ")
	fmt.Scanf("%d", &countOfRequests)

	t := float64(60.0 / float64(rateLimit))
	t1 := strconv.FormatFloat(t, 'f', 8, 64) + "s"
	timeout, err := (time.ParseDuration(t1))

	fmt.Println(t1)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	fmt.Println(timeout)
	fmt.Println(time.Now())

	dispatcher := NewDispatcher(MaxWorkers)
	dispatcher.Run()

	for i := 0; i < countOfRequests; i++ {
		go func() {
			foo := &Foo{Key: i, Value: i}
			work := Job{Foo: *foo}
			JobQueue <- work

		}()
		time.Sleep(timeout)
	}

	fmt.Println(time.Now())
}
