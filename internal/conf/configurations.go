package conf

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type MapleDatabaseConfigurations struct {
	Url string `required:"true"`
}

type MapleStorageConfigurations struct {
	ImagesPath string `required:"true" split_words:"true"`
}

type MapleLoggingConfigurations struct {
	Debug    bool
	Requests bool
}

type MapleAPIConfigurations struct {
	Host            string `default:"0.0.0.0:3000"`
	AdditionalHosts []string
}

type MapleAuthConfigurations struct {
	KeysEndpoint string `split_words:"true"`
}

type MapleSurgeConfigurations struct {
	URL        string
	ServiceKey string `split_words:"true"`
}

type MapleConfigurations struct {
	Database MapleDatabaseConfigurations `required:"true"`
	Storage  MapleStorageConfigurations  `required:"true"`
	Logging  MapleLoggingConfigurations
	API      MapleAPIConfigurations
	Auth     MapleAuthConfigurations
	Surge    MapleSurgeConfigurations

	Production bool `default:"false"`
}

func LoadFromEnvironments() (*MapleConfigurations, error) {
	// Load .env
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Warnln("Failed to load .env")
	}

	config := new(MapleConfigurations)

	if err := envconfig.Process("maple", config); err != nil {
		return nil, err
	}

	err := config.ApplyDefaults()
	if err != nil {
		return nil, err
	}

	if err = config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// ApplyDefaults apply defaults to MapleConfigurations
func (c *MapleConfigurations) ApplyDefaults() error {
	if c.API.AdditionalHosts == nil {
		c.API.AdditionalHosts = make([]string, 0)
	}

	return nil
}

func (c *MapleConfigurations) Validate() error {
	return nil
}
