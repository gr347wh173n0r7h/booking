package config

import (
	"os"
	"strconv"
)

// Config defines Booking service config
type Config struct {
	Hostname        string
	ListenPort      int
	LogLevel        string
	DBURL           string
	DBLog           bool
	MaxTimeBlockMin int
	SwaggerDistPath string
}

// NewDefaults returns a default Config
// "postgres://postgres:test@localhost:5432/booking?sslmode=disable"
// "/Users/jpetersen/Workspace/practice/swagger-ui/dist"
func NewDefaults() *Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8081
	}

	timeblocks, err := strconv.Atoi(os.Getenv("MAXTIMEBLOCK"))
	if err != nil {
		timeblocks = 60
	}

	return &Config{
		Hostname:        os.Getenv("HOST"),
		ListenPort:      port,
		LogLevel:        os.Getenv("LOGLEVEL"),
		DBURL:           os.Getenv("DBURL"),
		DBLog:           false,
		MaxTimeBlockMin: timeblocks,
		SwaggerDistPath: os.Getenv("SWAGGERDIST"),
	}
}
