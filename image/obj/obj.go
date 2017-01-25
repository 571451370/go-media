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
	Verts   []f64.Vec4
	Coords  []f64.Vec4
	Normals []f64.Vec4
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
		case strings.HasPrefix(line, "mtllib "):
			m.Mats, err = addMat(m.Mats, line)
			if err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func addVert(verts []f64.Vec4, line string) []f64.Vec4 {
	var (
		t string
		p [4]float64
	)
	p[3] = 1

	n, _ := fmt.Sscan(line, &t, &p[0], &p[1], &p[2], &p[3])
	if n != 4 {
		n, _ = fmt.Sscan(line, &t, &p[0], &p[1], &p[2])
	}
	if n != 4 && n != 3 {
		n, _ = fmt.Sscan(line, &t, &p[0], &p[1])
	}
	if n != 4 && n != 3 && n != 2 {
		return verts
	}

	return append(verts, f64.Vec4{p[0], p[1], p[2], p[3]})
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

func addMat(mat []Material, line string) ([]Material, error) {
	var (
		t string
		e string
	)
	n, _ := fmt.Sscan(line, &t, &e)
	if n != 2 {
		return mat, nil
	}

	f, err := os.Open(e)
	if err != nil {
		return mat, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for {
		line = s.Text()
		switch {
		case strings.HasPrefix(line, "newmtl "):
		case strings.HasPrefix(line, "Ka "):
		case strings.HasPrefix(line, "Kd "):
		case strings.HasPrefix(line, "Ks "):
		case strings.HasPrefix(line, "map_Ka "):
		case strings.HasPrefix(line, "map_Kd "):
		case strings.HasPrefix(line, "map_Ks "):
		case strings.HasPrefix(line, "map_Ns "):
		case strings.HasPrefix(line, "map_d "):
		}
	}
}
