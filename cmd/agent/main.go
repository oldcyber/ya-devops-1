package main

import (
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"ya-devops-1/internal/agent"
	"ya-devops-1/internal/tools"
)

func main() {
	checkEnv := func(key string) bool {
		_, ok := os.LookupEnv(key)
		if !ok {
			return false
		} else {
			return true
		}
	}

	Address := flag.String("a", "", "address")
	ReportInterval := flag.Int("r", 0, "report interval")
	PoolInterval := flag.Int("p", 0, "pool interval")
	flag.Parse()

	log.Println("Start agent")
	if !checkEnv("ADDRESS") && *Address != "" {
		tools.Conf.Address = *Address
	}
	if !checkEnv("REPORT_INTERVAL") && *ReportInterval != 0 {
		tools.Conf.ReportInterval = time.Duration(*ReportInterval) * time.Second
	}
	if !checkEnv("POLL_INTERVAL") && *PoolInterval != 0 {
		tools.Conf.PollInterval = time.Duration(*PoolInterval) * time.Second
	}
	log.Println("loading config. Address:", tools.Conf.Address, "Poll interval:", tools.Conf.PollInterval.Seconds(), "Report interval", tools.Conf.ReportInterval.Seconds())
	// log.Println("loading config. Address:", *Address, "Poll interval:", *PoolInterval, "Report interval", *ReportInterval)

	// tools.Conf.Address = "localhost:8080"
	// cfg := tools.NewConfig()
	//cfg := config{}
	//if err := env.Parse(&cfg); err != nil {
	//	fmt.Printf("%+v\n", err)
	//}
	//log.Printf("%+v\n", Cfg)
	agent.WorkWithMetrics()
}
