package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpcode(t *testing.T) {
	// 0x00E0
	// 0x00EE
	// 0x1000: // 1NNN	Jumps to address NNN.
	// 0x2000: // 2NNN	Calls subroutine at NNN.
	//0x3000: // 3XNN	Skips t    he next instruction if VX equals NN.

	//--- Test 00E0 clear ----
	chp := InitChip8([]byte{0x00, 0xE0}, map[uint16]byte{})
	chp.Videomem[1][2] = 1
	upDisp, errExec := chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, true, upDisp)
	assert.Equal(t, byte(0), chp.Videomem[1][2])
	//RET
	chp = InitChip8([]byte{0x00, 0xEE}, map[uint16]byte{})
	chp.SP = 3
	chp.Stack[chp.SP] = 0xDEAD

	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0xDEAD+2), chp.PC)

	//0x00FE: // Disable extended screen mode
	chp = InitChip8([]byte{0x00, 0xFE}, map[uint16]byte{})
	chp.HiRes = true
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, true, upDisp)
	assert.Equal(t, false, chp.HiRes)

	//0x00FF:
	chp = InitChip8([]byte{0x00, 0xFF}, map[uint16]byte{})
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, true, upDisp)
	assert.Equal(t, true, chp.HiRes)

	//0x00CN
	chp = InitChip8([]byte{0x00, 0xC5}, map[uint16]byte{})
	for i := 0; i < 64; i++ {
		for j := 0; j < i; j++ {
			chp.Videomem[j][i] = 1
		}
	}
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, true, upDisp)
	for x := 0; x < 128; x++ {
		for y := 0; y < 64; y++ {
			if (x + 5) < y {
				assert.Equal(t, uint8(1), chp.Videomem[x][y])
			} else {
				assert.Equal(t, uint8(0), chp.Videomem[x][y])
			}
		}
	}
	//0x1000: // 1NNN	Jumps to address NNN.
	chp = InitChip8([]byte{0x12, 0x34}, map[uint16]byte{})
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x234), chp.PC)

	//0x2000: // 2NNN	Calls subroutine at NNN.
	chp = InitChip8([]byte{0x21, 0x23}, map[uint16]byte{})
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x123), chp.PC)
	assert.Equal(t, byte(1), chp.SP)

	// 0x3000: // 3XNN	Skips t    he next instruction if VX equals NN.
	chp = InitChip8([]byte{0x32, 0x42}, map[uint16]byte{})
	chp.Reg[2] = 0x42
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x200+4), chp.PC)
	chp = InitChip8([]byte{0x32, 0x42}, map[uint16]byte{})
	chp.Reg[2] = 0x41
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x200+2), chp.PC)
	//0x4000: // 4XNN	Skips the next instruction if VX doesn't equal NN.
	chp = InitChip8([]byte{0x42, 0x42}, map[uint16]byte{})
	chp.Reg[2] = 0x41
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x200+4), chp.PC)
	chp = InitChip8([]byte{0x42, 0x42}, map[uint16]byte{})
	chp.Reg[2] = 0x42
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x200+2), chp.PC)
	// 0x5000: // 5XY0	Skips the next instruction if VX equals VY.
	chp = InitChip8([]byte{0x51, 0x20}, map[uint16]byte{})
	chp.Reg[1] = 0x12
	chp.Reg[2] = 0x12
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x200+4), chp.PC)
	chp = InitChip8([]byte{0x51, 0x20}, map[uint16]byte{})
	chp.Reg[1] = 0x10
	chp.Reg[2] = 0x12
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x200+2), chp.PC)
	//0x6000: // 6XNN	Sets VX to NN.
	chp = InitChip8([]byte{0x62, 0x69}, map[uint16]byte{})
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, byte(0x69), chp.Reg[2])
	//0x7000: // 7XNN	Adds NN to VX.
	chp = InitChip8([]byte{0x72, 0x12}, map[uint16]byte{})
	chp.Reg[2] = 0x10
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, byte(0x10+0x12), chp.Reg[2])
	//0x8000:
	// 8XY0	Sets VX to the value of VY.
	// 8XY2	Sets VX to VX and VY.
	// 8XY3	Sets VX to VX xor VY.
	// 8XY4	Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
	chp = InitChip8([]byte{0x81, 0x24}, map[uint16]byte{})
	chp.Reg[2] = 0xFE
	chp.Reg[1] = 0x3
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, byte(1), chp.Reg[0xF])

	chp = InitChip8([]byte{0x81, 0x24}, map[uint16]byte{})
	chp.Reg[2] = 0xFE
	chp.Reg[1] = 0x1
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, byte(0), chp.Reg[0xF])

	// 8XY5	VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
	// 8X06 Set Vx = Vx SHR 1.If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
	// 8XY7  SUBN Vx, Vy - Sub Vx from Vy into Vx. Sets VF
	// 8xyE SHL Vx, Vy - Most sig bit of Vx into VF, shift Vx left to mult by 2

	chp = InitChip8([]byte{0xC2, 0x03}, map[uint16]byte{}) // CXNN	Sets VX to a random number and NN.
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)

	// DXYN	Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels
	chp = InitChip8([]byte{0xD2, 0x32}, map[uint16]byte{})
	chp.I = 0x300
	chp.Mem[0x300] = 0xFF
	chp.Mem[0x301] = 1<<7 | 1
	chp.Reg[2] = 13
	chp.Reg[3] = 5
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, true, upDisp)

	for i := 0; i < 8; i++ {
		assert.Equal(t, byte(1), chp.Videomem[13+i][5])
	}
	assert.Equal(t, byte(0), chp.Videomem[12][5])
	assert.Equal(t, byte(0), chp.Videomem[13+8][5])
	assert.Equal(t, byte(0), chp.Videomem[13+6][4])

	assert.Equal(t, byte(0), chp.Videomem[13-1][6])
	assert.Equal(t, byte(1), chp.Videomem[13][6])
	assert.Equal(t, byte(0), chp.Videomem[13+1][6])

	assert.Equal(t, byte(0), chp.Videomem[13+6][6])
	assert.Equal(t, byte(1), chp.Videomem[13+7][6])
	assert.Equal(t, byte(0), chp.Videomem[13+8][6])

	assert.Equal(t, uint16(0x300), chp.I)

	// EXA1	Skips the next instruction if the key stored in VX isn't pressed.
	chp = InitChip8([]byte{0xE3, 0xA1}, map[uint16]byte{})
	chp.Reg[3] = 6
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x204), chp.PC)

	chp = InitChip8([]byte{0xE3, 0xA1}, map[uint16]byte{})
	chp.Reg[3] = 6
	chp.Keyboard[6] = true
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x202), chp.PC)

	//0xE09E // Skips the next instruction if the key stored in VX is pressed
	chp = InitChip8([]byte{0xE3, 0x9E}, map[uint16]byte{})
	chp.Reg[3] = 6
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x202), chp.PC)

	chp = InitChip8([]byte{0xE3, 0x9E}, map[uint16]byte{})
	chp.Reg[3] = 6
	chp.Keyboard[6] = true
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x204), chp.PC)

	// FX07 Set Vx = delay timer value.
	chp = InitChip8([]byte{0xF1, 0x07}, map[uint16]byte{})
	chp.DT = 123
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, byte(123), chp.Reg[1])

	// FX0A	A key press is awaited, and then stored in VX.
	chp = InitChip8([]byte{0xF1, 0x0A}, map[uint16]byte{})
	chp.Reg[1] = 3
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x200), chp.PC)

	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x200), chp.PC)

	chp.Keyboard[8] = true
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, uint16(0x202), chp.PC)
	assert.Equal(t, byte(8), chp.Reg[1])

	//0x18: // FX18	Set sound timer = Vx.
	chp = InitChip8([]byte{0xF1, 0x18}, map[uint16]byte{})
	chp.Reg[1] = 3
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, byte(3), chp.ST)

	// FX29	Sets I to the location of the sprite for the character in VX.
	chp = InitChip8([]byte{0xF0, 0x29}, map[uint16]byte{})
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)

	// FX33	Stores BCD
	chp = InitChip8([]byte{0xF2, 0x33}, map[uint16]byte{})
	chp.Reg[2] = 0xFF
	chp.I = 0x220
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, byte(2), chp.Mem[0x220])
	assert.Equal(t, byte(5), chp.Mem[0x221])
	assert.Equal(t, byte(5), chp.Mem[0x222])
	assert.Equal(t, byte(0), chp.Mem[0x223])

	// FX55	Stores V0 to VX in memory starting at address I.
	chp = InitChip8([]byte{0xF2, 0x55}, map[uint16]byte{})
	chp.Reg[0] = 0x03
	chp.Reg[1] = 0x02
	chp.Reg[2] = 0x01

	chp.I = 0x220
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, byte(3), chp.Mem[0x220])
	assert.Equal(t, byte(2), chp.Mem[0x221])
	assert.Equal(t, byte(1), chp.Mem[0x222])
	assert.Equal(t, byte(0), chp.Mem[0x223])

	// FX65	Fills V0 to VX with values from memory starting at address I.
	chp = InitChip8([]byte{0xF2, 0x65}, map[uint16]byte{})
	chp.I = 0x200
	upDisp, errExec = chp.ExecOp()
	assert.Equal(t, nil, errExec)
	assert.Equal(t, false, upDisp)
	assert.Equal(t, byte(0xF2), chp.Reg[0])
	assert.Equal(t, byte(0x65), chp.Reg[1])
	assert.Equal(t, byte(0), chp.Reg[2])
}
