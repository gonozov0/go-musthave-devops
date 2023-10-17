package main

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gonozov0/go-musthave-devops/internal/agent"
	"github.com/gonozov0/go-musthave-devops/internal/shared"
)

func main() {
	cfg, err := agent.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %s", err.Error())
	}

	collectTicker := time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)

	var metrics []shared.Metric

	for {
		select {
		case <-collectTicker.C:
			metrics = append(metrics, agent.CollectMetrics()...)
		case <-sendTicker.C:
			log.Infof("Sending %d metrics", len(metrics))
			metrics, err = agent.SendMetrics(metrics, cfg.ServerAddress)
			if err != nil {
				log.Fatalf("Could not send metrics: %s", err.Error())
			}
		}
	}
}
