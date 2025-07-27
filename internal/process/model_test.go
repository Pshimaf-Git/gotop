package process

import (
	"context"
	"fmt"
	"testing"

	"github.com/shirou/gopsutil/v3/process"
	"github.com/stretchr/testify/assert"
)

func Test_toMB(t *testing.T) {
	t.Run("base case", func(t *testing.T) {

	})
}

func TestFromProcess(t *testing.T) {
	t.Run("base case", func(t *testing.T) {
		proc := process.Process{
			Pid: 1000,
		}

		pInfo := FromProcess(context.Background(), &proc)
		assert.Equal(t, proc.Pid, pInfo.PID, "FromProcess(proc) = pInfo.Pid != proc.Pid")
	})
}

func TestProcessInfo_String(t *testing.T) {
	t.Run("base case", func(t *testing.T) {
		p := ProcessInfo{
			PID:        1,
			PPID:       12,
			Name:       UNKOWN_NAME,
			MemoryMB:   1000,
			CPUPercent: 1.21,
			ComandLine: "./main.exe",
			Username:   "root",
			Status:     "status",
		}

		s := fmt.Sprintf(
			"PID: %d, PPID: %d, Name %s, CPU: %.1f, Memory: %dMB, Cmd: %s, User: %s, Status: %s",
			p.PID, p.PPID, p.Name, p.CPUPercent, p.MemoryMB, p.ComandLine, p.Username, p.Status,
		)

		assert.Equal(t, s, p.String())
	})
}
