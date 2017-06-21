package cams

import (
	"testing"
	"os"
	"math"
	"github.com/hschendel/stl"
	"reflect"
	"path/filepath"
)

func Test_CreateCams(t *testing.T) {
	inputCsv := filepath.Join("testfiles", "star_path.csv")
	leftOutputFile := filepath.Join("testfiles", "left.stl")
	rightOutputFile := filepath.Join("testfiles", "rigth.stl")
	leftCorrectFile := filepath.Join("testfiles", "correct_star_left.stl")
	rightCorrectFile := filepath.Join("testfiles", "correct_star_right.stl")

	//NOTE: My examples output models 1000X larger because Trimble Sketchup has rounding bugs with small objects
	err := CreateCams(
		5,
		0.045,
		0.045,
		Coordinate{-21750, -30536},
		Coordinate{21750, -30536},
		7060,
		inputCsv,
		leftOutputFile,
		rightOutputFile,
	)
	if err != nil {
		t.Error(err)
	}

	if same, sameerr := areStlsTheSame(leftCorrectFile, leftOutputFile);
		!same || sameerr != nil {
		t.Error("Left cam wrong")
	}
	if same, sameerr := areStlsTheSame(rightCorrectFile, rightOutputFile);
		!same || sameerr != nil {
		t.Error("Right cam wrong")
	}

	if _, fileErr := os.Stat(leftOutputFile); os.IsNotExist(fileErr) {
		t.Error("Did not write file")
	}
	os.Remove(leftOutputFile)
	if _, fileErr := os.Stat(rightOutputFile); os.IsNotExist(fileErr) {
		t.Error("Did not write file")
	}
	os.Remove(rightOutputFile)
}

func Test_Scale(t *testing.T) {
	input := []Coordinate{
		Coordinate{ 1.0, 2.0 },
		Coordinate{ 4.0, 1.0 },
	}
	expected := []Coordinate{
		Coordinate{ 10.0, 4.0 },
		Coordinate{ 40.0, 2.0 },
	}
	result := Scale(input, 10.0, 2.0)
	if !reflect.DeepEqual(expected, result) {
		t.Error("Did not scale properly")
	}
}

func Test_Interpolate(t *testing.T) {
	input := []Coordinate{
		Coordinate{ 10.0, 100.0 },
		Coordinate{ 20, 200.0 },
	}
	expected := []Coordinate{
		Coordinate{ 10.0, 100.0 },
		Coordinate{ 15.0, 150.0 },
		Coordinate{ 20, 200.0 },
		Coordinate{ 15.0, 150.0 },
	}
	result := Interpolate(input, 2)
	if !reflect.DeepEqual(expected, result) {
		t.Error("Did not interpolate properly")
	}
}

func Test_WriteRealCam(t *testing.T) {
	outputFile := filepath.Join("testfiles", "real_cam.stl")
	inputCsv := filepath.Join("testfiles", "star_cam_coordinates.csv")

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
	outputFile := filepath.Join("testfiles", "simple_cam.stl")

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

func Test_ConvertCoord(t *testing.T) {
	result := ConvertCoord(
		Coordinate{10.0, 11.0},
		Coordinate{0.0, 0.0},
		35,
		22,
		1.0,
		false)
	if math.Abs(result.X - -13.86168) > 0.01 {
		t.Error("Incorrect X coord")
	}
	if math.Abs(result.Y - -0.34859) > 0.01 {
		t.Error("Incorrect Y coord")
	}
}

func areStlsTheSame(a, b string) (bool, error) {
	solidA, errA := stl.ReadFile(a)
	if errA != nil {
		return false, errA
	}
	solidB, errB := stl.ReadFile(b)
	if errB != nil {
		return false, errB
	}

	return reflect.DeepEqual(solidA.Triangles, solidB.Triangles), nil
}
