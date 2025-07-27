package process

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/process"
)

func GetProcessInfoChan(ctx context.Context, interval time.Duration) <-chan ProcessInfo {
	info := make(chan ProcessInfo)

	go func() {
		defer close(info)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			processes, err := process.ProcessesWithContext(ctx)
			if err != nil {
				time.Sleep(interval)
				continue
			}

			for _, proc := range processes {
				select {
				case <-ctx.Done():
					return

				case info <- FromProcess(ctx, proc):
				}
			}

			time.Sleep(interval)
		}
	}()

	return info
}
