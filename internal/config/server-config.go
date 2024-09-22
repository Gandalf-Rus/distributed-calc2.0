package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultServerPort       = "8080"
	defaultGrpcPort         = "5000"
	defaultDBHost           = "localhost" // "db"
	defaultDBPort           = "5432"
	defaultDBUser           = "postgres"
	defaultDBPassword       = "postgres"
	defaultDBName           = "distributedcalc"
	defaultJWTTokenTimeout  = 30 * time.Hour
	defaultAgentLostTimeout = 10 * time.Second
)

type Config struct {
	ServerPort       int
	GrpcHost         string
	GrpcPort         int
	Dbhost           string
	Dbport           int
	Dbuser           string
	Dbpassword       string
	Dbname           string
	OperatorsDelay   OperatorsDelay
	AgentLostTimeout time.Duration
	JwtTokenTimeout  time.Duration
}

var Cfg Config

func InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file")
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = defaultServerPort
	}

	port, err := strconv.Atoi(serverPort)
	if err != nil {
		return fmt.Errorf("failed to parse %s as int: %w", os.Getenv("SERVER_PORT"), err)
	}

	grpcHost := os.Getenv("GRPC_HOST")

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = defaultGrpcPort
	}

	grpcport, err := strconv.Atoi(grpcPort)
	if err != nil {
		return fmt.Errorf("failed to parse %s as int: %w", os.Getenv("GRPC_PORT"), err)
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = defaultDBHost
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = defaultDBPort
	}
	dbport, err := strconv.Atoi(dbPort)
	if err != nil {
		return fmt.Errorf("failed to parse %s as int: %w", os.Getenv("DB_PORT"), err)
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = defaultDBUser
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = defaultDBPassword
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = defaultDBName
	}

	Cfg.ServerPort = port
	Cfg.GrpcHost = grpcHost
	Cfg.GrpcPort = grpcport
	Cfg.Dbhost = dbHost
	Cfg.Dbport = dbport
	Cfg.Dbuser = dbUser
	Cfg.Dbpassword = dbPassword
	Cfg.Dbname = dbName

	Cfg.AgentLostTimeout = defaultAgentLostTimeout
	Cfg.JwtTokenTimeout = defaultJWTTokenTimeout
	Cfg.OperatorsDelay = OperatorsDelay{
		DelayForAdd: 0,
		DelayForSub: 0,
		DelayForMul: 0,
		DelayForDiv: 0,
	}

	return nil
}

func (cfg *Config) ChangeOpDuration(op string, duration int) error {
	if duration < 0 {
		return fmt.Errorf("operation duration must be positive")
	}
	switch op {
	case "+":
		cfg.OperatorsDelay.DelayForAdd = duration
	case "-":
		cfg.OperatorsDelay.DelayForSub = duration
	case "*":
		cfg.OperatorsDelay.DelayForMul = duration
	case "/":
		cfg.OperatorsDelay.DelayForDiv = duration
	default:
		return fmt.Errorf("unsupported operator")
	}
	return nil
}
