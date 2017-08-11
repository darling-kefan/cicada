package cicada

import (
	"fmt"
	"log"
	"os"
	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Mysql     map[string]mysqlDb  `toml:"mysql"`
	InfluxDB  influxDB            `toml:"influxDB"`
}

type mysqlDb struct {
	Host      string  `toml:"host"`
	Port      int     `toml:"port"`
	Username  string  `toml:"username"`
	Password  string  `toml:"password"`
	Database  string  `toml:"database"`
}

type influxDB struct {
	Server    string  `toml:"server"`
	Username  string  `toml:"username"`
	Password  string  `toml:"password"`
}

var conf tomlConfig


func init1() {
	tomlFile := fmt.Sprintf("%s%cconfig.toml", runPath(), os.PathSeparator)
	// @todo fix config file
	tomlFile = "/home/shouqiang/go/src/github.com/darling-kefan/cicada/config.toml"
	if _, err := toml.DecodeFile(tomlFile, &conf); err != nil {
		log.Fatal(err)
	}
}

