package main

import (
	"os"
	"path/filepath"
	"strings"
)

const DefaultMapSeparate = ";"

// 内部数据类型
const (
	Invalid = iota
	KindBool
	KindString
	KindInt32
	KindInt64
	KindInt
	KindFloat32
	KindFloat64
	KindDouble
	KindArrBool
	KindArrString
	KindArrInt32
	KindArrInt64
	KindArrInt
	KindArrFloat32
	KindArrFloat64
	KindArrDouble
	KindMap
)

// XlsxToKind xlsx数据类型和内部数据类型映射表
var XlsxToKind = map[string]int{
	"(bool)":       KindBool,
	"(string)":     KindString,
	"(str)":        KindString,
	"(int32)":      KindInt32,
	"(int64)":      KindInt64,
	"(int)":        KindInt,
	"(float32)":    KindFloat32,
	"(float64)":    KindFloat64,
	"(double)":     KindDouble,
	"(arrbool)":    KindArrBool,
	"(arrstring)":  KindArrString,
	"(arrstr)":     KindArrString,
	"(arrint32)":   KindArrInt32,
	"(arrint64)":   KindArrInt64,
	"(arrint)":     KindArrInt,
	"(arrfloat32)": KindArrFloat32,
	"(arrfloat64)": KindArrFloat64,
	"(arrdouble)":  KindArrDouble,
	"(map)":        KindMap,
}

// KindToProto 内部数据类型和protobuf数据类型映射表
var KindToProto = map[int]string{
	KindBool:       "bool",
	KindString:     "string",
	KindInt32:      "int32",
	KindInt64:      "int64",
	KindInt:        "int64",
	KindFloat32:    "float",
	KindFloat64:    "double",
	KindDouble:     "double",
	KindArrBool:    "bool",
	KindArrString:  "string",
	KindArrInt32:   "int32",
	KindArrInt64:   "int64",
	KindArrInt:     "int64",
	KindArrFloat32: "float",
	KindArrFloat64: "double",
	KindArrDouble:  "double",
	KindMap:        "map<int64, int64>",
}

func GetKindToProto(k int) string {
	switch k {
	case KindInt:
		if Cfg.IntToInt32 {
			return KindToProto[KindInt32]
		}
	case KindArrInt:
		if Cfg.IntToInt32 {
			return KindToProto[KindArrInt32]
		}
	default:
	}
	return KindToProto[k]
}

var KindToGo = map[int]string{
	KindBool:       "bool",
	KindString:     "string",
	KindInt32:      "int32",
	KindInt64:      "int64",
	KindInt:        "int64",
	KindFloat32:    "float32",
	KindFloat64:    "float64",
	KindDouble:     "float64",
	KindArrBool:    "[]bool",
	KindArrString:  "[]string",
	KindArrInt32:   "[]int32",
	KindArrInt64:   "[]int64",
	KindArrInt:     "[]int64",
	KindArrFloat32: "[]float32",
	KindArrFloat64: "[]float64",
	KindArrDouble:  "[]float64",
	KindMap:        "map[int64]int64",
}

func GetKindToGo(k int) string {
	switch k {
	case KindInt:
		if Cfg.IntToInt32 {
			return KindToGo[KindInt32]
		}
	case KindArrInt:
		if Cfg.IntToInt32 {
			return KindToGo[KindArrInt32]
		}
	default:
	}
	return KindToGo[k]
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func NewDir(path string) error {
	exist, err := PathExists(path)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	if err := os.Mkdir(path, 0666); err != nil {
		return err
	}
	return nil
}

func FindFilepath(dir string, f func(info os.FileInfo) bool) (infos []string, err error) {
	info, err := os.Stat(dir)
	if err != nil {
		return
	}
	if info.IsDir() {
		files, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(files); i++ {
			res, err := FindFilepath(filepath.Join(dir, files[i].Name()), f)
			if err != nil {
				return nil, err
			}
			infos = append(infos, res...)
		}
		return infos, nil
	}
	if f(info) {
		infos = append(infos, dir)
	}
	return
}

func FindXlsxPath(dir string) (infos []string, err error) {
	return FindFilepath(dir, func(info os.FileInfo) bool {
		return strings.HasSuffix(info.Name(), ".xlsx")
	})
}
