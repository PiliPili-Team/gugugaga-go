module.exports = {
  apps: [
    {
      name: "gd-webhook",
      script: "./gd-webhook-server",
      cwd: "./",
      instances: 1,
      autorestart: true,
      watch: false,
      max_memory_restart: "512M",
      env: {
        NODE_ENV: "production"
      },
      // 日志时间格式
      log_date_format: "YYYY-MM-DD HH:mm:ss",
      // 错误日志路径 (可选)
      error_file: "./userdata/logs/pm2-error.log",
      // 输出日志路径 (可选)
      out_file: "./userdata/logs/pm2-out.log",
      // 合并日志
      merge_logs: true
    }
  ]
};
