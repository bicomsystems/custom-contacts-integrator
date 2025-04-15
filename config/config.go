package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type DbConf struct {
	Address            string
	Port               int
	User               string
	Password           string
	Database           string
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifeTime    int
}

type HttpServerConf struct {
	Address string
	Port    int
}

type Configuration struct {
	MySQL      DbConf
	HttpServer HttpServerConf
}

var (
	Conf *Configuration = DefaultConf()
)

// default MySQL conf
func getMySqlConf() DbConf {
	return DbConf{
		Address:            "127.0.0.1",
		Port:               3306,
		User:               "root",
		Password:           "Test123!",
		Database:           "",
		MaxOpenConnections: 10,
		MaxIdleConnections: 2,
		ConnMaxLifeTime:    10,
	}
}

// default HTTP server conf (it should always retain localhost address but the port could be changed if needed. By default it listens to port 5555
func getHttpServerConf() HttpServerConf {
	return HttpServerConf{
		Address: "127.0.0.1",
		Port:    5555,
	}
}

// DefaultConfiguration for app
func DefaultConf() *Configuration {
	return &Configuration{
		MySQL:      getMySqlConf(),
		HttpServer: getHttpServerConf(),
	}
}

// InitUsingFile decodes filename on the system (e.x csm.toml) and initalizes Conf variable with appropriate fields from the file
// same as if function DefaultConf is called
func InitUsingFile(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		return err
	}

	if _, err := toml.DecodeFile(filename, Conf); err != nil {
		return err
	}

	return nil
}
