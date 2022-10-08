package video2h265mp4

import (
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func TestUnit(t *testing.T) {
	src := "/Users/zen/Github/video2h265mp4/DB"
	dst := "/Users/zen/Github/video2h265mp4/DB/h265"
	pattern := "mp4"
	threads := "4"
	ConvToH265(src, dst, pattern, threads)
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
func TestReadfile(t *testing.T) {
	if fi, err := os.Stat("/Users/zen/Downloads/Docker.dmg"); err != nil {

	} else {
		t.Logf("name:%v\n", fi.Name())
		t.Logf("size:%v\n", fi.Size()/Megabyte)
		t.Logf("is dir:%v\n", fi.IsDir())
		t.Logf("mode:%v\n", fi.Mode())
		t.Logf("modTime:%v\n", fi.ModTime())
	}
	ret := getSize("/Users/zen/Downloads/Docker.dmg")
	t.Logf("ret==%s\n", ret)
}
func TestGetSize(t *testing.T) {
	src := "/Users/zen/Github/video2h265mp4/DB/h265/不同丹数丝袜之间的差别.mp4"
	ret := getSize(src)
	t.Log(ret)
}
