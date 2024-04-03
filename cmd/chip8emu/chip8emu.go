package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"time"
)

//var disp wiointerface.Display_ILI9341

// Avoid blowing up stack, really bad coding style!
type Chip8 struct {
	Mem      []uint8  // [0x1000]uint8
	Reg      []uint8  //[0x10]uint8
	Stack    []uint16 //[0x10]uint16
	I        uint16
	PC       uint16
	SP       uint8
	DT       uint8
	ST       uint8
	Videomem [][]byte //[64][32]bool
	Keyboard []bool   //[16]bool
	Run      bool     //Mahdollistaa pysÃÂ¤ytyksen

	HiRes bool

	prevTickTime time.Time

	colorGuide map[uint16]byte
}

// ac *bwGoScreenManager.AppConnection

func InitChip8(code []byte, colorGuide map[uint16]byte) Chip8 {
	chip := Chip8{
		Mem:          make([]uint8, 0x1000),
		Reg:          make([]uint8, 0x10),
		Stack:        make([]uint16, 0x10),
		PC:           0x200,
		Videomem:     make([][]byte, 128),
		Keyboard:     make([]bool, 16),
		Run:          true,
		prevTickTime: time.Now(),
		colorGuide:   colorGuide,
	}

	for x := range chip.Videomem {
		chip.Videomem[x] = make([]byte, 64)
	}

	font := []uint8{
		0xf0, 0x90, 0x90, 0x90, 0xf0,
		0x20, 0x60, 0x20, 0x20, 0x70,
		0xf0, 0x10, 0xf0, 0x80, 0xf0,
		0xf0, 0x10, 0xf0, 0x10, 0xf0,
		0x90, 0x90, 0xf0, 0x10, 0x10,
		0xf0, 0x80, 0xf0, 0x10, 0xf0,
		0xf0, 0x80, 0xf0, 0x90, 0xf0,
		0xf0, 0x10, 0x20, 0x40, 0x40,
		0xf0, 0x90, 0xf0, 0x90, 0xf0,
		0xf0, 0x90, 0xf0, 0x10, 0xf0,
		0xf0, 0x90, 0xf0, 0x90, 0x90,
		0xe0, 0x90, 0xe0, 0x90, 0xe0,
		0xf0, 0x80, 0x80, 0x80, 0xf0,
		0xe0, 0x90, 0x90, 0x90, 0xe0,
		0xf0, 0x80, 0xf0, 0x80, 0xf0,
		0xf0, 0x80, 0xf0, 0x80, 0x80,
	}

	for i := 0; i < len(font); i++ {
		chip.Mem[i] = font[i]
	}

	for i, v := range code {
		chip.Mem[0x200+i] = v
	}
	return chip
}

func (p *Chip8) UpdateTimer(tNow time.Time) {
	changeSteps := p.prevTickTime.Sub(tNow).Milliseconds() / 16
	if 0 < changeSteps {
		if 0 < p.DT {
			p.DT -= uint8(changeSteps)
		}
		if 0 < p.ST {
			p.ST -= uint8(changeSteps)
		}
		p.prevTickTime = tNow
	}
}

func (p *Chip8) ScrollDown(n int) {
	//can scroll only down..
	fmt.Printf("scrolling down %v\n", n)

	for y := 64 - 1; y >= n; y-- {
		for x := 0; x < 128; x++ {
			p.Videomem[x][y] = p.Videomem[x][y-n]
		}
	}
	// Wipe the remaining top n rows of pixels
	for y := 0; y < n; y++ {
		for x := 0; x < 128; x++ {
			p.Videomem[x][y] = 0
		}
	}

}

func (p *Chip8) OpCodeNow() uint16 {
	return uint16(p.Mem[p.PC])<<8 + uint16(p.Mem[p.PC+1])
}

func OpCodeClass(opcode uint16) uint16 {
	switch opcode & 0xF000 {
	case 0x0000:
		return opcode
	case 0x1000, 0x2000, 0x3000, 0x4000, 0x5000, 0x6000, 0x7000:
		return opcode & 0xF000
	case 0x8000:
		return opcode & 0xF00F
	case 0x9000, 0xA000, 0xB000, 0xC000, 0xD000:
		return opcode & 0xF000
	case 0xF000:
		return opcode & 0xF0FF
	}
	return 0
}

func (p *Chip8) execOp8000() error {
	opcode := uint16(p.Mem[p.PC])<<8 + uint16(p.Mem[p.PC+1])

	X := byte((opcode & 0x0F00) >> 8)
	Y := byte((opcode & 0x00F0) >> 4)
	NN := byte(opcode & 0x00FF)

	switch opcode & 0xF {
	case 0x0: // 8XY0	Sets VX to the value of VY.
		p.Reg[X] = p.Reg[Y]
	case 0x1: // 8XY1	Sets VX to VX or VY.
		p.Reg[X] = p.Reg[X] | p.Reg[Y]
	case 0x2: // 8XY2	Sets VX to VX and VY.
		p.Reg[X] = p.Reg[X] & p.Reg[Y]
	case 0x3: // 8XY3	Sets VX to VX xor VY.
		p.Reg[X] = p.Reg[X] ^ p.Reg[Y]
	case 0x4: // 8XY4	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
		sum := uint16(p.Reg[X]) + uint16(p.Reg[Y])
		p.Reg[0xF] = 0
		if 0xFF < sum {
			p.Reg[0xF] = 1
		}
		p.Reg[X] = NN

	case 0x5: // 8XY5	VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
		//TODO BUGI? sub := Reg[X] + Reg[Y]
		if p.Reg[X] < p.Reg[Y] {
			p.Reg[0xF] = 0
		} else {
			p.Reg[0xF] = 1
		}
		p.Reg[X] = p.Reg[X] - p.Reg[Y]
	case 0x6: // Set Vx = Vx SHR 1.If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
		//MODERN
		p.Reg[0xF] = p.Reg[X] & 1
		p.Reg[X] >>= 1
		//DEFAULT
	//	p.Reg[0xF] = p.Reg[Y] & 1
	//	p.Reg[X]= p.Reg[Y]

	case 0x7: // SUBN Vx, Vy - Sub Vx from Vy into Vx. Sets VF
		p.Reg[0xF] = 1
		if p.Reg[Y] < p.Reg[X] {
			p.Reg[0xF] = 0
		}
		p.Reg[X] = p.Reg[Y] - p.Reg[X]
	case 0x000E: //SHL Vx, Vy - Most sig bit of Vx into VF, shift Vx left to mult by 2

		//MODERN
		p.Reg[0xF] = p.Reg[X] >> 7
		p.Reg[X] <<= 1
		//DEFAULT
		/*p.Reg[0xF] = p.Reg[Y] >> 7
		p.Reg[X] = p.Reg[Y] << 2*/

	default:
		return fmt.Errorf("Unknown OP 0x%04X", opcode)
	}
	return nil
}

/*
func (p *Chip8) execOp0000() error {
	opcode := uint16(p.Mem[p.PC])<<8 + uint16(p.Mem[p.PC+1])
}*/

// Need to update display and erro return
func (p *Chip8) ExecOp() (bool, error) {
	opcode := uint16(p.Mem[p.PC])<<8 + uint16(p.Mem[p.PC+1])
	//fmt.Printf("OP= 0x%04X\n", opcode)

	nextPC := p.PC + 2 //normaalisti

	X := byte((opcode & 0x0F00) >> 8)
	Y := byte((opcode & 0x00F0) >> 4)

	N := int(opcode & 0x000F)
	NN := byte(opcode & 0x00FF)
	NNN := uint16(opcode & 0x0FFF)
	//fmt.Printf("X=%v,Y=%v ja op oli %#X\n", X, Y, opcode)

	updateDisplay := false
	switch opcode & 0xF000 {
	/*
		OPCODE 0
	*/
	case 0x0000: // 0x0***
		switch opcode {
		case 0x00E0: //CLS
			for x := 0; x < len(p.Videomem); x++ {
				for y := 0; y < len(p.Videomem[x]); y++ {
					p.Videomem[x][y] = 0
				}
			}
			updateDisplay = true
		case 0x00EE: //RET
			if p.SP <= 0 {
				return false, fmt.Errorf("stack bottom")
			}
			nextPC = p.Stack[p.SP] + 2
			p.SP--

			/*
				case 0x00FB: // Scroll display 4 pixels right
					for y := 0; y < 64; y++ {
						var x int
						for x = 128 - 1; x >= 4; x-- {
							p.Videomem[x-4][y] = p.Videomem[x][y]
						}
						p.Videomem[3][y] = 0
						p.Videomem[2][y] = 0
						p.Videomem[1][y] = 0
						p.Videomem[0][y] = 0
					}
					updateDisplay = true
				case 0x00FC: // Scroll display 4 pixels left
					for y := 0; y < 64; y++ {
						var x int
						for x = 0; x < 128-4; x++ {
							p.Videomem[x][y] = p.Videomem[x+4][y]
						}
					}
					for y := 0; y < 64; y++ {
						p.Videomem[124][y] = 0
						p.Videomem[125][y] = 0
						p.Videomem[126][y] = 0
						p.Videomem[127][y] = 0
					}
					updateDisplay = true
			*/
		case 0x00FD: // Exit CHIP interpreter
			return true, fmt.Errorf("EXIT")
		case 0x00FE: // Disable extended screen mode
			p.HiRes = false
			updateDisplay = true
		case 0x00FF: // Enable extended screen mode for full-screen graphics
			p.HiRes = true
			updateDisplay = true
		default:
			if opcode&0xFFF0 == 0x00C0 { //0x00CN: //Scroll display N lines down
				p.ScrollDown(N)
				updateDisplay = true
			} else {
				return false, fmt.Errorf("Unknown OP 0x%04X  (last=%v)", opcode, opcode)
			}
		}
	case 0x1000: // 1NNN	Jumps to address NNN.
		nextPC = NNN
	case 0x2000: // 2NNN	Calls subroutine at NNN.
		if p.SP >= 15 {
			return false, fmt.Errorf("stack exhausted")
		}
		p.SP++
		p.Stack[p.SP] = p.PC
		nextPC = NNN
	case 0x3000: // 3XNN	Skips t    he next instruction if VX equals NN.
		if p.Reg[X] == NN {
			nextPC += 2
		}
	case 0x4000: // 4XNN	Skips the next instruction if VX doesn't equal NN.
		if p.Reg[X] != NN {
			nextPC += 2
		}
	case 0x5000: // 5XY0	Skips the next instruction if VX equals VY.
		if opcode&0xF != 0 {
			return false, fmt.Errorf("invalid opcode 0x%04X\n", opcode)
		}
		if p.Reg[X] == p.Reg[Y] {
			nextPC += 2
		}
	case 0x6000: // 6XNN	Sets VX to NN.
		p.Reg[X] = NN
	case 0x7000: // 7XNN	Adds NN to VX.
		p.Reg[X] = p.Reg[X] + NN
	case 0x8000:
		err8000 := p.execOp8000()
		if err8000 != nil {
			return false, err8000
		}

	case 0x9000: // 9XY0	Skips the next instruction if VX doesn't equal VY.
		if p.Reg[X] != p.Reg[Y] {
			nextPC += 2
		}
	case 0xA000: // ANNN	Sets I to the address NNN.
		p.I = NNN
	case 0xB000: // BNNN	Jumps to the address NNN plus V0.
		nextPC = NNN + uint16(p.Reg[0])
	case 0xC000: // CXNN	Sets VX to a random number and NN.
		p.Reg[X] = byte(uint16(rand.Uint32())) & NN
	case 0xD000:
		// DXYN	Draws a sprite at coordinate (VX, VY) that has a width of 8
		// pixels and a height of N pixels. Each row of 8 pixels is read as
		// bit-coded (with the most significant bit of each byte displayed on
		// the left) starting from memory location I; I value doesn't change
		// after the execution of this instruction.

		//w := len(p.Videomem)
		//h := len(p.Videomem[0])

		w := 64
		h := 32
		if p.HiRes {
			w = 128
			h = 64
		}

		Xo := int(p.Reg[X])
		Yo := int(p.Reg[Y])

		guideColor, hazGuide := p.colorGuide[p.I]
		if !hazGuide {
			guideColor = 1 //Default
		}

		//On super chip, N=0 draws 16x16 bitmap  2bytes per row
		p.Reg[0xF] = 0

		if N == 0 {
			fmt.Printf("DRAW SUPERCHIP PIC\n")
			for row := 0; row < 16; row++ {
				rowdata1 := p.Mem[int(p.I)+row*2]
				rowdata2 := p.Mem[int(p.I)+row*2+1]
				yp := (Yo + row) % 64
				for a := 0; a < 8; a++ {
					b1 := rowdata1&(byte(0x80)>>a) != 0
					xp1 := (Xo + a) % 128
					b2 := rowdata2&(byte(0x80)>>a) != 0
					xp2 := (xp1 + 8) % 128

					fmt.Printf("xp1=%v xp2=%v\n", xp1, xp2)

					if b1 {
						if 0 < p.Videomem[xp1][yp] {
							p.Reg[0xF] = 1
							p.Videomem[xp1][yp] = 0 //always clear
						} else {
							p.Videomem[xp1][yp] = guideColor
						}
					}
					if b2 {
						if 0 < p.Videomem[xp2][yp] {
							p.Reg[0xF] = 1
							p.Videomem[xp2][yp] = 0 //always clear
						} else {
							p.Videomem[xp2][yp] = guideColor
						}
					}

				}
			}
		} else {

			for yoff := 0; yoff < N; yoff++ {
				rowdata := p.Mem[int(p.I)+yoff]
				yp := (Yo + yoff) % h
				for xoff := 0; xoff < 8; xoff++ {
					b := rowdata&(byte(0x80)>>xoff) != 0
					if b {
						xp := (Xo + xoff) % w
						if 0 < p.Videomem[xp][yp] {
							p.Reg[0xF] = 1
							p.Videomem[xp][yp] = 0 //always clear
						} else {
							p.Videomem[xp][yp] = guideColor
						}
						//fmt.Printf("changed x:%v y:%v\n", xp, yp)
					}
				}
				//fmt.Printf("\n\n")
			}
		}

		updateDisplay = true
	case 0xE000: //
		//fmt.Printf("\n-----------TESTATAAN NAPPIA %v-------------\n,", p.Reg[X]&0xF)
		switch opcode & 0x00FF {
		case 0x9E:
			if p.Keyboard[p.Reg[X]&0xF] {
				nextPC += 2
			}
		case 0xA1: // EXA1	Skips the next instruction if the key stored in VX isn't pressed.
			if (p.Keyboard[p.Reg[X]&0xF]) == false { //TODO mutex??
				nextPC += 2
			}
		default:
			return false, fmt.Errorf("Invalid op code 0x%04X", opcode)
		}
	case 0xF000:
		switch byte(opcode & 0x00FF) {
		case 0x07: // FX07 Set Vx = delay timer value.
			p.Reg[X] = p.DT
		case 0x0A: // FX0A	A key press is awaited, and then stored in VX.
			//fmt.Printf("\n-----------OOTETAAN ANYKEY-------------\n,")
			fmt.Printf("WAIT KEY!!\n")
			haskey := false
			for i := 0; i < 16; i++ {
				if p.Keyboard[i] {
					p.Reg[X] = uint8(i)
					haskey = true
					break
				}
			}
			if !haskey {
				return false, nil //Jumps back to same op
			}

		case 0x15: // FX15	Set delay timer = Vx.
			p.DT = p.Reg[X]
		case 0x18: // FX18	Set sound timer = Vx.
			p.ST = p.Reg[X]

		case 0x1E: // FX1E	Adds VX to I.
			// VF is set to 1 when range overflow (I+VX>0xFFF), and 0 when

			p.I = p.I + uint16(p.Reg[X])

			/*p.Reg[0xF] = 0x0
			if p.I > 0xFFF {
				p.Reg[0xF] = 1
			}
			*/

		case 0x29: // FX29	Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
			//TODO ONKO BUGI? I = 0x50 + uint16(p.Reg[X])*4 // XXX assuming fontset at the beginning of memory.
			p.I = uint16(p.Reg[X]) * 5

		//case 0x0030	FX30 (i := bighex vx) Set i to a large hexadecimal character based on the value of vx.

		case 0x33:
			// FX33	Stores the Binary-coded decimal representation of VX, with
			// the most significant of three digits at the address in I, the
			// middle digit at I plus 1, and the least significant digit at I
			// plus 2.
			r := p.Reg[X]
			p.Mem[p.I] = r / 100
			p.Mem[p.I+1] = (r % 100) / 10
			p.Mem[p.I+2] = r % 10
			//TODO RATKAISE TESTILLÄ!

		case 0x0055:
			// FX55	Stores V0 to VX in memory starting at address I.
			// On the original interpreter, when the operation is done, I=I+X+1
			//for i, reg := range p.Reg[0:X] {
			for i := 0; i <= int(X); i++ {
				p.Mem[p.I+uint16(i)] = p.Reg[i]
			}
			p.I = uint16(X + 1)
		case 0x0065:
			// FX65	Fills V0 to VX with values from memory starting at address I.
			for i := uint16(0); i <= uint16(X); i++ {
				p.Reg[i] = p.Mem[p.I+i]
			}
			p.I = uint16(X + 1)

		//case 0x0075: FX75 (saveflags vx) Save v0-vX to flag registers.
		//case 0x0085: FX85 (loadflags vx) Restore v0-vX from flag registers.

		default:
			//fmt.Printf("TUNTEMATON OP KOODI")
			return false, fmt.Errorf("Unknown OP 0x%04X", opcode)
		}
	default:
		//fmt.Printf("TUNTEMATON OP KOODI")
		return false, fmt.Errorf("Unknown OP 0x%04X", opcode)

	}

	if nextPC == p.PC {
		p.Run = false
		return false, fmt.Errorf("jamming loop")
	}

	p.PC = nextPC
	return updateDisplay, nil
}

/*
const (
	PIXEOFFCOLOR uint16 = 0
	PIXEONCOLOR  uint16 = 0x1F
)

func runchip(code []byte, keymask []uint32) {
	//var doDisplayUpdate bool
	chp := InitChip8(code)

	timeRunned := time.Now()
	for {
		chp.SpyKeys()
		chp.SpySprites()
		doDisplayUpdate, errExec := chp.ExecOp()
		if errExec != nil {
			//fmt.Printf("EXEC ERR %s\n", errExec)
			return
		}
		if doDisplayUpdate {
			disp.SetWindow(0, 0, 64*4, 32*4)
			disp.StartWrite()
			for y := 0; y < 32; y++ {
				for rep := 0; rep < 4; rep++ {
					for x := 0; x < 64; x++ {
						if chp.Videomem[x][y] {
							disp.Write16bit([]uint16{PIXEONCOLOR, PIXEONCOLOR, PIXEONCOLOR, PIXEONCOLOR})
						} else {
							disp.Write16bit([]uint16{PIXEOFFCOLOR, PIXEOFFCOLOR, PIXEOFFCOLOR, PIXEOFFCOLOR})
						}
					}
				}
			}
			disp.EndWrite()

		}

		for time.Since(timeRunned) < time.Millisecond*2 {

			keys := disp.GetWioKeys()

			for i, mask := range keymask {
				chp.Keyboard[i] = 0 < keys&mask
			}

		}

		timeRunned = time.Now()
		chp.UpdateTimer(time.Now())
	}
}
*/
