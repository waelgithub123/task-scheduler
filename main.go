package main

import (
	"bytes"
	"database/sql"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Task struct {
	ID             int
	Name           string
	Command        string
	Interval       int
}

var db *sql.DB

func main() {
	// Database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Create DSN and connect to MySQL
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true&loc=UTC"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Verify database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Successfully connected to database")

	// Start scheduler loop
	go startScheduler()

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down scheduler...")
}

func startScheduler() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tasks, err := fetchAndScheduleTasks()
			if err != nil {
				log.Println("Error fetching tasks:", err)
				continue
			}

			for _, task := range tasks {
				go executeTask(task)
			}
		}
	}
}

func fetchAndScheduleTasks() ([]Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Lock and retrieve due tasks
	rows, err := tx.Query(`
		SELECT id, name, command, interval_seconds 
		FROM tasks 
		WHERE next_run <= UTC_TIMESTAMP() AND status = 'enabled'
		FOR UPDATE
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Command, &task.Interval); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Update next_run for locked tasks
	for _, task := range tasks {
		_, err := tx.Exec(`
			UPDATE tasks 
			SET next_run = DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? SECOND) 
			WHERE id = ?
		`, task.Interval, task.ID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func executeTask(task Task) {
	startTime := time.Now().UTC()
	var output bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("sh", "-c", task.Command)
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	err := cmd.Run()
	endTime := time.Now().UTC()

	success := err == nil
	errorMessage := ""
	if err != nil {
		errorMessage = err.Error() + "\n" + stderr.String()
	}

	// Log task execution
	_, err = db.Exec(`
		INSERT INTO task_executions (
			task_id, start_time, end_time, success, output, error
		) VALUES (?, ?, ?, ?, ?, ?)
	`, task.ID, startTime, endTime, success, output.String(), errorMessage)
	
	if err != nil {
		log.Printf("Failed to log execution for task %d: %v", task.ID, err)
	}
}