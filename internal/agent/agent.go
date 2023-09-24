package agent

import (
	"github.com/gonozov0/go-musthave-devops/internal/agent/internal"
	"github.com/gonozov0/go-musthave-devops/internal/config"
	"log"
	"time"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("ERROR loading config: %v", err)
	}

	collectTicker := time.NewTicker(time.Duration(cfg.CollectInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(cfg.SendInterval) * time.Second)

	var metrics []internal.Metric

	for {
		select {
		case <-collectTicker.C:
			metrics = append(metrics, internal.CollectMetrics()...)
		case <-sendTicker.C:
			metrics, err = internal.SendMetrics(metrics, cfg.ServerAddress)
			if err != nil {
				log.Fatalf("ERROR sending metrics: %v", err)
			}
		}
	}
}
