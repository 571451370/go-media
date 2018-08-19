package obj

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/qeedquan/go-media/image/imageutil"
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
	Name           string
	Colors         [3]f64.Vec3
	SpecularFactor float64
	DissolveFactor float64
	Illumination   int
	Texture        struct {
		Displacement    *Texture
		Diffuse         *Texture
		Ambient         *Texture
		SpecularColor   *Texture
		SpecularHilight *Texture
		Alpha           *Texture
		Bump            *Texture
		Dissolve        *Texture
	}
}

type Texture struct {
	Blend       [2]bool
	MipmapBoost float64
	Origin      f64.Vec3
	Scale       f64.Vec3
	Turbulence  f64.Vec3
	Clamp       bool
	BumpFactor  bool
	IMF         uint64
	Map         *image.RGBA
}

func Load(name string, r io.Reader) (*Model, error) {
	var err error
	m := &Model{}
	s := bufio.NewScanner(r)
	dir := filepath.Dir(name)
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
			m.Mats, err = addMat(dir, m.Mats, line)
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

func addMat(dir string, mat []Material, line string) ([]Material, error) {
	var (
		t string
		e string
	)
	n, _ := fmt.Sscan(line, &t, &e)
	if n != 2 {
		return mat, nil
	}

	f, err := os.Open(filepath.Join(dir, e))
	if err != nil {
		return mat, err
	}
	defer f.Close()

	p := make([]Material, 1)
	m := &p[0]

	s := bufio.NewScanner(f)
	for s.Scan() {
		line = s.Text()
		switch {
		case strings.HasPrefix(line, "newmtl "):
			var name string
			fmt.Sscan(line, &t, &name)
			if len(p) == 1 {
				m.Name = name
			} else {
				p = append(p, Material{Name: name})
				m = &p[len(p)-1]
			}
		case strings.HasPrefix(line, "illum "):
			fmt.Sscan(line, &t, &m.Illumination)
		case strings.HasPrefix(line, "Ka "):
			fmt.Sscan(line, &t, &m.Colors[0].X, &m.Colors[0].Y, &m.Colors[0].Z)
		case strings.HasPrefix(line, "Kd "):
			fmt.Sscan(line, &t, &m.Colors[1].X, &m.Colors[1].Y, &m.Colors[1].Z)
		case strings.HasPrefix(line, "Ks "):
			fmt.Sscan(line, &t, &m.Colors[2].X, &m.Colors[2].Y, &m.Colors[2].Z)
		case strings.HasPrefix(line, "Ns "):
			fmt.Sscan(line, &t, &m.SpecularFactor)
		case strings.HasPrefix(line, "d "):
			fmt.Sscan(line, &t, &m.DissolveFactor)
		case strings.HasPrefix(line, "map_disp ") || strings.HasPrefix(line, "disp "):
			m.Texture.Displacement, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_bump ") || strings.HasPrefix(line, "bump "):
			m.Texture.Bump, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_Ka "):
			m.Texture.Ambient, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_Kd "):
			m.Texture.Diffuse, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_Ks "):
			m.Texture.SpecularColor, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_Ns "):
			m.Texture.SpecularHilight, err = loadTexture(dir, line)
		case strings.HasPrefix(line, "map_d "):
			m.Texture.Alpha, err = loadTexture(dir, line)
		}

		if err != nil {
			return nil, err
		}
	}

	if len(p) > 1 || m.Name != "" {
		mat = append(mat, p...)
	}

	return mat, nil
}

func loadTexture(dir, line string) (*Texture, error) {
	t := &Texture{
		Blend: [2]bool{true, true},
		Clamp: true,
	}

	var (
		s   scanner.Scanner
		err error
	)
	s.Init(strings.NewReader(line))
	for {
		tok := s.Scan()
		if tok == '-' {
			tok = s.Scan()
		}
		if tok != scanner.Ident {
			continue
		}

		v := s.TokenText()
		switch v {
		case "bm":
		case "clamp":
			t.Clamp = scanBool(&s)
		case "blendu":
			t.Blend[0] = scanBool(&s)
		case "blendv":
			t.Blend[1] = scanBool(&s)
		case "imfchan":
		case "mm":
		case "o":
			t.Origin.X = scanFloat(&s)
			t.Origin.Y = scanFloat(&s)
			t.Origin.Z = scanFloat(&s)
		case "s":
			t.Scale.X = scanFloat(&s)
			t.Scale.Y = scanFloat(&s)
			t.Scale.Z = scanFloat(&s)
		case "t":
			t.Turbulence.X = scanFloat(&s)
			t.Turbulence.Y = scanFloat(&s)
			t.Turbulence.Z = scanFloat(&s)
		case "texres":
		default:
			t.Map, err = imageutil.LoadRGBAFile(filepath.Join(dir, v))
			if err != nil {
				return t, fmt.Errorf("%s: failed to load texture file: %v", v, err)
			}
		}
	}

	return t, nil
}

func scanFloat(s *scanner.Scanner) float64 {
	s.Scan()
	n, _ := strconv.ParseFloat(s.TokenText(), 64)
	return n
}

func scanBool(s *scanner.Scanner) bool {
	s.Scan()
	switch strings.ToLower(s.TokenText()) {
	case "on", "true", "1":
		return true
	}
	return false
}