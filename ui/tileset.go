package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tileset struct {
	image     *ebiten.Image
	tileSizeX int
	tileSizeY int
	nTilesW   int
	nTilesH   int

	TileWidth  float64
	TileHeight float64
}

func (t *Tileset) TileAtIndex(index int) *ebiten.Image {
	x := index % t.nTilesW
	y := index / (t.nTilesW)
	startX := x * t.tileSizeX
	startY := y * t.tileSizeY
	return t.image.SubImage(image.Rect(startX, startY, startX+t.tileSizeX, startY+t.tileSizeY)).(*ebiten.Image)
}

func (t *Tileset) TileAtIJ(i, j int) *ebiten.Image {
	startX := i * t.tileSizeX
	startY := j * t.tileSizeY
	return t.image.SubImage(image.Rect(startX, startY, startX+t.tileSizeX, startY+t.tileSizeY)).(*ebiten.Image)
}

func NewTileset(image *ebiten.Image, tileSizeX, tileSizeY int) *Tileset {
	w, h := image.Bounds().Dx(), image.Bounds().Dy()

	return &Tileset{
		image:      image,
		tileSizeX:  tileSizeX,
		tileSizeY:  tileSizeY,
		nTilesW:    w / tileSizeX,
		nTilesH:    h / tileSizeY,
		TileWidth:  float64(tileSizeX),
		TileHeight: float64(tileSizeY),
	}
}
