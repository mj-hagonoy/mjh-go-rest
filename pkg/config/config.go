package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Credentials struct {
		GoogleCloud string `yaml:"google_app_creds"`
	} `yaml:"credentials"`
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
	FileStorage struct {
		Default     string `yaml:"default"`
		GoogleCloud struct {
			ProjectID  string `yaml:"project_id"`
			BucketName string `yaml:"bucket_name"`
			UploadPath string `yaml:"upload_path"`
		} `yaml:"google_cloud"`
	} `yaml:"file_storage"`
	Messaging struct {
		GoogleCloud struct {
			ProjectID string `yaml:"project_id"`
			TopicID   string `yaml:"topic_id"`
		} `yaml:"google_cloud"`
	} `yaml:"messaging"`
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
