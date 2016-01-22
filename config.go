package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/oursky/gcfg"
)

// Configuration is Skygear's configuration
type Configuration struct {
	HTTP struct {
		Host string
	}
	App struct {
		Name          string
		APIKey        string `gcfg:"api-key"`
		AccessControl string `gcfg:"access-control"`
		Dev           bool
	}
	DB struct {
		ImplName string `gcfg:"implementation"`
		Option   string
	}
	TokenStore struct {
		ImplName string `gcfg:"implementation"`
		Path     string `gcfg:"path"`
	} `gcfg:"token-store"`
	AssetStore struct {
		ImplName string `gcfg:"implementation"`
		Public   bool

		// followings only used when ImplName = fs
		Path string

		// followings only used when ImplName = s3
		AccessToken string `gcfg:"access-key"`
		SecretToken string `gcfg:"secret-key"`
		Region      string
		Bucket      string
	} `gcfg:"asset-store"`
	AssetURLSigner struct {
		URLPrefix string `gcfg:"url-prefix"`
		Secret    string
	} `gcfg:"asset-url-signer"`
	APNS struct {
		Enable   bool
		Env      string
		Cert     string
		Key      string
		CertPath string `gcfg:"cert-path"`
		KeyPath  string `gcfg:"key-path"`
	}
	GCM struct {
		Enable bool
		APIKey string `gcfg:"api-key"`
	}
	LOG struct {
		Level string
	}
	LogHook struct {
		SentryDSN   string `gcfg:"sentry-dsn"`
		SentryLevel string `gcfg:"sentry-level"`
	} `gcfg:"log-hook"`
	Plugin map[string]*struct {
		Transport string
		Path      string
		Args      []string
	}
	// the alembic section here is to make the config be parsed correctly
	// the values should not be used
	UselessAlembic struct {
		ScriptLocation string `gcfg:"script_location"`
	} `gcfg:"alembic"`
}

// ReadFileInto reads a configuration from file specified by path
func ReadFileInto(config *Configuration, path string) error {
	if err := gcfg.ReadFileInto(config, path); err != nil {
		return err
	}
	if config.HTTP.Host == "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}
		config.HTTP.Host = ":" + port
	}

	appAPIKey := os.Getenv("API_KEY")
	if appAPIKey != "" {
		config.App.APIKey = appAPIKey
	}

	appName := os.Getenv("APP_NAME")
	if appName != "" {
		config.App.Name = appName
	}
	if config.App.Name == "" {
		return errors.New("app name is not set")
	}
	if !regexp.MustCompile("^[A-Za-z0-9_]+$").MatchString(config.App.Name) {
		return fmt.Errorf("app name '%s' contains invalid characters other than alphanumberics or underscores", config.App.Name)
	}

	if config.App.AccessControl == "" {
		config.App.AccessControl = "role"
	}

	dbImplName := os.Getenv("DB_IMPL_NAME")
	if dbImplName != "" {
		config.DB.ImplName = dbImplName
	}

	if config.DB.ImplName == "pq" && os.Getenv("DATABASE_URL") != "" {
		config.DB.Option = os.Getenv("DATABASE_URL")
	}

	err := readAPNS(config)
	if err != nil {
		return err
	}
	return nil
}

func readAPNS(config *Configuration) error {
	shouldEnableAPNS := os.Getenv("APNS_ENABLE")
	if shouldEnableAPNS != "" {
		config.APNS.Enable = shouldEnableAPNS == "1"
	}
	if !config.APNS.Enable {
		return nil
	}

	env := os.Getenv("APNS_ENV")
	if env != "" {
		config.APNS.Env = env
	}

	cert, key := os.Getenv("APNS_CERTIFICATE"), os.Getenv("APNS_PRIVATE_KEY")
	if cert != "" {
		config.APNS.Cert = cert
	}
	if key != "" {
		config.APNS.Key = key
	}

	if config.APNS.Cert == "" && config.APNS.CertPath != "" {
		certPEMBlock, err := ioutil.ReadFile(config.APNS.CertPath)
		if err != nil {
			return err
		}
		config.APNS.Cert = string(certPEMBlock)
	}

	if config.APNS.Key == "" && config.APNS.KeyPath != "" {
		keyPEMBlock, err := ioutil.ReadFile(config.APNS.KeyPath)
		if err != nil {
			return err
		}
		config.APNS.Key = string(keyPEMBlock)
	}

	return nil
}
