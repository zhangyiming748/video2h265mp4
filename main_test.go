package video2h265mp4

import (
	. "github.com/zhangyiming748/video2h265mp4/convert"
	"testing"
)

func TestUnit(t *testing.T) {
	src := "/Users/zen/Github/video2h265mp4/DB"
	dst := "/Users/zen/Github/video2h265mp4/DB/h265"
	pattern := "mp4"
	threads := 4
	ConvertToH265(src, dst, pattern, threads)
}
