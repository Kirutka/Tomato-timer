package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TimerState struct {
	workDuration    time.Duration
	breakDuration   time.Duration
	currentDuration time.Duration
	remaining       time.Duration
	isRunning       bool
	isWorkPhase     bool
	pomodoros       int
}

func main() {
	app := tview.NewApplication()
	
	state := &TimerState{
		workDuration:    25 * time.Minute,
		breakDuration:   5 * time.Minute,
		currentDuration: 25 * time.Minute,
		remaining:       25 * time.Minute,
		isRunning:       false,
		isWorkPhase:     true,
		pomodoros:       0,
	}

	timerText := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	statusText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	updateTimer := func() {
		for {
			time.Sleep(100 * time.Millisecond)
			
			app.QueueUpdateDraw(func() {
				if state.isRunning {
					state.remaining -= 100 * time.Millisecond
					
					if state.remaining <= 0 {
						state.isRunning = false
						if state.isWorkPhase {
							state.pomodoros++
							state.isWorkPhase = false
							state.currentDuration = state.breakDuration
							state.remaining = state.breakDuration
							statusText.SetText("[red]⏰ Перерыв! Время отдохнуть!")
						} else {
							state.isWorkPhase = true
							state.currentDuration = state.workDuration
							state.remaining = state.workDuration
							statusText.SetText("[green]🍅 Работаем! Фокус!")
						}
					}
				}

				minutes := int(state.remaining.Minutes())
				seconds := int(state.remaining.Seconds()) % 60
				
				color := "green"
				if !state.isWorkPhase {
					color = "blue"
				}
				if state.remaining < time.Minute {
					color = "red" 
				}

				timerText.SetText(fmt.Sprintf("[%s]%02d:%02d[-]", color, minutes, seconds))
				
				phase := "РАБОТА"
				if !state.isWorkPhase {
					phase = "ОТДЫХ"
				}
				status := fmt.Sprintf("🍅 [yellow]%d[-] | [white]%s[-] | ", state.pomodoros, phase)
				if state.isRunning {
					status += "[green]▶ В процессе[-]"
				} else {
					status += "[red]⏸ Остановлено[-]"
				}
				statusText.SetText(status)
			})
		}
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			state.isRunning = !state.isRunning
		case tcell.KeyRune:
			switch event.Rune() {
			case ' ':
				state.isRunning = !state.isRunning
			case 'r', 'R':
				state.remaining = state.currentDuration
				state.isRunning = false
			case 'w', 'W':
				state.isWorkPhase = true
				state.currentDuration = state.workDuration
				state.remaining = state.workDuration
				state.isRunning = false
			case 'b', 'B':
				state.isWorkPhase = false
				state.currentDuration = state.breakDuration
				state.remaining = state.breakDuration
				state.isRunning = false
			case 'q', 'Q':
				app.Stop()
			}
		}
		return event
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 1, 1, false).
		AddItem(timerText, 3, 1, false).
		AddItem(tview.NewBox(), 1, 1, false).
		AddItem(statusText, 3, 1, false).
		AddItem(tview.NewBox(), 1, 1, false)

	mainFlex := tview.NewFlex().
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(flex, 0, 1, false).
		AddItem(tview.NewBox(), 0, 1, false)

	timerText.SetTextAlign(tview.AlignCenter).SetText("[green]25:00[-]")
	statusText.SetText("🍅 [yellow]0[-] | [white]РАБОТА[-] | [red]⏸ Остановлено[-]")

	go updateTimer()

	helpText := tview.NewTextView().
		SetDynamicColors(true).
		SetText("\n[white]Управление:[-]\n" +
			"[green]Пробел/Enter[-] - Старт/Пауза\n" +
			"[yellow]R[-] - Сброс таймера\n" +
			"[blue]W[-] - Переключить на работу\n" +
			"[cyan]B[-] - Переключить на перерыв\n" +
			"[red]Q[-] - Выход")

	finalFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, false).
		AddItem(helpText, 0, 1, false)

	if err := app.SetRoot(finalFlex, true).Run(); err != nil {
		panic(err)
	}
}