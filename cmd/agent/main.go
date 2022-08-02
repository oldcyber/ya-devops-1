package main

import (
	log "github.com/sirupsen/logrus"
	"ya-devops-1/internal/agent"
	"ya-devops-1/internal/tools"
)

func main() {
	log.Println("Start agent")
	log.Println("loading config. Address:", tools.Conf.Address, "Poll interval:", tools.Conf.PollInterval.Seconds(), "Report interval", tools.Conf.ReportInterval.Seconds())

	// tools.Conf.Address = "localhost:8080"
	// cfg := tools.NewConfig()
	//cfg := config{}
	//if err := env.Parse(&cfg); err != nil {
	//	fmt.Printf("%+v\n", err)
	//}
	//log.Printf("%+v\n", Cfg)
	agent.WorkWithMetrics()
}
