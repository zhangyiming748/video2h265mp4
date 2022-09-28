package video2h265mp4

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestUnit(t *testing.T) {
	src := "/Users/zen/Github/video2h265mp4/DB"
	dst := "/Users/zen/Github/video2h265mp4/DB/h265"
	pattern := "mp4"
	threads := "4"
	ConvToH265(src, dst, pattern, threads)
}

func TestVoice(t *testing.T) {
	voice(1)
	voice(2)
	voice(3)
}

func BenchmarkBeep(b *testing.B) {
	var cmd *exec.Cmd
	cmd = exec.Command("echo", "-e", "\\a")
	for i := 0; i < b.N; i++ {
		cmd.Run()
	}
}

func TestBeep(t *testing.T) {
	var cmd *exec.Cmd
	cmd = exec.Command("echo", "-e", "\\a")
	fmt.Println(cmd)
	for i := 0; i < 10; i++ {
		cmd.Run()
	}
}
