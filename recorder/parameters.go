package recorder

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/lino-network/lino-service-discovery/pkg/k8s"
	"github.com/rs/zerolog/log"
	yaml "gopkg.in/yaml.v2"
)

type configKey = string

// Configs contains all config strings.
type Configs struct {
	store        map[configKey]string
	serviceAddrs k8s.ServiceAddr
}

func readKeys(c *Configs, path string, store []configKey) error {
	if path != "" {
		yamlFile, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		m := make(map[string]string)
		err = yaml.Unmarshal(yamlFile, &m)
		if err != nil {
			return err
		}
		for _, key := range store {
			v, ok := m[key]
			if !ok {
				v = defaultConfig[key]
			}
			c.store[key] = v
		}
	} else {
		for _, key := range store {
			value := os.Getenv(strings.ToUpper(key))
			if value == "" {
				value = defaultConfig[key]
			}
			c.store[key] = value
		}
	}
	return nil
}

func NewConfig(configPath string) (*Configs, error) {
	var c = &Configs{}
	c.store = make(map[configKey]string)
	if err := readKeys(c, configPath, configKeys); err != nil {
		return nil, err
	}
	c.serviceAddrs = k8s.NewServiceAddrs(k8s.DefaultVals{})
	log.Debug().Msgf("Configs:%+v", c.store)
	return c, nil
}

const (
	dbUsernameKey configKey = "db_username"
	dbPasswordKey configKey = "db_password"
	dbHostKey     configKey = "db_host"
	dbPortKey     configKey = "db_port"
	dbNameKey     configKey = "db_name"
)

var (
	configKeys = []configKey{dbUsernameKey, dbHostKey, dbPasswordKey, dbPortKey,
		dbNameKey}
	defaultConfig = map[configKey]string{
		dbHostKey:     "localhost",
		dbUsernameKey: "root",
		dbPasswordKey: "my-secret",
		dbPortKey:     "3308",
		dbNameKey:     "lino_db",
	}
)

// DBUsername returns MySQL username
func (c *Configs) DBUsername() string {
	return c.store[dbUsernameKey]
}

// DBPassword returns MySQL password
func (c *Configs) DBPassword() string {
	return c.store[dbPasswordKey]
}

// // DBHost returns MySQL host
func (c *Configs) DBHost() string {
	if os.Getenv("GO_ENV") == "" {
		return c.store[dbHostKey]
	}
	return c.serviceAddrs.GetAddr(k8s.RdsRecorder)
}

// DBPort returns MySQL port; default to 3306
func (c *Configs) DBPort() int {
	if port, ok := c.store[dbPortKey]; ok {
		if port, err := strconv.Atoi(port); err == nil {
			return port
		}
	}
	return 3306
}

// DBName returns MySQL db name
func (c *Configs) DBName() string {
	return c.store[dbNameKey]
}
