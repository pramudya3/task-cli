package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pramudya3/task-cli/helper"
)

type Task struct {
	Id          string        `json:"id"`
	Description string        `json:"description"`
	Priority    string        `json:"priority"`
	Completed   bool          `json:"completed"`
	CreatedAt   time.Time     `json:"created_at"`
	CompletedAt time.Time     `json:"completed_at"`
	TookTime    time.Duration `json:"took_time"`
}

type TaskTracker struct {
	Tasks    []Task `json:"tasks"`
	FilePath string `json:"file_path"`
}

func LoadTaskTracker(filePath string) (*TaskTracker, error) {
	tt := &TaskTracker{
		Tasks:    make([]Task, 0),
		FilePath: filePath,
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return tt, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, tt); err != nil {
		return nil, err
	}

	return tt, nil
}

func (tt *TaskTracker) Save() error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(tt.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(tt, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal tasks: %w", err)
	}
	if err := os.WriteFile(tt.FilePath, data, 0644); err != nil {
		return fmt.Errorf("could not write task file: %w", err)
	}
	return nil
}

func (tt *TaskTracker) Add(description, priority string) (*Task, error) {
	for _, t := range tt.Tasks {
		if t.Description == description {
			return nil, fmt.Errorf("task already exists: %s", description)
		}
	}

	task := Task{
		Id:          helper.GenerateId(description),
		Description: description,
		Priority:    priority,
		CreatedAt:   time.Now(),
	}

	tt.Tasks = append(tt.Tasks, task)
	return &task, nil
}

func (tt *TaskTracker) FindTask(id string) *Task {
	for i := range tt.Tasks {
		if tt.Tasks[i].Id == id {
			return &tt.Tasks[i]
		}
	}
	return nil
}

func (tt *TaskTracker) Complete(id string) error {
	task := tt.FindTask(id)
	if task == nil {
		return fmt.Errorf("task not found: %s", id)
	}
	task.Completed = true
	task.CompletedAt = time.Now()
	task.TookTime = time.Since(task.CreatedAt)
	return nil
}

func (tt *TaskTracker) Remove(id string) (*Task, error) {
	for i, t := range tt.Tasks {
		if t.Id == id {
			tt.Tasks = append(tt.Tasks[:i], tt.Tasks[i+1:]...)
			return &t, nil
		}
	}
	return nil, fmt.Errorf("task not found: %s", id)
}

func (tt *TaskTracker) CleanUp() {
	tt.Tasks = make([]Task, 0)
}

func IsValidPriority(p string) bool {
	p = strings.ToLower(p)
	return p == "low" || p == "medium" || p == "high"
}
