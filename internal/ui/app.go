package ui

import (
	"context"
	"time"

	"github.com/Pshimaf-Git/gotop/internal/config"
	"github.com/Pshimaf-Git/gotop/internal/process"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	app *tview.Application
}

func NewApp() App {
	return App{
		app: tview.NewApplication(),
	}
}

func (a *App) Run(ctx context.Context, processInfo <-chan process.ProcessInfo, cfg config.Config) error {
	processTable := NewProcessTable(cfg)

	a.app.SetRoot(processTable.GetTable(), true)

	a.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune {
			col := int(event.Rune() - '1')
			if col >= 0 && col < len(headers) {
				processTable.SetSortColumn(col)
			}
		}
		return event
	})

	go func() {
		allProcesses := make(map[int32]process.ProcessInfo)
		updateTicket := time.NewTicker(1000 * time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				return

			case pInfo, ok := <-processInfo:
				if !ok {
					return
				}
				allProcesses[pInfo.PID] = pInfo

			case <-updateTicket.C:
				a.app.QueueUpdateDraw(func() {
					processTable.UpdateData(allProcesses)
				})
			}
		}
	}()

	return a.app.Run()
}

func (a *App) Stop() {
	a.app.Stop()
}
