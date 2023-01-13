// SDR - Labo 4
// Nicolas Crausaz & Maxime Scharwath

package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Config struct {
	Servers []ServerConfig `json:"servers"`
}

type ServerConfig struct {
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Letter     string `json:"letter"`
	Neighbours []int  `json:"neighbours"`
}

func (d *ServerConfig) Address() (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", d.Host, d.Port))
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
	/*
		utils.SetLogEnabled(config.Logs)
		utils.SetDebugEnabled(config.Debug)
		utils.SetDebugDuration(config.DebugDuration * time.Millisecond)
	*/
	return config, nil
}
