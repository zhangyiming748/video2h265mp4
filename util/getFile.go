package util

import (
	"github.com/zhangyiming748/video2h265mp4/log"
	"os"
	"path"
	"strings"
)

func GetMultiFiles(dir, pattern string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Debug.Printf("读取文件目录产生的错误:%v\n", err)
	}
	var aim []string
	if strings.Contains(pattern, ";") {
		exts := strings.Split(pattern, ";")
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				log.Info.Println("跳过隐藏文件", file.Name())
				continue
			}
			ext := path.Ext(file.Name())
			log.Info.Printf("extname is %v\n", ext)
			for _, ex := range exts {
				if strings.Contains(ext, ex) {
					aim = append(aim, file.Name())
				}
			}
		}
	} else {
		for _, file := range files {
			if strings.HasPrefix(file.Name(), ".") {
				log.Info.Println("跳过隐藏文件", file.Name())
				continue
			}
			ext := path.Ext(file.Name())
			log.Info.Printf("extname is %v\n", ext)
			if strings.Contains(ext, pattern) {
				aim = append(aim, file.Name())
			}
		}
	}
	log.Debug.Printf("有效的目标文件: %v \n", aim)
	return aim
}

func GetFiles(dir, pattern string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Debug.Printf("读取文件目录产生的错误:%v\n", err)
	}
	var aim []string
	types := strings.Split(pattern, ";") //"wmv;rm"
	for _, f := range files {
		if l := strings.Split(f.Name(), ".")[0]; len(l) != 0 {
			for _, v := range types {
				if strings.HasSuffix(f.Name(), v) {
					log.Debug.Printf("有效的目标文件:%v\n", f.Name())
					aim = append(aim, f.Name())
				}
			}
		}
	}
	return aim
}
