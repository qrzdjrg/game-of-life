package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

const defaultSize int = 10

var neighs = [...]Coord2D{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

type Game struct {
	Grid
	M, N int
}

type Grid [][]byte

type OpCell struct {
	cell Coord2D
	val  int
}

type Coord2D struct {
	x, y int
}

func (c Coord2D) isInBound(lowerBound int, upperBound int) bool {
	return c.x < upperBound && c.y < upperBound &&
		c.x >= lowerBound && c.y >= lowerBound
}

func (c Coord2D) Add(coord Coord2D) Coord2D {
	return Coord2D{c.x + coord.x, c.y + coord.y}
}

func (g *Game) New() {
	g.NewSquare(defaultSize)
}

func (g *Game) NewSquare(m int) {
	g.Grid.New(m, m)
	g.M = m
	g.N = m
}

func (g *Grid) New(m, n int) {
	*g = make([][]byte, m)
	for i := range defaultSize {
		(*g)[i] = make([]byte, n)
	}
}

func (g *Grid) Clear() {
	for i := range defaultSize {
		clear((*g)[i])
	}
	clear(*g)
}

func (g *Game) RunN(n int, millis int64) {
	for i := range n {
		if millis == 0 {
			g.Iterate()
		} else {
			fmt.Printf("It %v\n", i+1)
			g.Iterate()
			time.Sleep(time.Duration(millis) * time.Millisecond)
			g.Print()
		}
	}
}

func (g *Game) initOscillator() {
	g.Set(1, 1)
	g.Set(1, 2)
	g.Set(1, 3)
}

func (g *Game) initOscillator2() {
	g.Set(4, 4)
	g.Set(4, 5)
	g.Set(5, 4)
	g.Set(5, 5)

	g.Set(6, 6)
	g.Set(6, 7)
	g.Set(7, 6)
	g.Set(7, 7)

}

func (g *Grid) Iterate() {

	mods := make([]OpCell, 0)
	for i := range defaultSize {
		for j := range defaultSize {
			// fmt.Print((*g)[i][j])
			live := 0
			currentCell := Coord2D{i, j}
			for _, neighCell := range g.neighCoords(currentCell) {
				if g.isAlive(neighCell) {
					live++
				}
			}
			g.processCell(currentCell, live)
			mods = append(mods, g.processCell(currentCell, live))
		}
	}

	for _, mod := range mods {
		if mod.val == -1 {
			g.UnsetCoord(mod.cell)
		} else if mod.val == 1 {
			g.SetCoord(mod.cell)
		}
	}

}

func (g *Grid) processCell(coord Coord2D, liveNeigh int) OpCell {

	if g.isAlive(coord) {
		if liveNeigh != 2 && liveNeigh != 3 {
			return OpCell{coord, -1}
		}
	} else {
		if liveNeigh == 3 {
			return OpCell{coord, 1}
		}
	}
	return OpCell{coord, 0}
}

func (g *Grid) isAlive(coord Coord2D) bool {
	return (*g)[coord.x][coord.y] == 1
}

func (g *Grid) neighCoords(coord Coord2D) []Coord2D {

	res := make([]Coord2D, 0)

	for _, v := range neighs {
		c := coord.Add(v)
		if c.isInBound(0, defaultSize) {
			res = append(res, c)
		}
	}
	return res
}

func (g *Grid) SetCoord(coord Coord2D) {
	g.Set(coord.x, coord.y)
}

func (g *Grid) UnsetCoord(coord Coord2D) {
	g.Unset(coord.x, coord.y)
}

func (g *Grid) Set(i int, j int) {
	(*g)[i][j] = 1
}

func (g *Grid) Unset(i int, j int) {
	(*g)[i][j] = 0
}

func (g *Grid) String() string {
	var result strings.Builder
	for _, n := range *g {
		result.WriteString(fmt.Sprintf("%+v\n", n))
	}
	return result.String()
}

func (g *Grid) Print() {
	fmt.Println(g.String())
}

func main() {
	// var rtm runtime.MemStats
	// dumpMemStats("First Mem", rtm)

	fmt.Println("Game of Life")

	game := &Game{}
	game.NewSquare(10)

	// dumpMemStats("Alloc Mem", rtm)

	game.initOscillator()
	game.initOscillator2()
	game.Print()
	game.RunN(5, 500)
	game.Clear()

	// dumpMemStats("Last Mem", rtm)
}

func dumpMemStats(message string, rtm runtime.MemStats) {
	runtime.ReadMemStats(&rtm)
	fmt.Println(" \n=== ", message, " === ")
	fmt.Println(" Mallocs : ", rtm.Mallocs)
	fmt.Println(" Frees : ", rtm.Frees)
	fmt.Println(" Live Objects : ", rtm.Mallocs, rtm.Frees)
	fmt.Println(" PauseTotalNs : ", rtm.PauseTotalNs)
	fmt.Println(" NumGC : ", rtm.NumGC)
	fmt.Println(" LastGC : ", time.UnixMilli(int64(rtm.LastGC/1_000_000)))
	fmt.Println(" HeapObjects : ", rtm.HeapObjects)
	fmt.Println(" HeapAlloc : ", rtm.HeapAlloc)
}
