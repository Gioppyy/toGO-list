package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Task struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}

const JSON_FILE = "todo.json"

var (
	app       = tview.NewApplication()
	pages     = tview.NewPages()
	editIndex = 0
	tasks     []Task
	menuList  = tview.NewList()
)

func save() {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		panic(err)
	}
	os.WriteFile(JSON_FILE, data, 0644)
}

func load() {
	data, err := os.ReadFile(JSON_FILE)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
		tasks = []Task{}
		return
	}
	json.Unmarshal(data, &tasks)
}

func showList() {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)

	list := tview.NewList().ShowSecondaryText(false)
	for i, task := range tasks {
		status := " "
		if task.Completed {
			status = "[green]✓[white]"
		}
		list.AddItem(
			fmt.Sprintf("%d. [%s] %s | %s", i+1, status, task.CreatedAt, task.Name),
			"", 0, nil,
		)
	}
	list.AddItem(" ", "", 0, nil)

	testPanel := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("[ENTER] Mark a task as completed\n[e | E] Edit a task \n[DEL] Delete a task\n [ESC] Back to Main menu")

	testPanel.SetBorder(true).
		SetTitle(" How to Use ").
		SetTitleAlign(tview.AlignLeft)

	flex.AddItem(list, 0, 1, true).
		AddItem(testPanel, 6, 1, false)

	flex.SetBorder(true).
		SetTitle(" toGO-List | List Task ").
		SetTitleAlign(tview.AlignCenter)

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		selected := list.GetCurrentItem()
		switch event.Key() {
		case tcell.KeyEnter:
			if selected >= 0 && selected < len(tasks) {
				tasks[selected].Completed = !tasks[selected].Completed
				save()
				showList()
			}
		case tcell.KeyRune:
			if event.Rune() == 'e' || event.Rune() == 'E' {
				editIndex = selected
				editTask()
				return nil
			}
		case tcell.KeyDelete:
			if selected >= 0 && selected < len(tasks) {
				tasks = append(tasks[:selected], tasks[selected+1:]...)
				save()
				showList()
			}
		case tcell.KeyESC:
			pages.SwitchToPage("menu")
		}
		return event
	})

	pages.AddAndSwitchToPage("List", flex, true)
}

func addTask() {
	flex := tview.NewFlex()

	inputField := tview.NewInputField().
		SetLabel("Name: ").
		SetFieldWidth(30)

	form := tview.NewForm().
		AddFormItem(inputField).
		AddButton("Add", func() {
			taskName := inputField.GetText()

			if taskName != "" {
				tasks = append(tasks, Task{Name: taskName, Completed: false, CreatedAt: time.Now().Format(time.DateTime)})
				save()
				pages.SwitchToPage("menu")
			}
		}).
		AddButton("Cancel", func() {
			pages.SwitchToPage("menu")
		})

	form.SetLabelColor(tcell.ColorWhite)

	flex.AddItem(form, 0, 1, true)
	flex.SetBorder(true).
		SetTitle(" toGO-List | Add Task ").
		SetTitleAlign(tview.AlignCenter)

	pages.AddAndSwitchToPage("Add", flex, true)
}

func editTask() {
	flex := tview.NewFlex()

	inputField := tview.NewInputField().
		SetLabel("Edit:").
		SetText(tasks[editIndex].Name)

	form := tview.NewForm().
		AddFormItem(inputField).
		AddButton("Save", func() {
			newName := strings.TrimSpace(inputField.GetText())
			if newName != "" {
				tasks[editIndex].Name = newName
				save()
				showList()
			}
		}).
		AddButton("Cancel", func() {
			pages.SwitchToPage("List")
		})

	form.SetLabelColor(tcell.ColorWhite)

	flex.AddItem(form, 0, 1, true)
	flex.SetBorder(true).
		SetTitle(" toGO-List | Edit Task ").
		SetTitleAlign(tview.AlignCenter)

	pages.AddAndSwitchToPage("edit", flex, true)
}

func menu() {
	menuList.Clear()
	menuList.AddItem("List", "Show all the task", 'l', func() { showList() })
	menuList.AddItem("Add", "Add a task", 'a', func() { addTask() })
	menuList.AddItem("Exit", "Quit", 'q', func() { app.Stop() })

	menuList.SetBorder(true).SetTitle("   toGO-List | Menu   ").SetTitleAlign(tview.AlignCenter)
	pages.AddAndSwitchToPage("menu", menuList, true)
}

func main() {
	load()
	menu()

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
