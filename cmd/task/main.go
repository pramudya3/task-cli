package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pramudya3/task-cli/domain"
	"github.com/pramudya3/task-cli/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version     = "dev"
	TaskFile    string
	Priority    string
	ShowAll     bool
	taskTracker *domain.TaskTracker
)

var RootCmd = &cobra.Command{
	Use:     "task",
	Short:   "A personal task tracker",
	Version: Version,
	Long: `task is a CLI task tracker that helps you organize your work and personal tasks.
Store tasks locally with priorities, mark them complete, and keep track of your productivity.`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		setupConfig(cmd)
		tt, err := domain.LoadTaskTracker(getTaskFile())
		if err != nil {
			return err
		}
		taskTracker = tt
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if taskTracker != nil {
			return taskTracker.Save()
		}
		return nil
	},
}

var addCmd = &cobra.Command{
	Use:   "add [description]",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !domain.IsValidPriority(Priority) {
			return fmt.Errorf("invalid priority %s, must be low, medium, or high", Priority)
		}

		description := strings.Join(args, " ")
		task, err := taskTracker.Add(description, Priority)
		if err != nil {
			return err
		}
		fmt.Printf("Added task %s: %s [%s]\n", task.Id, task.Description, task.Priority)
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		PrintTasks(taskTracker.Tasks, ShowAll)
	},
}

var removeCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove a task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		task, err := taskTracker.Remove(args[0])
		if err != nil {
			return err
		}
		fmt.Printf("Removed task %s: %s\n", task.Id, task.Description)
		return nil
	},
}

var completeCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := taskTracker.Complete(args[0]); err != nil {
			return err
		}
		fmt.Printf("Task %s marked as completed.\n", args[0])
		return nil
	},
}

var cleanUpCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean all tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Are you sure you want to clean all tasks? (y/n): ")
		var response string
		fmt.Scanln(&response)

		response = strings.ToLower(strings.TrimSpace(response))
		if response != "y" && response != "yes" {
			fmt.Println("Cleanup cancelled.")
			return nil
		}
		taskTracker.CleanUp()

		fmt.Println("All tasks cleaned up.")
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of task",
	Run: func(cmd *cobra.Command, args []string) {
		v := Version

		if v == "dev" {
			if info, ok := debug.ReadBuildInfo(); ok {
				if info.Main.Version != "" && info.Main.Version != "(devel)" {
					v = info.Main.Version
				}
			}
		}

		fmt.Printf("task version: %s\n", v)
	},
}

func PrintTasks(tasks []domain.Task, showAll bool) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	header := fmt.Sprintf("%-8s %-10s %-50s %-10s %-22s %-22s %-12s", "ID", "STATUS", "DESCRIPTION", "PRIORITY", "CREATED", "COMPLETED", "TOOK TIME")
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))

	for _, t := range tasks {
		if !showAll && t.Completed {
			continue
		}

		status := "PENDING"
		completedAt := "-"
		tookTime := "-"
		if t.Completed {
			status = "DONE"
			completedAt = t.CompletedAt.Format("2006-01-02 15:04:05")
			tookTime = t.TookTime.Truncate(time.Second).String()
		}

		fmt.Printf("%-8s %-10s %-50s %-10s %-22s %-22s %-12s\n",
			t.Id,
			status,
			helper.Truncate(t.Description, 50),
			strings.ToUpper(t.Priority),
			t.CreatedAt.Format("2006-01-02 15:04:05"),
			completedAt,
			tookTime,
		)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&TaskFile, "file", "", "Path to storage file")
	RootCmd.PersistentFlags().StringVarP(&Priority, "priority", "p", "medium", "Task priority (low, medium, high)")
	listCmd.Flags().BoolVarP(&ShowAll, "all", "a", false, "Show completed tasks")

	RootCmd.AddCommand(addCmd, listCmd, removeCmd, completeCmd, cleanUpCmd, versionCmd)
}

func getTaskFile() string {
	if TaskFile != "" {
		return TaskFile
	}
	return viper.GetString("file")
}

func setupConfig(cmd *cobra.Command) {
	home, _ := os.UserHomeDir()
	configDir, _ := os.UserConfigDir()

	viper.SetConfigName("task")
	viper.SetConfigType("json")
	viper.SetDefault("priority", "medium")

	// Look for task.json in current directory or system config directory
	viper.AddConfigPath(".")
	if configDir != "" {
		viper.AddConfigPath(filepath.Join(configDir, "task"))
	}

	viper.BindPFlag("file", cmd.PersistentFlags().Lookup("file"))

	// Default storage path resolution
	var storagePath string
	if configDir != "" {
		storagePath = filepath.Join(configDir, "task", "tasks.json")
	} else if home != "" {
		storagePath = filepath.Join(home, ".task", "tasks.json")
	}

	if storagePath != "" {
		viper.SetDefault("file", storagePath)
	}

	_ = viper.ReadInConfig()
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
