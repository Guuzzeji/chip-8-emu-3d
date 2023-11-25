package main

import (
	"fmt"
)

type Screen struct {
	Height int
	Width  int
	Scale  int
	Canvas [33][65]int
}

func (s *Screen) Init() {
	s.Height = 32
	s.Width = 64
	s.Scale = 2
}

func (s *Screen) Clear() {
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			s.Canvas[y][x] = 0
		}
	}
}

func (s *Screen) DrawPixel(x int, y int) bool {

	// Checking bounds for width and height and offset if needed
	if x > s.Width {
		x -= s.Width
	} else if x < 0 {
		x += s.Width
	}

	if y > s.Height {
		y -= s.Height
	} else if y < 0 {
		y += s.Height
	}

	s.Canvas[y][x] ^= 1 // switch pixel on and off

	return s.Canvas[y][x] == 1 // return true if 1, else false
}

// DEBUG: Used for debugging into console
func (s *Screen) Render() {
	fmt.Print("\n")
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			fmt.Print(s.Canvas[y][x])
		}
		fmt.Print("\n")
	}
}
