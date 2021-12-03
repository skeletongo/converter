package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/tealeg/xlsx"
)

//go:embed templates
var templates embed.FS
var temps = template.Must(template.New("").Funcs(FM).ParseFS(templates, "templates/*"))

func main() {
	// 读取配置文件
	Load("config.json")

	// 查找所有xlsx文件路径
	paths, err := FindXlsxPath(Cfg.ExcelDir)
	if err != nil {
		log.Fatalf("find xlsx file error:%v workdir:%s\n", err, Cfg.ExcelDir)
	}

	g := sync.WaitGroup{}
	g.Add(len(paths) + 1)
	go func() {
		defer g.Done()
		CreateGoManager(paths)
	}()
	for _, v := range paths {
		go func(path string) {
			defer g.Done()
			Run(path)
		}(v)
	}
	g.Wait()
	log.Println("success")
}

func Run(path string) {
	xlsxFile, err := xlsx.OpenFile(path)
	if err != nil {
		log.Printf("open xlsx file error:%v filepath:%v\n", err, path)
		return
	}
	// 获取表结构
	meta, err := GetMetaData(xlsxFile)
	if err != nil {
		log.Printf("get xlsx file meta error:%v filepath:%s\n", err, path)
		return
	}
	meta.Path = path
	file, err := os.Stat(path)
	if err != nil {
		log.Printf("os.Stat error:%v filepath:%s\n", err, path)
		return
	}
	meta.Name = strings.TrimSuffix(file.Name(), ".xlsx")

	// 生成json文件
	if len(xlsxFile.Sheets) > 0 {
		rows := xlsxFile.Sheets[0].Rows
		if len(rows) > 2 {
			CreateJsonFile(meta, rows[2:])
		}
	}

	// 生成protobuf文件
	func() {
		f, err := os.Create(filepath.Join(Cfg.ProtocolDir, fmt.Sprint(strings.ToLower(meta.Name), ".proto")))
		if err != nil {
			log.Printf("create protobuf file error:%v filepath:%s\n", err, path)
			return
		}
		defer f.Close()
		if err = temps.ExecuteTemplate(f, "protocol", struct {
			PkgName string
			Meta    *MetaData
		}{
			filepath.Base(Cfg.ProtocolDir), meta,
		}); err != nil {
			log.Printf("generate protobuf error:%v filepath:%s\n", err, path)
			return
		}
	}()

	// 生成go代码文件
	func() {
		f, err := os.Create(filepath.Join(Cfg.GostructDir, fmt.Sprint(strings.ToLower(meta.Name), ".go")))
		if err != nil {
			log.Printf("create gostruct file error:%v filepath:%s\n", err, path)
			return
		}
		defer f.Close()
		if err = temps.ExecuteTemplate(f, "gostruct", struct {
			PkgName string
			Meta    *MetaData
		}{
			filepath.Base(Cfg.GostructDir), meta,
		}); err != nil {
			log.Printf("generate gostruct error:%v filepath:%s\n", err, path)
			return
		}
	}()
}
