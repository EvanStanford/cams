package main

import (
	"os"
	"path/filepath"
	"github.com/EvanStanford/cams/profiler"
)

func basename (filename string) string {
	ext := filepath.Ext(filename)
	return filename[0:len(filename) - len(ext)]
}

func main() {
	args := os.Args[1:]

	inputCsv := args[0]

	base := basename(filepath.Base(inputCsv))

	outputStlL := filepath.Join("out", base, base + "-L.stl")
	outputStlR := filepath.Join("out", base, base + "-R.stl")

	os.MkdirAll(filepath.Join("out", base), 0755)

	cams.CreateCams(
		5,
		0.045,
		0.045,
		cams.Coordinate{-21750, -30536},
		cams.Coordinate{21750, -30536},
		7060,
		inputCsv,
		outputStlL,
		outputStlR)
}