/*
Spy sprite and button usage.
Used as tool on PC
*/
package main

import (
	"fmt"
	"strings"
)

type ChipSpy struct {
	UsedKeys    [16]bool
	UsedSprites map[uint16]byte
	UsedOpCodes map[uint16]int64
}

func InitChipSpy() ChipSpy {
	return ChipSpy{
		UsedSprites: make(map[uint16]byte),
		UsedOpCodes: make(map[uint16]int64),
	}
}

func (a ChipSpy) String() string {
	var sb strings.Builder
	sb.WriteString("KEYS:")
	for i, b := range a.UsedKeys {
		if b {
			sb.WriteString(fmt.Sprintf(" %v", i))
		}
	}
	sb.WriteString("\nSPRITES:")
	for addr, rows := range a.UsedSprites {
		sb.WriteString(fmt.Sprintf("0x%02X (%v)\n", addr, rows))
	}
	return sb.String()
}

// For spying for assigning keys and sprites
func (p *ChipSpy) SpyKeys(chip *Chip8) bool { // print keys and sprite addesses
	opcode := uint16(chip.Mem[chip.PC])<<8 + uint16(chip.Mem[chip.PC+1])

	if (opcode&0xF0FF == 0xE0A1) || (opcode&0xF0FF == 0xE09E) {
		keynumber := chip.Reg[(opcode>>8)&0xF] & 0xF
		if !p.UsedKeys[keynumber] {
			p.UsedKeys[keynumber] = true
			return true
		}
	}
	return false
}

func debugPrintSprite(data []byte) {
	fmt.Printf("¤¤¤¤¤¤¤\n")
	for _, row := range data {
		fmt.Printf("%08b\n", row)
	}
	fmt.Printf("¤¤¤¤¤¤¤\n")

}

func (p *ChipSpy) SpySprites(chip *Chip8) bool {
	opcode := uint16(chip.Mem[chip.PC])<<8 + uint16(chip.Mem[chip.PC+1])
	n := byte(opcode & 0x000F)
	if opcode&0xF000 == 0xD000 {
		nInSprite, haz := p.UsedSprites[chip.I]
		if !haz {
			p.UsedSprites[chip.I] = n
			debugPrintSprite(chip.Mem[chip.I:(chip.I + uint16(n))])
			return true
		}
		if nInSprite != n {
			p.UsedSprites[chip.I] = n
			debugPrintSprite(chip.Mem[chip.I:(chip.I + uint16(n))])
			return true
		}
	}
	return false
}

func (p *ChipSpy) SpyCode(chip *Chip8) bool {
	key := OpCodeClass(chip.OpCodeNow())

	n, haz := p.UsedOpCodes[key]
	if !haz {
		p.UsedOpCodes[key] = 1
		return true //NEW
	}
	p.UsedOpCodes[key] = n + 1
	return false
}
