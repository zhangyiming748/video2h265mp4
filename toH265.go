package video2h265mp4

import (
	"github.com/zhangyiming748/video2h265mp4/log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

func ConvToH265(src, dst, pattern, threads string) {
	files := getFiles(src, pattern)
	log.Info.Println("文件目录", files)
	l := len(files)
	for index, file := range files {
		runtime.GC()
		toh265Help(src, dst, file, threads, index, l)
		runtime.GC()
	}

}
func toh265Help(src, dst, file, threads string, index, total int) {
	in := strings.Join([]string{src, file}, "/")
	log.Debug.Printf("开始处理文件:%v", in)
	extname := path.Ext(file)
	filename := strings.Trim(file, extname)
	filename = replace(filename)
	newFilename := strings.Join([]string{filename, "mp4"}, ".")
	out := strings.Join([]string{dst, newFilename}, "/")

	log.Info.Println("源文件目录:", src)
	log.Info.Println("输出文件目录:", dst)
	log.Info.Println("开始处理文件:", in)
	log.Info.Println("输出文件:", out)

	cmd := exec.Command("ffmpeg", "-threads", threads, "-i", in, "-c:v", "libx265", "-threads", threads, out)
	log.Debug.Printf("生成的命令是:%s\n", cmd)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		log.Debug.Panicf("cmd.StdoutPipe产生的错误:%v", err)
	}
	if err = cmd.Start(); err != nil {
		log.Debug.Panicf("cmd.Run产生的错误:%v", err)
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		//写成输出日志
		log.Info.Printf("正在处理第 %d/%d 个文件: %s\n", index+1, total, file)
		t := string(tmp)
		t = strings.Replace(t, "\u0000", "", -1)
		log.Info.Println(t)
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		log.Debug.Panicf("命令执行中有错误产生:%v", err)
	}
	log.Debug.Printf("完成当前文件的处理:源文件是%s\t目标文件是%s\n", in, file)
	err = os.RemoveAll(in)
	if err != nil {
		return
	} else {
		log.Debug.Printf("删除源文件:%s\n", in)
	}
}

func getFiles(dir, pattern string) []string {
	files, _ := os.ReadDir(dir)
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
func replace(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "《", "", -1)
	str = strings.Replace(str, "》", "", -1)
	str = strings.Replace(str, "【", "", -1)
	str = strings.Replace(str, "】", "", -1)
	str = strings.Replace(str, "(", "", -1)
	str = strings.Replace(str, "+", "", -1)
	str = strings.Replace(str, ")", "", -1)
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\u00A0", "", -1)
	str = strings.Replace(str, "\u0000", "", -1)
	return str
}
