package cams

import (
	"github.com/hschendel/stl"
	"fmt"
	"encoding/csv"
	"os"
	"strconv"
)

type Coordinate struct {
	x, y float32
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
		x, xerr := strconv.ParseFloat(record[0], 32)
		if xerr != nil {
			return nil, xerr
		}
		y, yerr := strconv.ParseFloat(record[1], 32)
		if yerr != nil {
			return nil, yerr
		}
		coords[i] = Coordinate{float32(x), float32(y)}
	}
	return coords, nil
}

//the cam is a 2d pizza. it is made of slices that all meet at the origin
func getPizzaSlice(a, b Coordinate) stl.Triangle {
	return stl.Triangle {
		Normal: stl.Vec3{ 0,0,1 },
		Vertices: [3]stl.Vec3{
			stl.Vec3{ 0,0,0 },
			stl.Vec3{ a.x,a.y,0 },
			stl.Vec3{ b.x,b.y,0 },
		},
	}
}

//directional arrow is used to show which direction the cam spins
func getDirectionalArrow(a, b Coordinate) stl.Triangle {
	return stl.Triangle {
		Normal: stl.Vec3{ 0,0,1 },
		Vertices: [3]stl.Vec3{
			stl.Vec3{ a.x * 1.2, a.y * 1.2 ,0 },
			stl.Vec3{ a.x * 1.4, a.y * 1.4 ,0 },
			stl.Vec3{ b.x * 1.3, b.y * 1.3 ,0 },
		},
	}
}

