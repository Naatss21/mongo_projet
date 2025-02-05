CREATE DATABASE monitoring_system;

USE monitoring_system;

CREATE TABLE system_metrics (
    id INT AUTO_INCREMENT PRIMARY KEY,
    hostname VARCHAR(255),
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cpu_usage FLOAT,
    memory_usage FLOAT,
    disk_usage FLOAT,
    network_sent BIGINT,
    network_received BIGINT
);
