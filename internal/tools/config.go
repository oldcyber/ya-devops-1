package tools

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type config struct {
	Address        string        `env:"ADDRESS" envDefault:"localhost:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
}

func (c *config) GetAddress() string {
	return c.Address
}

func (c *config) GetPollInterval() time.Duration {
	return c.PollInterval
}

func (c *config) GetStoreFile() string {
	return c.StoreFile
}

func (c *config) GetStoreInterval() time.Duration {
	return c.StoreInterval
}

func (c *config) GetReportInterval() time.Duration {
	return c.ReportInterval
}

func (c *config) GetRestore() bool {
	return c.Restore
}

func (c *config) InitFromEnv() error {
	if err := env.Parse(c); err != nil {
		log.Error(err)
		return err
	}
	log.Info("Config after env read:", *c)
	return nil
}

func checkEnv(key string) bool {
	_, ok := os.LookupEnv(key)
	if !ok {
		return false
	} else {
		return true
	}
}

func (c *config) InitFromFlags() error {
	Address := flag.String("a", "", "address")
	Restore := flag.Bool("r", true, "restore")
	StoreInterval := flag.Duration("i", 0, "store interval")
	StoreFile := flag.String("f", "", "store file")
	flag.Parse()
	if !checkEnv("ADDRESS") && *Address != "" {
		c.Address = *Address
	}
	if !checkEnv("RESTORE") && *Restore != c.Restore {
		c.Restore = *Restore
	}
	if !checkEnv("STORE_INTERVAL") && *StoreInterval != 0 {
		c.StoreInterval = *StoreInterval
	}
	if !checkEnv("STORE_FILE") && *StoreFile != "" {
		c.StoreFile = *StoreFile
	}
	log.Info("Config after flags read:", *c)
	return nil
}

func NewConfig() *config {
	return &config{
		Address:        "localhost:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
		StoreInterval:  300 * time.Second,
		StoreFile:      "/tmp/devops-metrics-db.json",
		Restore:        true,
	}
}

func (c *config) PrintConfig() {
	log.Info("Config after all init:", *c)
}
