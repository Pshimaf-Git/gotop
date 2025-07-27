package ui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Pshimaf-Git/gotop/internal/config"
	"github.com/Pshimaf-Git/gotop/internal/process"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var colors = tcell.ColorNames

type ProcessTable struct {
	*tview.Table
	sortByColomn int
	sortOrder    string
}

const (
	DESC = "desc"
	ASC  = "asc"
)

const (
	arrowUp   = "▲"
	arrowDown = "▼"
)

var headers = []string{
	"PID(1)", "PPID(2)", "Name(3)", "CPU (%)(4)", "Memory (MB)(5)", "User(6)", "Comand line(7)", "Status(8)",
}

const (
	pid = iota
	ppid
	name
	cpu
	memory
	user
	comandLine
	status
)

func NewProcessTable(cfg config.Config) ProcessTable {
	table := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false).SetFixed(1, 0)

	if cfg.ShowColumnSeparator {
		table.SetSeparator(tview.Borders.Vertical)
	}

	if color, ok := colors[cfg.ColumnBordersColor]; ok {
		table.SetBordersColor(color)
	} else {
		table.SetBordersColor(tcell.ColorWhite)
	}

	pt := ProcessTable{
		Table:        table,
		sortByColomn: cpu,
		sortOrder:    DESC,
	}

	pt.setupHeaders()

	return pt
}

func (pt *ProcessTable) SetSortColumn(col int) {
	if pt.sortByColomn == col {
		pt.SwapOrder()
	} else {
		pt.sortByColomn = col
		pt.sortOrder = DESC
	}
}

func (pt *ProcessTable) SwapOrder() {
	if pt.sortOrder == ASC {
		pt.sortOrder = DESC
	} else {
		pt.sortOrder = ASC
	}
}

func (pt *ProcessTable) GetTable() *tview.Table {
	return pt.Table
}

func (pt *ProcessTable) setupHeaders() {
	for col, header := range headers {
		cell := tview.NewTableCell(header).
			SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorYellow).SetSelectable(false)

		if col == pt.sortByColomn {
			arrow := addSpace(arrowUp)
			if pt.sortOrder == DESC {
				arrow = addSpace(arrowDown)
			}

			cell.SetText(header + arrow)
		}

		pt.Table.SetCell(0, col, cell)
	}

}

func (pt *ProcessTable) UpdateData(processes map[int32]process.ProcessInfo) {
	processesList := make([]process.ProcessInfo, 0, len(processes))

	for _, pInfo := range processes {
		processesList = append(processesList, pInfo)
	}

	sort.Slice(processesList, func(i, j int) bool {
		a := processesList[i]
		b := processesList[j]

		var less bool
		switch pt.sortByColomn {
		case pid:
			less = a.PID < b.PID
		case ppid:
			less = a.PPID < b.PPID
		case name:
			less = strings.ToLower(a.Name) < strings.ToLower(b.Name)
		case cpu:
			less = a.CPUPercent < b.CPUPercent
		case memory:
			less = a.MemoryMB < b.MemoryMB
		case user:
			less = strings.ToLower(a.Username) < strings.ToLower(b.Username)
		case comandLine:
			less = strings.ToLower(a.ComandLine) < strings.ToLower(b.ComandLine)
		case status:
			less = strings.ToLower(a.Status) < strings.ToLower(b.Status)
		default:
			less = a.PID < b.PID
		}

		if pt.sortOrder == DESC {
			return !less
		}

		return less
	})

	for i := pt.Table.GetRowCount() - 1; i > 0; i-- {
		pt.Table.RemoveRow(i)
	}

	row := 1
	for _, pInfo := range processesList {
		pt.SetCellPID(row, pInfo).
			SetCellPPID(row, pInfo).
			SetCellName(row, pInfo).
			SetCellCPUPercent(row, pInfo).
			SetCellMemoryMB(row, pInfo).
			SetCellUsername(row, pInfo).
			SetCellComandLine(row, pInfo).
			SetCellStatus(row, pInfo)

		row++
	}

	pt.setupHeaders()
}

func (pt *ProcessTable) SetCellPID(row int, pInfo process.ProcessInfo) *ProcessTable {
	pt.Table.SetCell(row, pid, tview.NewTableCell(strconv.Itoa(int(pInfo.PID))).SetAlign(tview.AlignCenter))
	return pt
}

func (pt *ProcessTable) SetCellPPID(row int, pInfo process.ProcessInfo) *ProcessTable {
	pt.Table.SetCell(row, ppid, tview.NewTableCell(strconv.Itoa(int(pInfo.PPID))).SetAlign(tview.AlignCenter))
	return pt
}

func (pt *ProcessTable) SetCellName(row int, pInfo process.ProcessInfo) *ProcessTable {
	pt.Table.SetCell(row, name, tview.NewTableCell(pInfo.Name).SetAlign(tview.AlignCenter))
	return pt
}

func (pt *ProcessTable) SetCellStatus(row int, pInfo process.ProcessInfo) *ProcessTable {
	pt.Table.SetCell(row, status, tview.NewTableCell(pInfo.Status).SetAlign(tview.AlignCenter))
	return pt
}

func (pt *ProcessTable) SetCellUsername(row int, pInfo process.ProcessInfo) *ProcessTable {
	pt.Table.SetCell(row, user, tview.NewTableCell(pInfo.Username).SetAlign(tview.AlignCenter))
	return pt
}

func (pt *ProcessTable) SetCellComandLine(row int, pInfo process.ProcessInfo) *ProcessTable {
	pt.Table.SetCell(row, comandLine, tview.NewTableCell(pInfo.ComandLine).SetAlign(tview.AlignCenter))
	return pt

}

func (pt *ProcessTable) SetCellCPUPercent(row int, pInfo process.ProcessInfo) *ProcessTable {
	pt.Table.SetCell(row, cpu, tview.NewTableCell(fmt.Sprintf("%.1f", pInfo.CPUPercent)).SetAlign(tview.AlignCenter))
	return pt
}

func (pt *ProcessTable) SetCellMemoryMB(row int, pInfo process.ProcessInfo) *ProcessTable {
	pt.Table.SetCell(row, memory, tview.NewTableCell(fmt.Sprintf("%d", pInfo.MemoryMB)).SetAlign(tview.AlignCenter))
	return pt
}

func addSpace(s string) string {
	return " " + s
}
