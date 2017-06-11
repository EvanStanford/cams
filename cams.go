package cams

import (
	"github.com/hschendel/stl"
	"fmt"
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

