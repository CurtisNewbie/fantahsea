package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type Configuration struct {
	DBConf     DBConfig     `json:"db"`
	ServerConf ServerConfig `json:"server"`
	FileConf   FileConfig   `json:"file"`
}

type DBConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type FileConfig struct {
	Base string `json:"base"`
}

/* Parse json config file */
func ParseJsonConfig(filePath string) (*Configuration, error) {

	file, err := os.Open(filePath)
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
	return &configuration, nil
}

/*
	Parse Cli Arg to extract a profile

	It looks for the arg that matches the pattern "profile=[profileName]"
	For example, for "profile=prod", the extract profile is "prod"
*/
func ParseProfile(args []string) string {
	profile := "dev" // the default one

	for _, s := range args {
		var eq int = strings.Index(s, "=")
		if eq != -1 {
			if key := s[:eq]; key == "profile" {
				profile = s[eq+1:]
				break
			}
		}
	}

	if strings.TrimSpace(profile) == "" {
		profile = "dev"
	}
	return profile
}
