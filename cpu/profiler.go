package cpu

import (
	"fmt"
	"log/slog"
	"time"
)

// =============================================
//            CPU PROFILER FINAL REPORT
// =============================================
// Program Status:  HALTED (TRAP 0x25)
// Total Runtime:   2.45ms (Simulated)
// ---------------------------------------------
// EXECUTION METRICS:
//   Total Instructions:   18
//   Total Clock Cycles:   82
//   Avg. CPI:             4.56 cycles/inst
//   Avg. IPC:             0.22 inst/cycle

// MEMORY STATS:
//   Total Memory Accesses: 24
//   Memory Reads:          21 (18 Fetch + 3 Data)
//   Memory Writes:         3
//   Effective Latency:     320ns (avg)

// CONTROL FLOW:
//   Branches Taken:        4
//   Branches Not Taken:    1
//   Branch Penalty Cycles: 12
// ---------------------------------------------

type Profiler struct {
	StartTime time.Time
	EndTime   time.Time

	InstructionCount int
	TotalCycles      int
	MemReads         int
	MemWrites        int

	BranchesTaken    int
	BranchesNotTaken int
}

func (p *Profiler) Report(logger slog.Logger, state string) {
	duration := p.EndTime.Sub(p.StartTime)

	var cpi, intensity float64
	if p.InstructionCount > 0 {
		cpi = float64(p.TotalCycles) / float64(p.InstructionCount)
		intensity = float64(p.MemReads+p.MemWrites) / float64(p.InstructionCount)
	}

	report := fmt.Sprintf(`
=============================================
          CPU PROFILER FINAL REPORT         
=============================================
Status:    %s
Runtime:   %v (Simulated)
---------------------------------------------
EXECUTION:
  Instructions:    %-10d
  Total Cycles:    %-10d
  Avg CPI:         %-10.2f
  
MEMORY:
  Reads:           %-10d
  Writes:          %-10d
  Total Accesses:  %-10d
  Intensity:       %.2f ops/inst

BRANCHING:
  Taken:           %-10d
  Not Taken:       %-10d
---------------------------------------------
`,
		state, duration, p.InstructionCount, p.TotalCycles, cpi,
		p.MemReads, p.MemWrites, p.MemReads+p.MemWrites, intensity,
		p.BranchesTaken, p.BranchesNotTaken)

	fmt.Print(report)
}
