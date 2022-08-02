package tools

import (
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  int           `env:"STORE_INTERVAL" envDefault:"10"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
}

var Conf = config{}

func init() {
	if err := env.Parse(&Conf); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
