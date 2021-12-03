package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var Cfg = new(Config)

type Config struct {
	ExcelDir    string `json:"exceldir"`    // Excel所在目录,默认./excel
	JsonDir     string `json:"jsondir"`     // 生成json文件的目录，默认./json
	ProtocolDir string `json:"protocoldir"` // 生成proto文件的目录，默认./protocol
	GostructDir string `json:"gostructdir"` // 生成go代码文件的目录，默认./gostruct
}

func (c *Config) Init() {
	if c.ExcelDir == "" {
		c.ExcelDir = "excel"
		if err := NewDir(c.ExcelDir); err != nil {
			log.Fatalln("new exceldir error:", err)
		}
	}
	if c.JsonDir == "" {
		c.JsonDir = "json"
		if err := NewDir(c.JsonDir); err != nil {
			log.Fatalln("new jsondir error:", err)
		}
	}
	if c.ProtocolDir == "" {
		c.ProtocolDir = "protocol"
		if err := NewDir(c.ProtocolDir); err != nil {
			log.Fatalln("new protocoldir error:", err)
		}
	}
	if c.GostructDir == "" {
		c.GostructDir = "gostruct"
		if err := NewDir(c.GostructDir); err != nil {
			log.Fatalln("new gostructdir error:", err)
		}
	}
}

func Load(filepath string) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		log.Fatalln("load config error:", err)
	}
	if len(bytes) > 0 {
		if err = json.Unmarshal(bytes, Cfg); err != nil {
			log.Fatalln("config unmarshal error:", err)
		}
	}
	Cfg.Init()
}
