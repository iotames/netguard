-- 创建流量监控记录表

CREATE TABLE IF NOT EXISTS ng_hook_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    remote_ip VARCHAR(45) NOT NULL,
    remote_port INTEGER,
    protocol VARCHAR(10),
    process_name VARCHAR(255),
    process_pid INTEGER,
    bytes_current_len BIGINT,
    inbound BOOLEAN,
    ip_country VARCHAR(100),
    ip_city VARCHAR(100),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引以优化查询性能
CREATE INDEX IF NOT EXISTS idx_logs_remote_ip ON ng_hook_logs(remote_ip);
CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON ng_hook_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_logs_process ON ng_hook_logs(process_name);