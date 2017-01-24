package obj

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/qeedquan/go-media/math/f64"
)

type Model struct {
	Verts   []f64.Vec3
	Coords  []f64.Vec3
	Normals []f64.Vec3
	Faces   [][3][3]int
	Mats    []Material
}

type Material struct {
	Name     string
	Colors   [3]f64.Vec3
	Diffuse  *image.RGBA
	Ambient  *image.RGBA
	Specular *image.RGBA
	Alpha    *image.RGBA
	Bump     *image.RGBA
}

func Load(obj string) (*Model, error) {
	f, err := os.Open(obj)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := &Model{}

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		switch {
		case strings.HasPrefix(line, "v "):
			m.Verts = addVert(m.Verts, line)
		case strings.HasPrefix(line, "vt "):
			m.Coords = addVert(m.Coords, line)
		case strings.HasPrefix(line, "vn "):
			m.Normals = addVert(m.Normals, line)
		case strings.HasPrefix(line, "f "):
			m.Faces = addFace(m.Faces, line)
		}
	}

	return m, nil
}

func addVert(verts []f64.Vec3, line string) []f64.Vec3 {
	var (
		t string
		p [3]float64
	)
	n, _ := fmt.Sscan(line, &t, &p[0], &p[1], &p[2])
	if n != 3 {
		n, _ = fmt.Sscan(line, &t, &p[0], &p[1])
	}
	if n != 3 && n != 2 {
		return verts
	}

	return append(verts, f64.Vec3{p[0], p[1], p[2]})
}

func addFace(faces [][3][3]int, line string) [][3][3]int {
	var (
		t string
		f [9]int
	)
	n, _ := fmt.Sscanf(line, "%s %d/%d/%d %d/%d/%d %d/%d/%d",
		&t, &f[0], &f[1], &f[2], &f[3], &f[4], &f[5], &f[6], &f[7], &f[8])
	if n != 10 {
		n, _ = fmt.Sscanf(line, "%s %d/%d %d/%d %d/%d",
			&t, &f[0], &f[1], &f[3], &f[4], &f[6], &f[7])
	}
	if n != 10 && n != 7 {
		n, _ = fmt.Sscanf(line, "%s %d %d %d",
			&t, &f[0], &f[3], &f[6])
	}
	if n != 10 && n != 7 && n != 4 {
		return faces
	}

	return append(faces, [3][3]int{
		{f[0], f[1], f[2]},
		{f[3], f[4], f[5]},
		{f[6], f[7], f[8]},
	})
}
