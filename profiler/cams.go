package cams

import (
	"github.com/hschendel/stl"
	"fmt"
	"encoding/csv"
	"os"
	"strconv"
	"math"
)

type Coordinate struct {
	X, Y float64
}

func CreateCams (numberOfInterpolations int, scaleX, scaleY float64,
		leftCamCenter, rightCamCenter Coordinate, rp float64,
		pathFile, leftOutput, rightOutput string) error {
	rawPath, err := ReadCoordsCsv(pathFile)
	if err != nil {
		return err
	}
	scaledPath := Scale(rawPath, scaleX, scaleY)
	interpolatedPath := Interpolate(scaledPath, numberOfInterpolations)

	left, right := GetCams(interpolatedPath, leftCamCenter, rightCamCenter, rp)

	err = WriteCam(left, leftOutput)
	if (err != nil) {
		return err
	}
	err = WriteCam(right, rightOutput)
	if (err != nil) {
		return err
	}
	return nil
}

//iteratively find where the cam centers should be to give the correct maxCamRadius
func GetCamCenters(path []Coordinate, rp, camBoreSpacing, maxCamRadius float64) (Coordinate, Coordinate) {
	minCamY := -1.0 * maxCamRadius
	maxCamY := 0.0
	var leftCenter, rightCenter Coordinate
	for i := 0; i < 10; i++ {
		testCamY := (minCamY + maxCamY)/2.0
		leftCenter = Coordinate{camBoreSpacing/-2.0, testCamY}
		rightCenter = Coordinate{camBoreSpacing/2.0, testCamY}
		
		if maxCamRadius < getMaxRadius(path, leftCenter, rightCenter, rp) {
			minCamY = testCamY
		} else {
			maxCamY = testCamY
		}
	}
	return leftCenter, rightCenter
}

func getMaxRadius(path []Coordinate, leftCamCenter, rightCamCenter Coordinate, rp float64) float64 {
	maxR := 0.0
	for _, target := range path {
		if r := getRadius(target, leftCamCenter, rp); r > maxR {
			maxR = r
		}
		if r := getRadius(target, rightCamCenter, rp); r > maxR {
			maxR = r
		}
	}
	return maxR
}

func Scale(path []Coordinate, sx, sy float64) []Coordinate {
	result := make([]Coordinate, len(path))
	for i, p := range path {
		result[i] = Coordinate{ p.X * sx, p.Y * sy}
	}
	return result
}

func Interpolate(path []Coordinate, n int) []Coordinate {
	result := make([]Coordinate, len(path) * n)
	for i, p := range path {
		nextp := path[(i + 1) % len(path)]
		dx := (nextp.X - p.X) / float64(n)
		dy := (nextp.Y - p.Y) / float64(n)

		for d := 0; d < n; d++ {
			result[i*n + d] = Coordinate {
				X: p.X + dx * float64(d),
				Y: p.Y + dy * float64(d),
			}
		}
	}
	return result
}

func GetCams (
	path []Coordinate, leftCamCenter, rightCamCenter Coordinate, rp float64) (
	left, right []Coordinate) {

	left = make([]Coordinate, len(path))
	right = make([]Coordinate, len(path))
	
	for i, target := range path {
		left[i] = ConvertCoord(
			target,
			leftCamCenter,
			len(path),
			i,
			rp,
			false)
		right[i] = ConvertCoord(
			target,
			rightCamCenter,
			len(path),
			i,
			rp,
			true)
	}
	return
}

//n- total number of coordinates
//i- index of this coordinate
//right- is this the right cam?
func ConvertCoord(target, camCenter Coordinate, n, i int, rp float64, right bool) Coordinate {
	rx := getRadius(target, camCenter, rp)

	B := math.Atan( (target.Y - camCenter.Y) / (target.X - camCenter.X) )

	Kx := 2.0 * math.Pi * float64(i) / float64(n)
	if right {
		Kx += math.Pi
		Kx *= -1.0
	}
	
	θ := 2.0 * math.Pi - Kx + B

	return Coordinate {
		X: rx * math.Cos(θ),
		Y: rx * math.Sin(θ),
	}
}

func getRadius(target, camCenter Coordinate, rp float64) float64 {
	return math.Sqrt(
		math.Pow(target.X - camCenter.X,2) +
		math.Pow(target.Y - camCenter.Y,2)) -
		rp
}

func WriteCam(path []Coordinate, outputFile string) error {

	pizza := make([]stl.Triangle, len(path) + 1)

	for i,_ := range path {
		a, b := i, (i+1)%len(path)
		pizza[i] = getPizzaSlice(path[a], path[b])
	}

	pizza[len(path)] = getDirectionalArrow(path[0], path[1])

	cam := stl.Solid {
		Name: "cam",
		IsAscii: true,
		Triangles: pizza,
	}

	err := cam.WriteFile(outputFile)
	if err != nil {
		fmt.Println("Fail on write")
		fmt.Println(err)
		return err
	}
	return nil
}

func ReadCoordsCsv(csvFile string) ([]Coordinate, error) {
	content, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(content)
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	coords := make([]Coordinate, len(records))
	for i, record := range records {
		X, xerr := strconv.ParseFloat(record[0], 64)
		if xerr != nil {
			return nil, xerr
		}
		Y, yerr := strconv.ParseFloat(record[1], 64)
		if yerr != nil {
			return nil, yerr
		}
		coords[i] = Coordinate{X, Y}
	}
	return coords, nil
}

//the cam is a 2d pizza. it is made of slices that all meet at the origin
func getPizzaSlice(a, b Coordinate) stl.Triangle {
	return stl.Triangle {
		Normal: stl.Vec3{ 0,0,1 },
		Vertices: [3]stl.Vec3{
			stl.Vec3{ 0,0,0 },
			stl.Vec3{ float32(a.X), float32(a.Y), 0 },
			stl.Vec3{ float32(b.X), float32(b.Y), 0 },
		},
	}
}

//directional arrow is used to show which direction the cam spins
func getDirectionalArrow(a, b Coordinate) stl.Triangle {
	return stl.Triangle {
		Normal: stl.Vec3{ 0,0,1 },
		Vertices: [3]stl.Vec3{
			stl.Vec3{ float32(a.X * 1.2), float32(a.Y * 1.2), 0 },
			stl.Vec3{ float32(a.X * 1.4), float32(a.Y * 1.4), 0 },
			stl.Vec3{ float32(b.X * 1.3), float32(b.Y * 1.3), 0 },
		},
	}
}

