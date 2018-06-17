package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c PostgresConfig) Dialect() string {
	return "postgres"
}

func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "jiji_dev_user",
		Password: "123test",
		Name:     "jiji_dev",
	}
}

func defaultMailgunConfig() MailgunConfig {
	return MailgunConfig{
		APIKey:       "key-dsafjlsdafh9ds890890sdfasfd",
		PublicAPIKey: "pubkey-124321498sdfdsf8u987dsfsdaf08",
		Domain:       "jiji_demo.com",
	}
}

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
	Mailgun  MailgunConfig  `json:"mailgun"`
}

type MailgunConfig struct {
	APIKey       string `json:"api_key"`
	PublicAPIKey string `json:"public_api_key"`
	Domain       string `json:"domain"`
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func DefaultConfig() Config {
	return Config{
		Port:     3000,
		Env:      "dev",
		Pepper:   "super-secret-pepper-for-password",
		HMACKey:  "super-secret-hmac-key",
		Database: DefaultPostgresConfig(),
		Mailgun:  defaultMailgunConfig(),
	}
}

func LoadConfig(configRequired bool) Config {
	// Open the config file
	configFile, err := os.Open(".config")
	if err != nil {
		if configRequired {
			panic(err)
		}
		// If there was an error opening the file, print out a
		// message saying we are using the default config and
		// return it.
		fmt.Println("Using the default config...")
		return DefaultConfig()
	}
	// If we opened the config file successfully we are going
	// to create a Config variable to load it into.
	var config Config
	// We also need a JSON decoder, which will read from the
	// file we opened when decoding
	dec := json.NewDecoder(configFile)
	// We then decode the file and place the results in c, the
	// Config variable we created for the results. The decoder
	// knows how to decode the data because of the struct tags
	// (eg `json:"port"`) we added to our Config and
	// PostgresConfig fields, much like GORM uses struct tags
	// to know which database column each field maps to.
	err = dec.Decode(&config)
	if err != nil {
		panic(err)
	}
	// If all goes well, return the loaded config.
	fmt.Println("Successfully loaded .config")
	return config
}
