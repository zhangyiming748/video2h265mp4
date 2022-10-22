package video2h265mp4

import (
	"fmt"
	"github.com/zhangyiming748/replace"
	"github.com/zhangyiming748/video2h265mp4/log"
	"github.com/zhangyiming748/video2h265mp4/util"
	"github.com/zhangyiming748/voiceAlert"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

const (
	success  = iota + 1 // 单次转码成功
	failed              // 转码失败,程序退出
	complete            // 转码进程完成
)
const (
	Byte     = 1
	Kilobyte = 1000
	Megabyte = 1000 * 1000
	Gigabyte = 1000 * 1000 * 1000
	Tegabyte = 1000 * 1000 * 1000 * 1000
	Pegabyte = 1000 * 1000 * 1000 * 1000 * 1000
	Exgabyte = 1000 * 1000 * 1000 * 1000 * 1000 * 1000
	Zegabyte = 1000 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000
	Yogabyte = 1000 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000
)

/*
转换h265并返回此次任务节省的磁盘空间
*/
func ConvToH265(src, dst, pattern, threads string) string {
	var sum int64
	defer func() {
		if err := recover(); err != nil {
			voiceAlert.Voice(failed)
			log.Debug.Printf("程序此次运行产生的错误:%v\n", err)
		} else {
			voiceAlert.Voice(complete)
		}
	}()
	if util.Illegal(src, dst, threads) {
		os.Exit(1)
	}
	files := util.GetFiles(src, pattern)
	log.Info.Println("文件目录", files)
	l := len(files)
	for index, file := range files {
		runtime.GC()
		before, err := os.Stat(strings.Join([]string{src, file}, "/"))
		if err != nil {
			log.Debug.Printf("获取源文件基础数据产生的错误:%v\n", err)
		}
		before_size := before.Size()
		fulldst := toh265Help(src, dst, file, threads, index, l)
		after, err := os.Stat(fulldst)
		if err != nil {
			log.Debug.Printf("获取目标文件基础数据产生的错误:%v\n", err)
		}
		after_size := after.Size()
		diff := diff(before_size, after_size)
		sum += (before_size - after_size)
		log.Debug.Printf("原始文件:%v\t处理前大小:%v\n", file, getSize(before_size))
		log.Debug.Printf("生成文件:%v\t处理后大小:%v\n", fulldst, getSize(after_size))
		log.Debug.Printf("节省了%v的空间\n", diff)
		runtime.GC()
	}
	log.Debug.Printf("共节省了%v的空间\n", getSize(sum))
	return getSize(sum)
}

/*
处理文件并反回输出文件全路径
*/
func toh265Help(src, dst, file, threads string, index, total int) string {
	defer func() {
		if err := recover(); err != nil {
			voiceAlert.Voice(failed)
		} else {
			voiceAlert.Voice(success)
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
	if err := os.RemoveAll(in); err != nil {
		log.Debug.Printf("删除源文件失败:%v\n", err)
	} else {
		log.Debug.Printf("删除源文件:%s\n", in)
	}
	return out
}

/*
输入文件路径,返回容量差
*/
func diff(before, after int64) string {
	b := before - after
	s := getSize(b)
	return s
}

func getSize(size int64) string {
	fsize := float64(size)
	var report string
	if fsize < Kilobyte {
		bsize := fsize / Byte
		s := fmt.Sprintf("%f", bsize)
		report = strings.Join([]string{s, "Byte"}, "")
	} else if fsize < Megabyte {
		ksize := fsize / Kilobyte
		s := fmt.Sprintf("%.1f", ksize)
		report = strings.Join([]string{s, "Kilobyte"}, "")
	} else if fsize < Gigabyte {
		msize := fsize / Megabyte
		s := fmt.Sprintf("%.2f", msize)
		report = strings.Join([]string{s, "Megabyte"}, "")
	} else if fsize < Tegabyte {
		gsize := fsize / Gigabyte
		s := fmt.Sprintf("%.3f", gsize)
		report = strings.Join([]string{s, "Gigabyte"}, "")
	} else if fsize < Pegabyte {
		tsize := fsize / Tegabyte
		s := fmt.Sprintf("%.4f", tsize)
		report = strings.Join([]string{s, "Tegabyte"}, "")
	} else if fsize < Exgabyte {
		psize := fsize / Pegabyte
		s := fmt.Sprintf("%.5f", psize)
		report = strings.Join([]string{s, "Pegabyte"}, "")
	} else if fsize < Zegabyte {
		zsize := fsize / Exgabyte
		s := fmt.Sprintf("%.6f", zsize)
		report = strings.Join([]string{s, "Exgabyte"}, "")
	} else if fsize < Yogabyte {
		ysize := fsize / Zegabyte
		s := fmt.Sprintf("%.7f", ysize)
		report = strings.Join([]string{s, "Zegabyte"}, "")
	} else {
		more := fsize / Yogabyte
		s := fmt.Sprintf("%.8f", more)
		report = strings.Join([]string{s, "Yogabyte"}, "")
	}
	return report
}
