package hardware

func assignLookup(arr *[256]instruction) {

	*arr = [256]instruction{

		{"BRK", (*Olc5602).brk, (*Olc5602).imm, false, 7}, {"ORA", (*Olc5602).ora, (*Olc5602).izx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 3}, {"ORA", (*Olc5602).ora, (*Olc5602).zp0, false, 3}, {"ASL", (*Olc5602).asl, (*Olc5602).zp0, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5}, {"PHP", (*Olc5602).php, (*Olc5602).imp, true, 3}, {"ORA", (*Olc5602).ora, (*Olc5602).imm, false, 2}, {"ASL", (*Olc5602).asl, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"ORA", (*Olc5602).ora, (*Olc5602).abs, false, 4}, {"ASL", (*Olc5602).asl, (*Olc5602).abs, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6},
		{"BPL", (*Olc5602).bpl, (*Olc5602).rel, false, 2}, {"ORA", (*Olc5602).ora, (*Olc5602).izy, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"ORA", (*Olc5602).ora, (*Olc5602).zpx, false, 4}, {"ASL", (*Olc5602).asl, (*Olc5602).zpx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6}, {"CLC", (*Olc5602).clc, (*Olc5602).imp, true, 2}, {"ORA", (*Olc5602).ora, (*Olc5602).aby, false, 4}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"ORA", (*Olc5602).ora, (*Olc5602).abx, false, 4}, {"ASL", (*Olc5602).asl, (*Olc5602).abx, false, 7}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7},
		{"JSR", (*Olc5602).jsr, (*Olc5602).abs, false, 6}, {"AND", (*Olc5602).and, (*Olc5602).izx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"BIT", (*Olc5602).bit, (*Olc5602).zp0, false, 3}, {"AND", (*Olc5602).and, (*Olc5602).zp0, false, 3}, {"ROL", (*Olc5602).rol, (*Olc5602).zp0, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5}, {"PLP", (*Olc5602).plp, (*Olc5602).imp, true, 4}, {"AND", (*Olc5602).and, (*Olc5602).imm, false, 2}, {"ROL", (*Olc5602).rol, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"BIT", (*Olc5602).bit, (*Olc5602).abs, false, 4}, {"AND", (*Olc5602).and, (*Olc5602).abs, false, 4}, {"ROL", (*Olc5602).rol, (*Olc5602).abs, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6},
		{"BMI", (*Olc5602).bmi, (*Olc5602).rel, false, 2}, {"AND", (*Olc5602).and, (*Olc5602).izy, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"AND", (*Olc5602).and, (*Olc5602).zpx, false, 4}, {"ROL", (*Olc5602).rol, (*Olc5602).zpx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6}, {"SEC", (*Olc5602).sec, (*Olc5602).imp, true, 2}, {"AND", (*Olc5602).and, (*Olc5602).aby, false, 4}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"AND", (*Olc5602).and, (*Olc5602).abx, false, 4}, {"ROL", (*Olc5602).rol, (*Olc5602).abx, false, 7}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7},
		{"RTI", (*Olc5602).rti, (*Olc5602).imp, true, 6}, {"EOR", (*Olc5602).eor, (*Olc5602).izx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 3}, {"EOR", (*Olc5602).eor, (*Olc5602).zp0, false, 3}, {"LSR", (*Olc5602).lsr, (*Olc5602).zp0, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5}, {"PHA", (*Olc5602).pha, (*Olc5602).imp, true, 3}, {"EOR", (*Olc5602).eor, (*Olc5602).imm, false, 2}, {"LSR", (*Olc5602).lsr, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"JMP", (*Olc5602).jmp, (*Olc5602).abs, false, 3}, {"EOR", (*Olc5602).eor, (*Olc5602).abs, false, 4}, {"LSR", (*Olc5602).lsr, (*Olc5602).abs, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6},
		{"BVC", (*Olc5602).bvc, (*Olc5602).rel, false, 2}, {"EOR", (*Olc5602).eor, (*Olc5602).izy, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"EOR", (*Olc5602).eor, (*Olc5602).zpx, false, 4}, {"LSR", (*Olc5602).lsr, (*Olc5602).zpx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6}, {"CLI", (*Olc5602).cli, (*Olc5602).imp, true, 2}, {"EOR", (*Olc5602).eor, (*Olc5602).aby, false, 4}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"EOR", (*Olc5602).eor, (*Olc5602).abx, false, 4}, {"LSR", (*Olc5602).lsr, (*Olc5602).abx, false, 7}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7},
		{"RTS", (*Olc5602).rts, (*Olc5602).imp, true, 6}, {"ADC", (*Olc5602).adc, (*Olc5602).izx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 3}, {"ADC", (*Olc5602).adc, (*Olc5602).zp0, false, 3}, {"ROR", (*Olc5602).ror, (*Olc5602).zp0, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5}, {"PLA", (*Olc5602).pla, (*Olc5602).imp, true, 4}, {"ADC", (*Olc5602).adc, (*Olc5602).imm, false, 2}, {"ROR", (*Olc5602).ror, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"JMP", (*Olc5602).jmp, (*Olc5602).ind, false, 5}, {"ADC", (*Olc5602).adc, (*Olc5602).abs, false, 4}, {"ROR", (*Olc5602).ror, (*Olc5602).abs, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6},
		{"BVS", (*Olc5602).bvs, (*Olc5602).rel, false, 2}, {"ADC", (*Olc5602).adc, (*Olc5602).izy, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"ADC", (*Olc5602).adc, (*Olc5602).zpx, false, 4}, {"ROR", (*Olc5602).ror, (*Olc5602).zpx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6}, {"SEI", (*Olc5602).sei, (*Olc5602).imp, true, 2}, {"ADC", (*Olc5602).adc, (*Olc5602).aby, false, 4}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"ADC", (*Olc5602).adc, (*Olc5602).abx, false, 4}, {"ROR", (*Olc5602).ror, (*Olc5602).abx, false, 7}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7},
		{"???", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"STA", (*Olc5602).sta, (*Olc5602).izx, false, 6}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6}, {"STY", (*Olc5602).sty, (*Olc5602).zp0, false, 3}, {"STA", (*Olc5602).sta, (*Olc5602).zp0, false, 3}, {"STX", (*Olc5602).stx, (*Olc5602).zp0, false, 3}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 3}, {"DEY", (*Olc5602).dey, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"TXA", (*Olc5602).txa, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"STY", (*Olc5602).sty, (*Olc5602).abs, false, 4}, {"STA", (*Olc5602).sta, (*Olc5602).abs, false, 4}, {"STX", (*Olc5602).stx, (*Olc5602).abs, false, 4}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 4},
		{"BCC", (*Olc5602).bcc, (*Olc5602).rel, false, 2}, {"STA", (*Olc5602).sta, (*Olc5602).izy, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6}, {"STY", (*Olc5602).sty, (*Olc5602).zpx, false, 4}, {"STA", (*Olc5602).sta, (*Olc5602).zpx, false, 4}, {"STX", (*Olc5602).stx, (*Olc5602).zpy, false, 4}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 4}, {"TYA", (*Olc5602).tya, (*Olc5602).imp, true, 2}, {"STA", (*Olc5602).sta, (*Olc5602).aby, false, 5}, {"TXS", (*Olc5602).txs, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 5}, {"STA", (*Olc5602).sta, (*Olc5602).abx, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5},
		{"LDY", (*Olc5602).ldy, (*Olc5602).imm, false, 2}, {"LDA", (*Olc5602).lda, (*Olc5602).izx, false, 6}, {"LDX", (*Olc5602).ldx, (*Olc5602).imm, false, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6}, {"LDY", (*Olc5602).ldy, (*Olc5602).zp0, false, 3}, {"LDA", (*Olc5602).lda, (*Olc5602).zp0, false, 3}, {"LDX", (*Olc5602).ldx, (*Olc5602).zp0, false, 3}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 3}, {"TAY", (*Olc5602).tay, (*Olc5602).imp, true, 2}, {"LDA", (*Olc5602).lda, (*Olc5602).imm, false, 2}, {"TAX", (*Olc5602).tax, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"LDY", (*Olc5602).ldy, (*Olc5602).abs, false, 4}, {"LDA", (*Olc5602).lda, (*Olc5602).abs, false, 4}, {"LDX", (*Olc5602).ldx, (*Olc5602).abs, false, 4}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 4},
		{"BCS", (*Olc5602).bcs, (*Olc5602).rel, false, 2}, {"LDA", (*Olc5602).lda, (*Olc5602).izy, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5}, {"LDY", (*Olc5602).ldy, (*Olc5602).zpx, false, 4}, {"LDA", (*Olc5602).lda, (*Olc5602).zpx, false, 4}, {"LDX", (*Olc5602).ldx, (*Olc5602).zpy, false, 4}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 4}, {"CLV", (*Olc5602).clv, (*Olc5602).imp, true, 2}, {"LDA", (*Olc5602).lda, (*Olc5602).aby, false, 4}, {"TSX", (*Olc5602).tsx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 4}, {"LDY", (*Olc5602).ldy, (*Olc5602).abx, false, 4}, {"LDA", (*Olc5602).lda, (*Olc5602).abx, false, 4}, {"LDX", (*Olc5602).ldx, (*Olc5602).aby, false, 4}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 4},
		{"CPY", (*Olc5602).cpy, (*Olc5602).imm, false, 2}, {"CMP", (*Olc5602).cmp, (*Olc5602).izx, false, 6}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"CPY", (*Olc5602).cpy, (*Olc5602).zp0, false, 3}, {"CMP", (*Olc5602).cmp, (*Olc5602).zp0, false, 3}, {"DEC", (*Olc5602).dec, (*Olc5602).zp0, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5}, {"INY", (*Olc5602).iny, (*Olc5602).imp, true, 2}, {"CMP", (*Olc5602).cmp, (*Olc5602).imm, false, 2}, {"DEX", (*Olc5602).dex, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"CPY", (*Olc5602).cpy, (*Olc5602).abs, false, 4}, {"CMP", (*Olc5602).cmp, (*Olc5602).abs, false, 4}, {"DEC", (*Olc5602).dec, (*Olc5602).abs, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6},
		{"BNE", (*Olc5602).bne, (*Olc5602).rel, false, 2}, {"CMP", (*Olc5602).cmp, (*Olc5602).izy, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"CMP", (*Olc5602).cmp, (*Olc5602).zpx, false, 4}, {"DEC", (*Olc5602).dec, (*Olc5602).zpx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6}, {"CLD", (*Olc5602).cld, (*Olc5602).imp, true, 2}, {"CMP", (*Olc5602).cmp, (*Olc5602).aby, false, 4}, {"NOP", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"CMP", (*Olc5602).cmp, (*Olc5602).abx, false, 4}, {"DEC", (*Olc5602).dec, (*Olc5602).abx, false, 7}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7},
		{"CPX", (*Olc5602).cpx, (*Olc5602).imm, false, 2}, {"SBC", (*Olc5602).sbc, (*Olc5602).izx, false, 6}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"CPX", (*Olc5602).cpx, (*Olc5602).zp0, false, 3}, {"SBC", (*Olc5602).sbc, (*Olc5602).zp0, false, 3}, {"INC", (*Olc5602).inc, (*Olc5602).zp0, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 5}, {"INX", (*Olc5602).inx, (*Olc5602).imp, true, 2}, {"SBC", (*Olc5602).sbc, (*Olc5602).imm, false, 2}, {"NOP", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).sbc, (*Olc5602).imp, true, 2}, {"CPX", (*Olc5602).cpx, (*Olc5602).abs, false, 4}, {"SBC", (*Olc5602).sbc, (*Olc5602).abs, false, 4}, {"INC", (*Olc5602).inc, (*Olc5602).abs, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6},
		{"BEQ", (*Olc5602).beq, (*Olc5602).rel, false, 2}, {"SBC", (*Olc5602).sbc, (*Olc5602).izy, false, 5}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 8}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"SBC", (*Olc5602).sbc, (*Olc5602).zpx, false, 4}, {"INC", (*Olc5602).inc, (*Olc5602).zpx, false, 6}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 6}, {"SED", (*Olc5602).sed, (*Olc5602).imp, true, 2}, {"SBC", (*Olc5602).sbc, (*Olc5602).aby, false, 4}, {"NOP", (*Olc5602).nop, (*Olc5602).imp, true, 2}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7}, {"???", (*Olc5602).nop, (*Olc5602).imp, true, 4}, {"SBC", (*Olc5602).sbc, (*Olc5602).abx, false, 4}, {"INC", (*Olc5602).inc, (*Olc5602).abx, false, 7}, {"???", (*Olc5602).xxx, (*Olc5602).imp, true, 7},
	}
}
