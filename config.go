package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var Cfg = new(Config)

type Config struct {
	ExcelDir string // Excel所在目录,默认./excel
	//JsonDir      string // 生成json文件的目录，默认./json
	//ProtoDir     string // 生成protobuf序列化文件的目录，默认./proto
	ProtocolDir  string // 生成proto文件的目录，默认./protocol
	GostructDir  string // 生成go代码文件的目录，默认./gostruct
	Mapseparate  string // map类型字段分隔符，默认为“;”
	GoPackage    string // protobuf go_package 参数
	ProtoPackage string // protobuf package 参数
	Protoc       string // protoc 可执行文件路径
	ProtocGenGo  string // protoc-gen-go 可执行文件路径
	ProtocGoOut  string // protoc-gen-go 参数
	IntToInt32   bool   // int对应int32
}

func (c *Config) Init() {
	if c.ExcelDir == "" {
		c.ExcelDir = "excel"
	}
	if err := NewDir(c.ExcelDir); err != nil {
		log.Fatalln("new exceldir error:", err)
	}
	//if c.JsonDir == "" {
	//	c.JsonDir = "json"
	//}
	//if err := NewDir(c.JsonDir); err != nil {
	//	log.Fatalln("new jsondir error:", err)
	//}
	//if c.ProtoDir == "" {
	//	c.ProtoDir = "proto"
	//}
	//if err := NewDir(c.ProtoDir); err != nil {
	//	log.Fatalln("new protodir error:", err)
	//}
	if c.ProtocolDir == "" {
		c.ProtocolDir = "protocol"
	}
	if err := NewDir(c.ProtocolDir); err != nil {
		log.Fatalln("new protocoldir error:", err)
	}
	if c.GostructDir == "" {
		c.GostructDir = "gostruct"
	}
	if err := NewDir(c.GostructDir); err != nil {
		log.Fatalln("new gostructdir error:", err)
	}
	if c.Mapseparate == "" {
		c.Mapseparate = DefaultMapSeparate
	}
	if c.GoPackage == "" {
		c.GoPackage = fmt.Sprintf("/%v", c.GostructDir)
	}
	if c.ProtoPackage == "" {
		c.ProtoPackage = c.GostructDir
	}
	if c.Protoc == "" {
		c.Protoc = "./protoc"
	}
	if c.ProtocGenGo == "" {
		c.ProtocGenGo = "./protoc-gen-go"
	}
	if c.ProtocGoOut == "" {
		c.ProtocGoOut = "."
	}

	log.Println(fmt.Sprintf("%#v", *c))
}

var paths = []string{
	"/etc/converter",
	"$HOME/.converter",
	".",
}

func GetViper(name string) *viper.Viper {
	vp := viper.New()
	vp.SetConfigName(name)
	vp.SetConfigType("json")
	for _, v := range paths {
		vp.AddConfigPath(v)
	}

	err := vp.ReadInConfig()
	if err != nil {
		vp = viper.New()
		vp.SetConfigName(name)
		vp.SetConfigType("yaml")
		for _, v := range paths {
			vp.AddConfigPath(v)
		}
		if err = vp.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
	return vp
}

func InitConfig() {
	err := GetViper("config").Unmarshal(Cfg)
	if err != nil {
		panic(err)
	}
	Cfg.Init()
}
