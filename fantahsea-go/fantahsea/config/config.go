package config

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {

	DBConf DBConfig `json:"db"`
	ServerConf ServerConfig `json:"server"`

}

type DBConfig struct {
	User string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
} 

/* Parse json config file */
func ParseJsonConfig(filePath string) (*Configuration, error) {

	file, err := os.Open(filePath);
	if err != nil {
		log.Printf("Failed to open config file, %v\n", err)
		return nil, err
	}

	defer file.Close()

	jsonDecoder := json.NewDecoder(file)

	configuration := Configuration{}
	err = jsonDecoder.Decode(&configuration)
	if err != nil {
		log.Printf("Failed to decode config file as json, %v\n", err)
		return nil, err
	}

	log.Printf("Parsed json config file: '%v'\n", filePath)
	return &configuration, nil;
}
