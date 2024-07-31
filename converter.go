package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tealeg/xlsx"
)

// MetaData 元数据
// 从文件中提取的数据结构信息
type MetaData struct {
	Path string    // 路径
	Name string    // 文件名称
	Cols []*Column // 字段
}

func (m *MetaData) String() string {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("绝对路径：%s, 文件名：%s\n", m.Path, m.Name))
	for _, v := range m.Cols {
		buf.WriteString(fmt.Sprintf("%v\t", v.Name))
	}
	if len(m.Cols) > 0 {
		buf.WriteString("\n")
	}
	for _, v := range m.Cols {
		buf.WriteString(fmt.Sprintf("%v\t", v.Comment))
	}
	if len(m.Cols) > 0 {
		buf.WriteString("\n")
	}
	return buf.String()
}

// Column 字段信息
type Column struct {
	Index   int
	Name    string // 字段名称
	Kind    int    // 数据类型
	IsArray bool   // 是否数组
	Comment string // 描述
}

func GetMetaData(xlsxFile *xlsx.File) (m *MetaData, err error) {
	m = &MetaData{}

	if len(xlsxFile.Sheets) == 0 {
		return
	}

	rows := xlsxFile.Sheets[0].Rows
	if len(rows) == 0 {
		return
	}

	// 字段信息
	for k, v := range rows[0].Cells {
		name := v.String()
		l := len(name)
		if l <= 3 || name[l-1] != ')' {
			continue
		}
		i := strings.LastIndex(name, "(")
		if i <= 0 {
			continue
		}
		kind := XlsxToKind[name[i:]]
		if kind == Invalid {
			continue
		}
		m.Cols = append(m.Cols, &Column{
			Index:   k,
			Name:    name[:i],
			Kind:    kind,
			IsArray: strings.HasPrefix(name[i:], "(arr"),
		})
	}
	if len(rows) > 1 {
		for _, v := range m.Cols {
			v.Comment = rows[1].Cells[v.Index].String()
		}
	}
	return
}

//func CreateJsonFile(meta *MetaData, rows []*xlsx.Row) {
//	data := struct {
//		Cols []*Column
//		Data [][]interface{}
//	}{
//		Cols: meta.Cols,
//	}
//
//	for i, v := range rows {
//		var r []interface{}
//		for _, c := range meta.Cols {
//			var cell *xlsx.Cell
//			if c.Index < len(v.Cells) {
//				cell = v.Cells[c.Index]
//			}
//			var val interface{}
//			switch c.Kind {
//			case KindBool:
//				if cell == nil {
//					val = false
//				} else {
//					val = cell.Bool()
//				}
//
//			case KindString:
//				if cell == nil {
//					val = ""
//				} else {
//					val = fmt.Sprintf("%q", cell.String())
//				}
//
//			case KindInt32, KindInt64, KindInt, KindFloat32, KindFloat64, KindDouble:
//				if cell == nil {
//					val = 0
//				} else {
//					val = cell.Value
//				}
//
//			case KindArrString:
//				if cell == nil {
//					val = "[]"
//				} else {
//					s := strings.Join(strings.Split(cell.Value, ","), "\",\"")
//					if s != "" {
//						val = fmt.Sprintf("[\"%s\"]", s)
//					} else {
//						val = "[]"
//					}
//				}
//
//			case KindArrBool:
//				if cell == nil {
//					val = "[]"
//				} else {
//					var bs []string
//					for _, v := range strings.Split(cell.Value, ",") {
//						b, err := strconv.ParseBool(v)
//						if err != nil {
//							log.Printf("create json file error:%v column:%v index:%d filepath:%s\n",
//								err, c.Name, i+2, meta.Path)
//						}
//						bs = append(bs, strconv.FormatBool(b))
//					}
//					val = fmt.Sprintf("[%s]", strings.Join(bs, ","))
//				}
//
//			case KindArrInt32, KindArrInt64, KindArrInt,
//				KindArrFloat32, KindArrFloat64, KindArrDouble:
//				if cell == nil {
//					val = "[]"
//				} else {
//					val = fmt.Sprintf("[%s]", cell.String())
//				}
//			case KindMap:
//				if cell == nil {
//					val = "{}"
//				} else {
//					buf := bytes.NewBuffer(nil)
//					arrKV := strings.Split(cell.String(), Cfg.Mapseparate)
//					for k, v := range arrKV {
//						arr := strings.Split(v, ",")
//						if len(arr) != 2 {
//							continue
//						}
//						if k < len(arrKV)-1 {
//							buf.WriteString(fmt.Sprintf("\"%v\":%v,", arr[0], arr[1]))
//						} else {
//							buf.WriteString(fmt.Sprintf("\"%v\":%v", arr[0], arr[1]))
//						}
//					}
//					if buf.Len() > 0 {
//						val = fmt.Sprintf("{%s}", buf.String())
//					} else {
//						val = "{}"
//					}
//				}
//			}
//			r = append(r, val)
//		}
//		data.Data = append(data.Data, r)
//	}
//
//	func() {
//		f, err := os.Create(filepath.Join(Cfg.JsonDir, fmt.Sprint(meta.Name, ".json")))
//		if err != nil {
//			log.Printf("create json file error:%v filepath:%s\n", err, meta.Path)
//			return
//		}
//		defer f.Close()
//		if err = temps.ExecuteTemplate(f, "json", data); err != nil {
//			log.Printf("generate json error:%v filepath:%s\n", err, meta.Path)
//			return
//		}
//	}()
//}

func CreateGoManager(paths []string) {
	if len(paths) == 0 {
		return
	}
	p := filepath.Join(Cfg.GostructDir, "manager.go")
	f, err := os.Create(p)
	if err != nil {
		log.Printf("create manager.go file error:%v filepath:%s\n", err, p)
		return
	}
	defer f.Close()
	if err = temps.ExecuteTemplate(f, "gomanager", struct {
		PkgName string
		Paths   []string
	}{
		filepath.Base(Cfg.GostructDir), paths,
	}); err != nil {
		log.Printf("generate manager.go error:%v filepath:%s\n", err, p)
		return
	}
}
