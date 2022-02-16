package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database struct {
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		DbName string `yaml:"dbname"`
	} `yaml:"database"`
	Directory struct {
		UploadUsers   string `yaml:"import_users"`
		MailTemplates string `yaml:"mail_templates"`
	} `yaml:"directory"`
	Mail struct {
		EmaiFrom string `yaml:"from"`
		SmtpHost string `yaml:"smtp_host"`
		SmtpPort string `yaml:"smtp_port"`
		SmtpPass string `yaml:"smtp_user"`
		SmtpPwd  string `yaml:"smtp_pwd"`
	} `yaml:"mail"`
	Log struct {
		LogDir string `yaml:"log_dir"`
	} `yaml:"log"`
}

func (c Config) ApiUrl() string {
	return fmt.Sprintf("http://%s:%d/api/v1", c.Host, c.Port)
}

var conf Config

func ParseConfig(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal([]byte(data), &conf); err != nil {
		return err
	}
	return nil
}

func GetConfig() Config {
	return conf
}
