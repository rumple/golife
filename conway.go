/*
The MIT License (MIT)

Copyright (c) 2014 amrut.joshi@gmail.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

type grid struct {
	g   [][]bool
	gen int
}

func NewGridFromFile(filename string) *grid {
	var x, y int
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if len(l) == 0 {
			continue
		}
		stars := strings.Count(l, "*")
		if stars > 2 {
			y = stars - 2
			continue
		}
		x++
	}
	g := NewGrid(x, y)
	for xi, line := range lines {
		l := strings.TrimSpace(line)
		if len(l) == 0 {
			continue
		}
		stars := strings.Count(l, "*")
		if stars > 2 {
			continue
		}
		for yi, c := range l[1 : len(l)-1] {
			if c == '+' {
				g.setVal(xi-1, yi, true)
			}
		}
	}
	return g
}

func NewGrid(x, y int) *grid {
	var g grid
	g.g = make([][]bool, x+2)
	for i, _ := range g.g {
		g.g[i] = make([]bool, y+2)
	}
	return &g
}

func (g *grid) getY() int {
	return len((*g).g[0]) - 2
}

func (g *grid) getX() int {
	return len(g.g) - 2
}

// x and y are zero based
func (g *grid) getVal(x, y int) bool {
	if x < 0 || x > g.getX() {
		return false
	}
	return (g.g)[x+1][y+1]
}

// x and y are zero based
func (g *grid) setVal(x, y int, val bool) {
	if x < 0 || x > g.getX() {
		return
	}
	(g.g)[x+1][y+1] = val
}

func (g *grid) showGrid() {
	fmt.Print("\033[2J")
	fmt.Println("Gen:", g.gen)
	gx := g.getX()
	gy := g.getY()
	for x := 0; x < gx+2; x++ {
		for y := 0; y < gy+2; y++ {
			if x == 0 || y == 0 || x == gx+1 || y == gy+1 {
				fmt.Print("*")
				continue
			}
			if g.getVal(x-1, y-1) {
				fmt.Print("+")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (g *grid) addOscillator(x, y int) {
	g.setVal(x, y, true)
	g.setVal(x, y+1, true)
	g.setVal(x, y+2, true)
}

func (g *grid) addBeacon(x, y int) {
	g.setVal(x, y, true)
	g.setVal(x, y+1, true)
	g.setVal(x+1, y, true)
	g.setVal(x+2, y+3, true)
	g.setVal(x+3, y+3, true)
	g.setVal(x+3, y+2, true)
}

// x and y are zero based
func (g *grid) liveNeighbors(x, y int) int {
	n := 0
	for internalx := x; internalx <= x+2; internalx++ {
		for internaly := y; internaly <= y+2; internaly++ {
			if internalx == x+1 && internaly == y+1 {
				continue
			}
			if (g.g)[internalx][internaly] {
				n++
			}
		}
	}
	return n
}

func (g *grid) nextGen() int {
	active := 0
	newg := NewGrid(g.getX(), g.getY())
	for x := 0; x < g.getX(); x++ {
		for y := 0; y < g.getY(); y++ {
			n := g.liveNeighbors(x, y)
			val := g.getVal(x, y)
			newg.setVal(x, y, val)
			// Any live cell with fewer than two live neighbours dies, as if caused by under-population.
			if n < 2 && val == true {
				newg.setVal(x, y, false)
				active++
			}
			// Any live cell with more than three live neighbours dies, as if by overcrowding.
			if n > 3 && val == true {
				newg.setVal(x, y, false)
				active++
			}
			// Any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction.
			if n == 3 && val == false {
				newg.setVal(x, y, true)
				active++
			}
		}
	}
	for i, _ := range g.g {
		(g.g)[i] = (newg.g)[i]
	}
	g.gen++
	return active
}

func (g *grid) runLife(gen int) {
	for i := 0; i < gen; i++ {
		g.showGrid()
		active := g.nextGen()
		if active == 0 {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	g := NewGridFromFile("./conway.txt")
	g.showGrid()
	time.Sleep(2 * time.Second)
	g.runLife(1000000)
}
