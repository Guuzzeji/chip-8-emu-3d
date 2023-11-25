package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Keyboard struct {
	KeyMaps map[int32]int
}

type KeyAction interface {
	keyPress(uint8) int
	GetKeyPressed() int
}

// Keep map convert modern keys to chip 8  keys
var keymap = map[int32]int{
	49:  0x1, // 1
	50:  0x2, // 2
	51:  0x3, // 3
	52:  0xc, // 4
	113: 0x4, // Q
	119: 0x5, // W
	101: 0x6, // E
	114: 0xD, // R
	97:  0x7, // A
	115: 0x8, // S
	100: 0x9, // D
	102: 0xE, // F
	122: 0xA, // Z
	120: 0x0, // X
	99:  0xB, // C
	118: 0xF, // V
}

func (k *Keyboard) Init() {
	k.KeyMaps = keymap
}

// Part of keyboard struct, pass in key value we want to check for
// and return 1 or 0 if we did get that correct key value
func (k *Keyboard) KeyPress(keyValue uint8) int {
	key := rl.GetCharPressed()

	fmt.Println("KeyPress", keyValue, key)
	chip8Key := k.KeyMaps[key]

	if uint8(chip8Key) == keyValue {
		return 1
	}

	// if rl.IsKeyDown(rl.KeyKp1) {
	// 	fmt.Printf("pressed")
	// 	return 1
	// }

	// return 0

	return 0
}

// Checks to see if any key is pressed and returns that key value that was pressed
// Key value is base on keymap array
// Returns int and bool, keyvalue and wether key was pressed
func (k *Keyboard) GetKeyPressed() (int, bool) {
	key := rl.GetCharPressed()

	chip8Key := k.KeyMaps[key]
	isKeyPressed := false

	if key != 0 {
		isKeyPressed = true
		fmt.Println(key)
		// fmt.Println(k.KeyMaps[key])
	}

	return chip8Key, isKeyPressed //return zero if not in map
}
