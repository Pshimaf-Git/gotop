package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	PID        int32
	PPID       int32
	Name       string
	CPUPercent float64
	MemoryMB   uint64
	ComandLine string
	Username   string
	Status     string
}

const (
	UNKOWN_PPID = -1
	UNKOWN_NAME = "N/A"
)

const (
	STATUS_SEP = ", "
)

func FromProcess(ctx context.Context, process *process.Process) ProcessInfo {
	name, err := process.NameWithContext(ctx)
	if err != nil {
		name = UNKOWN_NAME
	}

	cpuPercent, err := process.CPUPercentWithContext(ctx)
	if err != nil {
		cpuPercent = 0.0
	}

	var rssMB uint64
	memInfo, err := process.MemoryInfoWithContext(ctx)
	if err == nil && memInfo != nil {
		rssMB = toMB(memInfo.RSS)
	}

	ppid, err := process.PpidWithContext(ctx)
	if err != nil {
		ppid = UNKOWN_PPID
	}

	cmdLine, err := process.CmdlineWithContext(ctx)
	if err != nil {
		cmdLine = UNKOWN_NAME
	}
	cmdLine = strings.TrimSpace(cmdLine)

	username, err := process.UsernameWithContext(ctx)
	if err != nil {
		username = UNKOWN_NAME
	}

	statuses, err := process.StatusWithContext(ctx)
	if err != nil {
		statuses = []string{UNKOWN_NAME}
	}
	status := strings.Join(statuses, STATUS_SEP)

	return ProcessInfo{
		PID:        process.Pid,
		PPID:       ppid,
		Name:       name,
		CPUPercent: cpuPercent,
		MemoryMB:   rssMB,
		ComandLine: cmdLine,
		Username:   username,
		Status:     status,
	}
}

func (pi *ProcessInfo) String() string {
	return fmt.Sprintf(
		"PID: %d, PPID: %d, Name %s, CPU: %.1f, Memory: %dMB, Cmd: %s, User: %s, Status: %s",
		pi.PID, pi.PPID, pi.Name, pi.CPUPercent, pi.MemoryMB, pi.ComandLine, pi.Username, pi.Status,
	)
}

func toMB(rss uint64) uint64 {
	return rss / (1024 * 1024)
}
