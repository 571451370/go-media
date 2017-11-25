package tiled

import (
	"encoding/xml"

	"github.com/qeedquan/go-media/ioe"
)

type Map struct {
}

type TMX struct {
	XMLName      xml.Name `xml:"map"`
	Version      string   `xml:"version,attr"`
	TiledVersion string   `xml:"tiledversion,attr"`
	Orientation  string   `xml:"orientation,attr"`
	RenderOrder  string   `xml:"renderorder,attr"`
	Width        int      `xml:"width,attr"`
	Height       int      `xml:"height,attr"`
	TileWidth    int      `xml:"tilewidth,attr"`
	TileHeight   int      `xml:"tileheight,attr"`
	NextObjectID int      `xml:"nextobjectid,attr"`
	Tileset      []*TSX   `xml:"tileset"`
	Layer        []struct {
		Name    string `xml:"name,attr"`
		Width   int    `xml:"width,attr"`
		Height  int    `xml:"height,attr"`
		Visible *int   `xml:"visible,attr"`
		Data    struct {
			Encoding    string `xml:"encoding,attr"`
			Compression string `xml:"compression,attr"`
			Tile        []struct {
				Gid int `xml:"gid,attr"`
			} `xml:"tile"`
			Chardata string `xml:",chardata"`
		} `xml:"data"`
	} `xml:"layer"`
}

type TSX struct {
	XMLName    xml.Name `xml:"tileset"`
	FirstGID   int      `xml:"firstgid,attr"`
	Source     string   `xml:"source,attr"`
	Name       string   `xml:"name,attr"`
	TileWidth  int      `xml:"tilewidth,attr"`
	TileHeight int      `xml:"tileheight,attr"`
	TileCount  int      `xml:"tilecount,attr"`
	Columns    int      `xml:"columns,attr"`
	Margin     int      `xml:"margin,attr"`
	Spacing    int      `xml:"spacing,attr"`
	Image      struct {
		Source string `xml:"source,attr"`
		Trans  string `xml:"trans,attr"`
		Width  int    `xml:"width,attr"`
		Height int    `xml:"height,attr"`
	} `xml:"image"`
}

func OpenMap(fs ioe.FS, name string) (*Map, error) {
	buf, err := ioe.ReadFile(fs, name)
	if err != nil {
		return nil, err
	}

	var tmx TMX
	err = xml.Unmarshal(buf, &tmx)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
