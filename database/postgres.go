package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

var DB *sql.DB

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

func Init() error {
	// Read the configuration file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		return err
	}

	// Parse the configuration file
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// Create the database connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.SSLMode,
	)

	// Open the database connection
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Test the database connection
	err = DB.Ping()
	if err != nil {
		return err
	}

	return nil
}
