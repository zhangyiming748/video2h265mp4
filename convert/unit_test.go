package convert

import "testing"

func TestUnit(t *testing.T) {
	src := "/Users/zen/Github/video2h265mp4/DB"
	dst := "/Users/zen/Github/video2h265mp4/DB/h265"
	pattern := "mp4"
	threads := 4
	ConvToH265(src, dst, pattern, threads)
}
