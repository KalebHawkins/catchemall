// Copyright 2024 Kaleb Hawkins <KalebHawkins@outlook.com>

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth   = 800
	screenHeight  = 600
	spawnInterval = 60
	objWidth      = 40
	objHeight     = 40
	objSpeed      = 5
)

type GameState int

const (
	StatePlaying GameState = iota
	StateGameOver
)

type Player struct {
	X, Y          float64
	Width, Height float64
}

type Object struct {
	X, Y          float64
	Width, Height float64
	Speed         float64
	DangerCube    bool
}

type Game struct {
	Player       *Player
	Objects      []*Object
	spawnCounter int

	playerScore int
	playerLives int

	state GameState
}

func (g *Game) Update() error {
	if g.state == StateGameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			g.Reset()
		}
		return nil
	}

	mx, _ := ebiten.CursorPosition()
	g.Player.X = float64(mx) - g.Player.Width/2

	if g.Player.X < 0 {
		g.Player.X = 0
	} else if g.Player.X > screenWidth-g.Player.Width {
		g.Player.X = screenWidth - g.Player.Width
	}

	g.spawnCounter++
	if g.spawnCounter >= spawnInterval {
		randomX := float64(rand.Intn(screenWidth - objWidth))
		isDangerCube := rand.Float64() < 0.3 // 30% chance of a cube being a danger cube

		g.Objects = append(g.Objects, &Object{
			X:          randomX,
			Y:          0,
			Width:      objWidth,
			Height:     objHeight,
			Speed:      objSpeed,
			DangerCube: isDangerCube,
		})
		g.spawnCounter = 0
	}

	// Iterate backwards to safely remove objects that have fallen offscreen.
	for i := len(g.Objects) - 1; i >= 0; i-- {
		g.Objects[i].Y += g.Objects[i].Speed

		if g.Objects[i].Y > screenHeight {

			if !g.Objects[i].DangerCube {
				g.playerLives--
			}

			g.Objects = append(g.Objects[:i], g.Objects[i+1:]...)
		}

		if g.Player.X < g.Objects[i].X+g.Objects[i].Width &&
			g.Player.X+g.Player.Width > g.Objects[i].X &&
			g.Player.Y < g.Objects[i].Y+g.Objects[i].Height &&
			g.Player.Y+g.Player.Height > g.Objects[i].Y {

			if g.Objects[i].DangerCube {
				g.playerLives--
			} else {
				g.playerScore++
			}

			g.Objects = append(g.Objects[:i], g.Objects[i+1:]...)

		}
	}

	if g.playerLives <= 0 {
		g.state = StateGameOver
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x33, 0x33, 0x33, 0xff})

	if g.state == StateGameOver {
		ebitenutil.DebugPrint(screen, "Game Over!\nPress Space or Click to Restart")
		return
	}

	vector.DrawFilledRect(screen, float32(g.Player.X), float32(g.Player.Y), float32(g.Player.Width), float32(g.Player.Height), color.Black, true)

	for _, obj := range g.Objects {
		if obj.DangerCube {
			vector.DrawFilledRect(screen, float32(obj.X), float32(obj.Y), float32(obj.Width), float32(obj.Height), color.RGBA{0xff, 0x0, 0x00, 0xff}, true)
		} else {
			vector.DrawFilledRect(screen, float32(obj.X), float32(obj.Y), float32(obj.Width), float32(obj.Height), color.RGBA{0x00, 0xff, 0x00, 0xff}, true)
		}
	}

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Score: %d\nLives: %d\nObject Count: %d", g.playerScore, g.playerLives, len(g.Objects)),
	)
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func (g *Game) Reset() {
	g.Objects = make([]*Object, 0)
	g.playerScore = 0
	g.playerLives = 3
	g.spawnCounter = 0
	g.state = StatePlaying
}

func main() {
	ebiten.SetWindowTitle("Catch 'Em All")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	g := &Game{
		Player:       &Player{X: (screenWidth - 120) / 2, Y: (screenHeight - 20 - 20), Width: 120, Height: 20},
		Objects:      make([]*Object, 0),
		spawnCounter: 0,
		playerScore:  0,
		playerLives:  3,
		state:        StatePlaying,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
