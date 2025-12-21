# Task Tracker CLI

Task Tracker is a professional, lightweight Command Line Interface (CLI) task tracker written in Go. It helps you stay productive by organizing your tasks with priorities, tracking progress, and measuring completion time—all stored locally in a standard, cross-platform configuration directory.

## Features

- **Zero-Config**: Works out of the box with zero setup required.
- **Priority Management**: Assign `low`, `medium`, or `high` priorities.
- **Unique Short IDs**: Uses Sqids for user-friendly 5-character IDs.
- **Auto-Persistence**: Saves tasks to a JSON file automatically.
- **Time Tracking**: Records creation and completion timestamps, calculating duration taken.
- **Smart Pathing**: Uses standard OS configuration directories to keep your home folder clean.

## Installation

Install the binary directly to your Go bin folder:
```bash
go install github.com/pramudya3/task-cli/cmd/task@latest
```

---

**Note for macOS/Linux users**: Ensure your Go bin directory is in your PATH. Add this to your `~/.zshrc` or `~/.bashrc`:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Usage

```bash
task add "Complete project documentation" -p high  # Add task
task list                                          # List pending
task list --all                                    # List everything
task complete <id>                                 # Mark as done
task remove <id>                                   # Delete task
task clean                                         # Clear all tasks
task version                                       # Show version info
```

## Configuration

Task Tracker follows OS standards for storage:
- **macOS**: `~/Library/Application Support/task/tasks.json`
- **Linux**: `~/.config/task/tasks.json`
- **Windows**: `%AppData%\task\tasks.json`

## Uninstall

To completely remove Task Tracker from your system:

### 1. Remove binary
```bash
rm $(which task)
```

### 2. Clean Data & Config
Delete the storage directory:
- **macOS**: `rm -rf ~/Library/Application\ Support/task`
- **Linux**: `rm -rf ~/.config/task`
- **Windows**: `rmdir /s /q %AppData%\task`

## License
MIT License.
