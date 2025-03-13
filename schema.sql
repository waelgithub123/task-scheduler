CREATE DATABASE IF NOT EXISTS scheduler;
USE scheduler;

CREATE TABLE IF NOT EXISTS tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    command TEXT NOT NULL,
    interval_seconds INT NOT NULL,
    next_run DATETIME NOT NULL,
    status ENUM('enabled', 'disabled') DEFAULT 'enabled'
);

CREATE TABLE IF NOT EXISTS task_executions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    task_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    success BOOLEAN,
    output TEXT,
    error TEXT,
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);