package main

import (
	"fmt"
	"math/rand"
	"os"
)

type Cpu struct {
	Memory         [4096]uint16
	Registers      [16]uint8
	Stack          [16]uint16
	StackPointer   uint16
	ProgramCounter uint16
	IndexRegisters uint16
	SoundTimer     uint16
	DelayTimer     uint16
	Speed          int
	Display        *Screen
	Keyboard       *Keyboard
	Pause          bool
}

var sprites = [80]uint16{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func (c *Cpu) Init() {

	for i := 1; i < len(sprites); i++ {
		c.Memory[i] = sprites[i]
	}

	c.StackPointer = 0
	c.ProgramCounter = 0
	c.IndexRegisters = 0
	c.SoundTimer = 0
	c.DelayTimer = 0

	c.ProgramCounter = 0x200
	c.Speed = 10
	c.Pause = false

	var keyboard Keyboard
	keyboard.Init()
	c.Keyboard = &keyboard

	var Display Screen
	Display.Init()

	c.Display = &Display
}

// Reads infomation and adds it to memory
func (c *Cpu) loadProgramIntoMemory(program []byte) {
	for i := 0; i < len(program); i++ {
		c.Memory[0x200+i] = uint16(program[i])
	}
}

// Used to load rom into memmory base on filepath
func (c *Cpu) LoadRom(filepath string) {
	data, err := os.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	c.loadProgramIntoMemory(data)
}

// Run CPU Cycle base on the speed set in CPU struct
func (c *Cpu) Cycle() {
	for i := 0; i < c.Speed; i++ {
		opcode := (c.Memory[c.ProgramCounter]<<8 | c.Memory[c.ProgramCounter+1])
		c.ExecuteInstruction(uint16(opcode))
	}
}

// Run Opcode and decodes it to do correct action
// ! BUG: For some programs, it will just explode and not run correctly,
// ! error in rendering, but program is still running correctly
func (c *Cpu) ExecuteInstruction(opcode uint16) {
	vX := (opcode & 0x0F00) >> 8
	vY := (opcode & 0x00F0) >> 4

	if c.Pause {
		// Used for 0x0A
		keyValue, isKeyPressed := c.Keyboard.GetKeyPressed()

		if isKeyPressed {
			c.ProgramCounter += 2
			c.Registers[vX] = uint8(keyValue)
			c.Pause = false
		}
	}

	// fmt.Printf("%X, %X, %X, %X, %X, %X \n", opcode, vX, vY, c.IndexRegisters, c.ProgramCounter, c.StackPointer)
	// fmt.Printf("OPCODE = %X", opcode)

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode & 0x000F {

		case 0x0000:
			c.Display.Clear()
			c.ProgramCounter += 2

		case 0x000E:
			c.StackPointer -= 1
			c.ProgramCounter = c.Stack[c.StackPointer]
			c.ProgramCounter += 2

		default:
			fmt.Printf("Something broke 0x0000")

		}

	case 0x1000:
		//fmt.Printf("Hit 0x1000")
		c.ProgramCounter = opcode & 0x0FFF

	case 0x2000:
		c.Stack[c.StackPointer] = c.ProgramCounter
		c.StackPointer += 1
		c.ProgramCounter = opcode & 0x0FFF

	case 0x3000:
		//fmt.Printf("Hit 0x3000")
		if c.Registers[vX] == uint8(opcode&0x00FF) {
			c.ProgramCounter += 2
		}
		c.ProgramCounter += 2

	case 0x4000:
		if c.Registers[vX] != uint8(opcode&0x00FF) {
			c.ProgramCounter += 2
		}
		c.ProgramCounter += 2

	case 0x5000:
		if c.Registers[vX] == c.Registers[vY] {
			c.ProgramCounter += 2
		}
		c.ProgramCounter += 2

	case 0x6000:
		c.Registers[vX] = uint8(opcode & 0x00FF)
		c.ProgramCounter += 2

	case 0x7000:
		c.Registers[vX] += uint8(opcode & 0x00FF)
		c.ProgramCounter += 2

	case 0x8000:
		switch opcode & 0x000F {
		case 0x0:
			c.Registers[vX] = c.Registers[vY]
			c.ProgramCounter += 2

		case 0x1:
			c.Registers[vX] |= c.Registers[vY]
			c.ProgramCounter += 2

		case 0x2:
			c.Registers[vX] &= c.Registers[vY]
			c.ProgramCounter += 2

		case 0x3:
			c.Registers[vX] ^= c.Registers[vY]
			c.ProgramCounter += 2

		case 0x4:
			sum := c.Registers[vX] + c.Registers[vY]

			c.Registers[0xF] = 0

			if sum > 0xFF {
				c.Registers[0xF] = 1
			}

			c.Registers[vX] = sum
			c.ProgramCounter += 2

		case 0x5:
			c.Registers[0xF] = 0

			if c.Registers[vX] > c.Registers[vY] {
				c.Registers[0xF] = 1
			}

			c.Registers[vX] -= c.Registers[vY]
			c.ProgramCounter += 2

		case 0x6:
			c.Registers[0xF] = c.Registers[vX] & 0x0001
			c.Registers[vX] >>= 1
			c.ProgramCounter += 2

		case 0x7:
			c.Registers[0xF] = 0

			if c.Registers[vY] > c.Registers[vX] {
				c.Registers[0xF] = 1
			}

			c.Registers[vX] = c.Registers[vY] - c.Registers[vX]
			c.ProgramCounter += 2

		case 0xE:
			c.Registers[0xF] = c.Registers[vX] >> 7
			c.Registers[vX] <<= 1
			c.ProgramCounter += 2

		default:
			//fmt.Printf("Something broke 0x8000")

		}

	case 0x9000:
		if c.Registers[vX] != c.Registers[vY] {
			c.ProgramCounter += 2
		}
		c.ProgramCounter += 2

	case 0xA000:
		c.IndexRegisters = opcode & 0x0FFF
		c.ProgramCounter += 2

	case 0xB000:
		c.ProgramCounter = (opcode & 0x0FFF) + uint16(c.Registers[0])

	case 0xC000:
		rand := rand.Int() * 0x00FF
		c.Registers[vX] = uint8(uint16(rand) & (opcode & 0x00FF))
		c.ProgramCounter += 2

	case 0xD000:
		width := 8
		height := (opcode & 0x000F)

		c.Registers[0xF] = 0

		for row := 0; row < int(height); row++ {
			sprite := c.Memory[int(c.IndexRegisters)+row]

			for col := 0; col < width; col++ {
				if (sprite & 0x80) > 0 {
					if c.Display.DrawPixel(int(c.Registers[vX])+col, int(c.Registers[vY])+row) {
						c.Registers[0xF] = 1
					}
				}

				sprite <<= 1
			}
		}
		c.ProgramCounter += 2

	case 0xE000:
		switch opcode & 0x00FF {
		case 0x9E:
			if c.Keyboard.KeyPress(c.Registers[vX]) == 1 {
				c.ProgramCounter += 2
			}
			c.ProgramCounter += 2

		case 0xA1:
			if c.Keyboard.KeyPress(c.Registers[vX]) == 0 {
				c.ProgramCounter += 2
			}
			c.ProgramCounter += 2

		default:
			fmt.Printf("Something broke 0xE000")

		}

	case 0xF000:
		//fmt.Printf("Hit 0xF000")
		switch opcode & 0x00FF {
		case 0x07:
			c.Registers[vX] = uint8(c.DelayTimer)
			c.ProgramCounter += 2

		case 0x0A:
			c.Pause = true

		case 0x15:
			c.DelayTimer = uint16(c.Registers[vX])
			c.ProgramCounter += 2

		case 0x18:
			c.SoundTimer = uint16(c.Registers[vX])
			c.ProgramCounter += 2

		case 0x1E:
			if c.IndexRegisters+vX > 0x0FFF {
				c.Registers[0xF] = 1
			} else {
				c.Registers[0xF] = 1
			}

			c.IndexRegisters += uint16(c.Registers[vX])
			c.ProgramCounter += 2

		case 0x29:
			c.IndexRegisters = uint16(c.Registers[vX]) * 5
			c.ProgramCounter += 2

		case 0x33:
			c.Memory[c.IndexRegisters] = uint16(c.Registers[vX] / 100)
			c.Memory[c.IndexRegisters+1] = uint16((c.Registers[vX] & 100) / 10)
			c.Memory[c.IndexRegisters+2] = uint16((c.Registers[vX] % 10))
			c.ProgramCounter += 2

		case 0x55:
			for i := 0; i <= int(vX); i++ {
				c.Memory[c.IndexRegisters+uint16(i)] = uint16(c.Registers[i])
			}

			//c.IndexRegisters = (vX) + 1
			c.ProgramCounter += 2

		case 0x65:
			for i := 0; i <= int(vX); i++ {
				c.Registers[i] = uint8(c.Memory[c.IndexRegisters+uint16(i)])
			}

			//c.IndexRegisters = (vX) + 1
			c.ProgramCounter += 2

		default:
			fmt.Printf("Something broke 0xF000")

		}

	default:
		fmt.Printf("error OOF")
	}

	if c.DelayTimer > 0 {
		c.DelayTimer = c.DelayTimer - 1
	}

	if c.SoundTimer > 0 {
		if c.SoundTimer == 1 {
			//c.beeper() need to impliment this for sound
			fmt.Println("Beep")
		}
		c.SoundTimer = c.SoundTimer - 1
	}
}
