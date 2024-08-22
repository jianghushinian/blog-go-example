package main

import (
	"log"
	"time"

	genericapiserver "k8s.io/apiserver/pkg/server"
	// genericapiserver "github.com/jianghushinian/blog-go-example/gracefulstop/pkg/server"
)

type Syncer struct {
	interval time.Duration
}

func (s *Syncer) Run(quit <-chan struct{}) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 业务逻辑
			log.Println("do something")
		case <-quit:
			log.Println("Stop loop")
			s.Stop()
			return nil
		}
	}
}

func (s *Syncer) Stop() {
	log.Println("Stop Syncer start")
	time.Sleep(time.Second * 5)
	log.Println("Stop Syncer done")
}

func main() {
	s := Syncer{interval: time.Second}
	quit := genericapiserver.SetupSignalHandler()

	if err := s.Run(quit); err != nil {
		log.Fatalf("Syncer run err: %s", err.Error())
	}
}
