package cpu

import (
	"fmt"
	"strings"
)

func (c *CPU) printStep(instruction uint16) {

	c.logger.Debug(fmt.Sprintf("--- Step: PC [0x%04X] | Instruction: [0x%04X] ---", c.pc.GetValue()-2, instruction))

	regs := c.registerFile.GetRegisters()
	c.logger.Debug("Registers:")
	for i := range 4 {
		c.logger.Debug(fmt.Sprintf("  R%d: 0x%04X (%5d)  |  R%d: 0x%04X (%5d)",
			i, regs[i], int16(regs[i]),
			i+4, regs[i+4], int16(regs[i+4])))
	}
	c.logger.Debug(fmt.Sprintf("Condition Codes: [ %s ]",
		c.ccString()))
	c.logger.Debug(fmt.Sprint(strings.Repeat("-", 45)))
}

func (c *CPU) ccString() string {
	res := []rune("---")
	if c.conditionalCodes.N {
		res[0] = 'N'
	}
	if c.conditionalCodes.Z {
		res[1] = 'Z'
	}
	if c.conditionalCodes.P {
		res[2] = 'P'
	}
	return string(res)
}

func (c *CPU) printRegisterFile() {
	regs := c.registerFile.GetRegisters()

	c.logger.Info("Registers:")
	for i := range 4 {
		c.logger.Info(fmt.Sprintf("  R%d: 0x%04X (%5d)  |  R%d: 0x%04X (%5d)",
			i, regs[i], int16(regs[i]),
			i+4, regs[i+4], int16(regs[i+4])))
	}

}
