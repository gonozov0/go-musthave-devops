package main

import (
	"log"
	"time"

	"github.com/gonozov0/go-musthave-devops/internal/agent"
)

func main() {
	config, err := agent.LoadConfig()
	if err != nil {
		log.Fatalf("ERROR loading config: %v", err)
	}

	collectTicker := time.NewTicker(time.Duration(config.CollectInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(config.SendInterval) * time.Second)

	var metrics []agent.Metric

	for {
		select {
		case <-collectTicker.C:
			metrics = append(metrics, agent.CollectMetrics()...)
		case <-sendTicker.C:
			metrics, err = agent.SendMetrics(metrics, config.ServerAddress)
			if err != nil {
				log.Fatalf("ERROR sending metrics: %v", err)
			}
		}
	}
}
