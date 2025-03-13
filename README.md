# Task Scheduler

A backend task scheduler built in Go, using MySQL for task storage and Docker for deployment. This system automates recurring tasks with customizable intervals and conditions, and includes failure handling and logging mechanisms for improved reliability and troubleshooting.

---

## Features

- **Recurring Task Scheduling**: Schedule tasks to run at customizable intervals.
- **Concurrent Execution**: Tasks are executed concurrently using goroutines.
- **Failure Handling**: Detailed error logging and task execution tracking.
- **Atomic Scheduling**: Ensures tasks are not missed or duplicated using MySQL transactions.
- **Docker Deployment**: Easy deployment using Docker and Docker Compose.
- **Logging**: Detailed logs for task executions, including output and errors.
- **Customizable Commands**: Execute any shell command as a task.

---

## Prerequisites

- Docker and Docker Compose installed on your system.

---

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/task-scheduler.git
cd task-scheduler