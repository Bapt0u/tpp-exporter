package tpp

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type VenafiTppConf struct {
	Username   string   `yaml:"username"`
	Password   string   `yaml:"password"`
	Url        string   `yaml:"url"`
	Scope      string   `yaml:"scope"`
	ClientId   string   `yaml:"client_id"`
	Policies   []string `yaml:"policies"`
	ScrapeTime int      `yaml:"scrapetime"`
	ExpireSoon int      `yaml:"expireSoon"`
}

type Config struct {
	VenafiTpp VenafiTppConf `yaml:"venafi_tpp"`
}

func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IsIn(l []string, s string) bool {
	for _, l := range l {
		if l == s {
			return true
		}
	}
	return false
}

func (c *Config) GetConf(path string) *Config {
	confFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("system; error=Error loading file: %e", err)
	}
	err = yaml.Unmarshal(confFile, c)
	if err != nil {
		log.Printf("system; error=Error unmarshaling conf: %e", err)
	}

	log.Printf("system; msg=%s successfully loaded. \n", path)
	return c
}
