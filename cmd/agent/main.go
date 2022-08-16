package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/oldcyber/ya-devops-1/internal/agent"
	"github.com/oldcyber/ya-devops-1/internal/tools"
)

func main() {
	cfg := tools.NewConfig()
	if err := cfg.InitFromEnv(); err != nil {
		log.Error(err)
		return
	}
	if err := cfg.InitFromFlags(); err != nil {
		log.Error(err)
		return
	}
	cfg.PrintConfig()
	err := agent.WorkWithMetrics(cfg)
	if err != nil {
		log.Error(err)
		return
	}
}
