// SDR - Labo 4
// Nicolas Crausaz & Maxime Scharwath

package config

import (
	"SDR-Labo4/src/utils/log"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type LogConfig struct {
	Enabled bool `json:"enabled"`
	Level   int  `json:"level"`
}

type Config struct {
	Servers []ServerConfig `json:"servers"`
	Log     LogConfig      `json:"log"`
}

type ServerConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Letter     string `json:"letter"`
	Neighbours []int  `json:"neighbours"`
}

func (s *ServerConfig) FullAddress() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s *ServerConfig) Address() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", s.FullAddress())
}

// LoadConfig loads the config file and returns a Config struct
func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	log.SetLogEnabled(config.Log.Enabled)
	log.SetLogLevelByValue(config.Log.Level)

	return config, nil
}
