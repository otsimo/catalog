package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	homedir "github.com/mitchellh/go-homedir"
)

type AuthConfig struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type Config struct {
	Auths map[string]AuthConfig `json:"auths"`
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

//NewConfig reads ~/otsimo/config.json, if it is not available than creates it
func NewConfig() (*Config, error) {
	cp := getConfigPath()
	if !exists(cp) {
		cnf := &Config{Auths: make(map[string]AuthConfig)}
		err := cnf.Save()
		if err != nil {
			logrus.Errorf("Unable to save config to file, error=%+v", err)
			return nil, err
		}
		return cnf, nil
	}
	configFile, err := os.Open(cp)
	defer configFile.Close()
	if err != nil {
		logrus.Errorf("Error open config file, error=%+v", err)
		return nil, err
	}
	var config Config
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		logrus.Errorf("Error decoding config file, error=%+v", err)
		return nil, err
	}
	return &config, nil
}

func getConfigDir() string {
	d, e := homedir.Dir()
	if e != nil {
		logrus.Fatalf("failed to get home dir: error=%v", e)
	}
	return path.Join(d, ".otsimo")
}

func getConfigPath() string {
	return path.Join(getConfigDir(), "config.json")
}

func (c *Config) Save() error {
	d, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	cdir := getConfigDir()
	if !exists(cdir) {
		os.Mkdir(cdir, os.ModePerm)
	}
	return ioutil.WriteFile(getConfigPath(), d, os.ModePerm)
}

func (c *Config) Token(host string) string {
	ac, ok := c.Auths[host]
	if !ok {
		return ""
	}
	return ac.Token
}
