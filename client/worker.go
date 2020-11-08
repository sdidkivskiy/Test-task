package main

import "net/http"

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
	Client     http.Client
}

func NewWorker(workerPool chan chan Job) Worker {

	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job),
		quit:       make(chan bool),
		Client:     http.Client{}}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				job.Foo.do(w.Client)
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
