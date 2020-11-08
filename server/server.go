package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

var (
	MaxWorkers int
	MaxQueue   int
	count      uint64
)

type Foo struct {
	Key   int `json:"key"`
	Value int `json:"value"`
}

type Job struct {
	Foo Foo
}

func (f *Foo) do() {
	log.Println(f, count)
	atomic.AddUint64(&count, 1)
}

func fooHandler(rw http.ResponseWriter, req *http.Request) {

	decoder := json.NewDecoder(req.Body)
	var foo Foo
	err := decoder.Decode(&foo)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	defer req.Body.Close()

	go func(foo Foo) {
		work := Job{foo}

		JobQueue <- work
	}(foo)

}

var JobQueue chan Job

func main() {

	fmt.Print("Enter count of worker: ")
	fmt.Scanf("%d", &MaxWorkers)

	fmt.Print("Enter size of queue: ")
	fmt.Scanf("%d", &MaxQueue)

	JobQueue = make(chan Job, MaxQueue)

	log.Println("main successful")

	dispatcher := NewDispatcher(MaxWorkers)
	dispatcher.Run()

	http.HandleFunc("/", fooHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("starting listening")
	} else {
		fmt.Errorf("server load error %s", err.Error())
	}
}
