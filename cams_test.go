package cams

import (
	"testing"
	"os"
)

func Test_WriteRealCam(t *testing.T) {
	outputFile := "testfiles\\real_cam.stl"
	inputCsv := "testfiles\\star_cam_coordinates.csv"

	coords, err := ReadCoordsCsv(inputCsv)
	if err != nil {
		t.Error(err)
	}

	err = WriteCam(coords, outputFile)
	if err != nil {
		t.Error(err)
	}
	if _, fileErr := os.Stat(outputFile); os.IsNotExist(fileErr) {
		t.Error("Did not write file")
	}
	os.Remove(outputFile)
}

func Test_WriteSimpleCam(t *testing.T) {
	outputFile := "testfiles\\simple_cam.stl"

	var input = []Coordinate{
		Coordinate{0,5},
		Coordinate{6,6},
		Coordinate{0,10},
		Coordinate{-3,-3},
		Coordinate{0,-8},
		Coordinate{-6,-7},
		Coordinate{-5,0},
		Coordinate{4,4},
	}
	
	err := WriteCam(input, outputFile)
	if err != nil {
		t.Error(err)
	}
	if _, fileErr := os.Stat(outputFile); os.IsNotExist(fileErr) {
		t.Error("Did not write file")
	}
	os.Remove(outputFile)	
}



