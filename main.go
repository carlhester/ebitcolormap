package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const TILES_WIDE = 5
const TILES_HIGH = 5
const TILE_SIZE = 100

type Game struct {
	gameMap *gameMap
	next    bool
}

type gameMap struct {
	ff        font.Face
	rand      *rand.Rand
	tiles     []*tile
	textColor color.Color
}

func newGameMap() *gameMap {
	tt, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	ff, _ := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return &gameMap{
		ff:        ff,
		rand:      r1,
		textColor: color.Black,
	}
}

func (m *gameMap) update() {
	m.tiles = []*tile{}
	for y := 0; y < TILES_HIGH; y++ {
		for x := 0; x < TILES_WIDE; x++ {
			r := uint8(m.rand.Intn(255))
			g := uint8(m.rand.Intn(255))
			b := uint8(m.rand.Intn(255))
			tile := newTile(x, y, color.RGBA{r, g, b, 0xff})
			text.Draw(tile.img, fmt.Sprintf("%d\n%d\n%d", r, g, b), m.ff, x+(TILE_SIZE/3), y+(TILE_SIZE/6), m.textColor)
			m.tiles = append(m.tiles, tile)
		}
	}
}

func (m *gameMap) Draw(screen *ebiten.Image) {
	for _, tile := range m.tiles {
		tile.Draw(screen)
	}
}

type tile struct {
	x     int
	y     int
	color color.Color
	img   *ebiten.Image
}

func newTile(x int, y int, color color.RGBA) *tile {
	img := ebiten.NewImage(TILE_SIZE, TILE_SIZE)
	img.Fill(color)
	return &tile{
		x:     x,
		y:     y,
		color: color,
		img:   img,
	}
}

func (b *tile) Draw(screen *ebiten.Image) {
	loc := &ebiten.DrawImageOptions{}
	loc.GeoM.Translate(float64(b.x*TILE_SIZE), float64(b.y*TILE_SIZE))
	screen.DrawImage(b.img, loc)
}

func NewGame() *Game {
	ebiten.SetWindowTitle("Ebiten Test Game")
	ebiten.SetWindowSize(TILE_SIZE*TILES_WIDE, TILE_SIZE*TILES_HIGH)
	ebiten.SetWindowResizable(true)

	gameMap := newGameMap()
	return &Game{
		gameMap: gameMap,
		next:    true,
	}
}

func (game *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		os.Exit(0)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		game.next = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		if game.gameMap.textColor == color.White {
			game.gameMap.textColor = color.Black
		} else {
			game.gameMap.textColor = color.White
		}
	}
	return nil
}

func (game *Game) Draw(screen *ebiten.Image) {
	if game.next {
		game.gameMap.update()
		game.next = false
	}
	screen.Fill(color.RGBA{0xff, 0xff, 0xff, 0xff})
	game.gameMap.Draw(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("(N) to change colors.\n(C) to change text color.\n(Q) to Quit.\nTPS: %0.2f", ebiten.CurrentTPS()))
}

func (game *Game) Layout(w, h int) (int, int) {
	return TILE_SIZE * TILES_WIDE, TILE_SIZE * TILES_HIGH
}

func main() {
	game := NewGame()
	ebiten.RunGame(game)

}
