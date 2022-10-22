package util

import "testing"

func TestGetMultiFiles(t *testing.T) {
	GetMultiFiles("/Users/zen/Github/video2h265mp4", "log")
	GetMultiFiles("/Users/zen/Github/video2h265mp4", "go;log")
}
