package tools

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL" envDefault:"10s"`
	PollInterval   time.Duration `env:"POLL_INTERVAL" envDefault:"2s"`
	StoreInterval  time.Duration `env:"STORE_INTERVAL" envDefault:"300s"`
	StoreFile      string        `env:"STORE_FILE" envDefault:"/tmp/devops-metrics-db.json"`
	Restore        bool          `env:"RESTORE" envDefault:"true"`
	Key            string        `env:"KEY" envDefault:""`
	DatabaseDSN    string        `env:"DATABASE_DSN"`
}

func (c *Config) GetAddress() string {
	return c.Address
}

func (c *Config) GetPollInterval() time.Duration {
	return c.PollInterval
}

func (c *Config) GetStoreFile() string {
	return c.StoreFile
}

func (c *Config) GetStoreInterval() time.Duration {
	return c.StoreInterval
}

func (c *Config) GetReportInterval() time.Duration {
	return c.ReportInterval
}

func (c *Config) GetRestore() bool {
	return c.Restore
}

func (c *Config) GetKey() string {
	return c.Key
}

func (c *Config) GetDatabaseDSN() string {
	return c.DatabaseDSN
}

func (c *Config) InitFromEnv() error {
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

func (c *Config) InitFromServerFlags() error {
	Address := flag.String("a", "", "address")
	Restore := flag.Bool("r", true, "restore")
	StoreInterval := flag.Duration("i", 0, "store interval")
	StoreFile := flag.String("f", "", "store file")
	Key := flag.String("k", "", "key")
	DatabaseDSN := flag.String("d", "", "database dsn")
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
	if !checkEnv("KEY") && *Key != "" {
		c.Key = *Key
	}
	if !checkEnv("DATABASE_DSN") && *DatabaseDSN != "" {
		c.DatabaseDSN = *DatabaseDSN
	}
	log.Info("Config after flags read:", *c)
	return nil
}

func (c *Config) InitFromAgentFlags() error {
	Address := flag.String("a", "", "address")
	ReportInterval := flag.Duration("r", 0, "report interval")
	PoolInterval := flag.Duration("p", 0, "poll interval")
	Key := flag.String("k", "", "key")
	flag.Parse()
	if !checkEnv("ADDRESS") && *Address != "" {
		c.Address = *Address
	}
	if !checkEnv("REPORT_INTERVAL") && *ReportInterval != 0 {
		c.ReportInterval = *ReportInterval
	}
	if !checkEnv("POLL_INTERVAL") && *PoolInterval != 0 {
		c.PollInterval = *PoolInterval
	}
	if !checkEnv("KEY") && *Key != "" {
		c.Key = *Key
	}
	log.Info("Config after flags read:", *c)
	return nil
}

func NewConfig() *Config {
	return &Config{
		Address:        "localhost:8080",
		ReportInterval: 10 * time.Second,
		PollInterval:   2 * time.Second,
		StoreInterval:  300 * time.Second,
		StoreFile:      "/tmp/devops-metrics-db.json",
		Restore:        true,
		Key:            "",
		DatabaseDSN:    "",
	}
}

func (c *Config) PrintConfig() {
	log.Info("Config after all init:", *c)
}
