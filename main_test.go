package video2h265mp4

import (
	"github.com/zhangyiming748/video2h265mp4/log"
	"os/exec"
	"strconv"
	"testing"
)

func TestUnit(t *testing.T) {
	src := "/Users/zen/Github/video2h265mp4/DB"
	dst := "/Users/zen/Github/video2h265mp4/DB/h265"
	pattern := "mp4"
	threads := "4"
	save, total := ConvToH265(src, dst, pattern, threads)
	log.Debug.Printf("节省的空间:%v\n", save)
	log.Debug.Printf("共处理的文件数:%v\n", total)

}

func BenchmarkBeep(b *testing.B) {
	var cmd *exec.Cmd
	cmd = exec.Command("echo", "-e", "\\a")
	for i := 0; i < b.N; i++ {
		cmd.Run()
	}
}

func TestFakeThreads(t *testing.T) {
	thread := "13"
	ret := fakeThreads(thread)
	t.Log(ret)

}
func fakeThreads(threads string) bool {
	maxThreads := 12
	if t, err := strconv.Atoi(threads); err != nil {
		return false
	} else if t >= maxThreads {
		return false
	} else {
		return true
	}
}

func TestGetSize(t *testing.T) {

}
