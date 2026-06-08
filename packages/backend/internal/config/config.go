package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Deployments struct {
	ChainID  uint64                   `json:"chainId"`
	RpcURL   string                   `json:"rpcUrl"`
	Proxies  map[string]common.Address `json:"proxies"`
}

func LoadDeployments(path string) (*Deployments, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var d Deployments
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func InitDB() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "linkworld")
	dbname := getEnv("DB_NAME", "linkworld")

	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("PostgreSQL connected successfully")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
