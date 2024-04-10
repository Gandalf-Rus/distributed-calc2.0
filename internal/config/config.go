package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	defaultServerPort = "8080"
)

type Config struct {
	Dbhost           string
	Dbuser           string
	Dbpassword       string
	Dbname           string
	Dbport           int
	ServerPort       int
	DelayForAdd      int
	DelayForSub      int
	DelayForMul      int
	DelayForDiv      int
	AgentLostTimeout int
}

var Cfg Config

func InitConfig() error {

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = defaultServerPort
	}

	port, err := strconv.Atoi(serverPort)
	if err != nil {
		return fmt.Errorf("failed to parse %s as int: %w", os.Getenv("SERVER_PORT"), err)
	}

	flag.StringVar(&Cfg.Dbhost, "dbhost", "localhost", "Postgress host")
	flag.StringVar(&Cfg.Dbuser, "dbuser", "postgres", "Postgress user")
	flag.StringVar(&Cfg.Dbpassword, "dbpassword", "postgres", "Postgress password")
	flag.IntVar(&Cfg.Dbport, "dbport", 5432, "Posgress port")
	flag.StringVar(&Cfg.Dbname, "dbname", "distribcalc", "Postgress database name")

	flag.IntVar(&Cfg.ServerPort, "httppport", port, "HTTP port to listen to")
	flag.IntVar(&Cfg.AgentLostTimeout, "agenttimeout", 60, "Timeout before agent considered lost (seconds)")

	flag.Parse()

	Cfg.DelayForAdd = 10
	Cfg.DelayForSub = 12
	Cfg.DelayForMul = 15
	Cfg.DelayForDiv = 20

	return nil
}
