package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type Task struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

const JSON_FILE = "todo.json"

func save(tasks []Task) {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		panic(err)
	}
	os.WriteFile(JSON_FILE, data, 0644)
}

func load() []Task {
	var tasks []Task
	jsonFile, err := os.Open(JSON_FILE)

	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}
		}
		panic(err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &tasks)
	return tasks
}

func menu() {
	fmt.Println("Use: todo [add | list | done]")
	fmt.Println("     todo add     <task_name>       		# Add a task")
	fmt.Println("     todo remove  <task_id>       		# Add a task")
	fmt.Println("     todo done    <task_id>         		# Mark a task as complete")
	fmt.Println("     todo list    [--completed | --pending]	# Add a task")
}

func getFilter(arg string) func(Task) bool {
	switch arg {
	case "--completed":
		return func(t Task) bool { return t.Completed }
	case "--pending":
		return func(t Task) bool { return !t.Completed }
	default:
		return func(t Task) bool { return true }
	}
}

func main() {
	args := os.Args
	logger := Logger{}

	if len(args) < 2 {
		menu()
		return
	}

	cmd := args[1]
	tasks := load()

	switch cmd {

	case "list":
		var filter func(Task) bool = func(t Task) bool { return true }

		if len(args) == 3 {
			filter = getFilter(args[2])
		}

		for i, task := range tasks {
			if !filter(task) {
				continue
			}

			status := " "
			if task.Completed {
				status = color.GreenString("✓")
			}

			fmt.Printf("%d. [%s] %s\n", i+1, status, task.Name)
		}

	case "add":
		if len(args) < 3 {
			logger.Error("Specify a name for the task")
			return
		}

		name := strings.Join(args[2:], " ")

		tasks = append(tasks, Task{Name: name, Completed: false})
		save(tasks)

		logger.SuccessF("Added task: %s\n", name)

	case "remove":
		if len(args) < 3 {
			logger.Error("Error! Specify a Task ID")
			return
		}

		id, err := strconv.Atoi(args[2])
		if err != nil {
			logger.Error("Error! Not a valid ID")
			return
		}

		if id < 1 || id > len(tasks) {
			logger.Error("Error! Not a valid ID")
			return
		}

		tasks = append(tasks[:id-1], tasks[id:]...)
		save(tasks)

		logger.SuccessF("Removed task n.%d", id)

	case "done":
		if len(args) != 3 {
			logger.Error("Error! Specify a Task ID")
			return
		}

		id, err := strconv.Atoi(args[2])
		if err != nil {
			logger.Error("Error! Not a valid ID")
			return
		}

		if id < 1 || id > len(tasks) {
			logger.Error("Error! Not a valid ID")
			return
		}

		if tasks[id-1].Completed {
			logger.Error("Error! That task is already completed.")
			return
		}

		tasks[id-1].Completed = true
		save(tasks)

		logger.SuccessF("Task n.%d completed!", id)

	default:
		menu()
	}

}
