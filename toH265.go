package video2h265mp4

import (
	"github.com/zhangyiming748/replace"
	"github.com/zhangyiming748/video2h265mp4/log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
)

const (
	success  = iota + 1 // 单次转码成功
	failed              // 转码失败,程序退出
	complete            // 转码进程完成
)

func ConvToH265(src, dst, pattern, threads string) {
	defer func() {
		if err := recover(); err != nil {
			voice(failed)
			log.Debug.Printf("程序此次运行产生的错误:%v\n", err)
		} else {
			voice(complete)
		}
	}()
	if illegal(src, dst, threads) {
		os.Exit(1)
	}
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
	defer func() {
		if err := recover(); err != nil {
			voice(failed)
		} else {
			voice(success)
		}
	}()
	in := strings.Join([]string{src, file}, "/")
	log.Debug.Printf("开始处理文件:%v\n", in)
	extname := path.Ext(file)
	filename := strings.Trim(file, extname)
	filename = replace.Replace(filename)
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
		log.Debug.Panicf("cmd.StdoutPipe产生的错误:%v\n", err)
	}
	if err = cmd.Start(); err != nil {
		log.Debug.Panicf("cmd.Run产生的错误:%v\n", err)
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
		log.Debug.Panicf("命令执行中有错误产生:%v\n", err)
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

//func replace(str string) string {
//	str = strings.Replace(str, "\n", "", -1)
//	str = strings.Replace(str, " ", "", -1)
//	str = strings.Replace(str, "《", "", -1)
//	str = strings.Replace(str, "》", "", -1)
//	str = strings.Replace(str, "【", "", -1)
//	str = strings.Replace(str, "】", "", -1)
//	str = strings.Replace(str, "(", "", -1)
//	str = strings.Replace(str, "+", "", -1)
//	str = strings.Replace(str, ")", "", -1)
//	str = strings.Replace(str, "`", "", -1)
//	str = strings.Replace(str, " ", "", -1)
//	str = strings.Replace(str, "\u00A0", "", -1)
//	str = strings.Replace(str, "\u0000", "", -1)
//	return str
//}

func voice(msg int) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		// 查询发音人 `say -v ?`
		switch msg {
		case success:
			//cmd = exec.Command("say", "-v", "Kate", "Rocket was launched successfully")
			cmd = exec.Command("say", "-v", "Victoria", "Rocket was launched successfully")
			cmd.Start()
		case failed:
			//cmd = exec.Command("say", "-v", "Bad News", "Rocket launch failed")
			cmd = exec.Command("say", "-v", "Victoria", "Rocket launch failed")
			cmd.Start()
		case complete:
			//cmd = exec.Command("say", "-v", "Kate", "mission complete!")
			cmd = exec.Command("say", "-v", "Victoria", "mission complete!")
			cmd.Start()
		}
	case "linux":
		cmd = exec.Command("echo", "-e", "\\a")
		switch msg {
		case success:
			for i := 0; i < 2; i++ {
				cmd.Start()
			}
		case failed:
			for i := 0; i < 50; i++ {
				cmd.Start()
			}
		case complete:
			for i := 0; i < 100; i++ {
				cmd.Start()
			}
		}
	}
	cmd.Wait()
}
func illegal(src, dst, threads string) bool {
	if src == dst {
		log.Debug.Println("输入输出目录相同")
		return true
	}
	if !exists(src) {
		log.Debug.Println("src目录不存在")
		return true
	}
	if !exists(dst) {
		log.Debug.Println("dst目录不存在")
		return true
	}
	if !isDir(src) {
		log.Debug.Println("src不是目录")
		return true
	}
	if !isDir(dst) {
		log.Debug.Println("dst不是目录")
		return true
	}
	if !allowThreads(threads) {
		log.Debug.Println("不允许的线程数")
		return true
	}
	return false
}

// 判断所给路径文件/文件夹是否存在
func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断给定线程数是否合法
func allowThreads(threads string) bool {
	maxThreads := runtime.NumCPU()
	if t, err := strconv.Atoi(threads); err != nil {
		return false
	} else if t >= maxThreads {
		return false
	} else {
		return true
	}
}
