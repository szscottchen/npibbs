#!/usr/bin/env node

const { spawn } = require('child_process');
const fs = require('fs');

// 创建日志文件
const logFile = 'vite-debug.log';
const logStream = fs.createWriteStream(logFile, { flags: 'a' });

console.log(`启动Vite调试模式，日志将输出到: ${logFile}`);
console.log('当前工作目录:', process.cwd());

// 启动Vite开发服务器
const viteProcess = spawn('npx', ['vite', '--config', './config/vite.config.dev.ts', '--debug'], {
  stdio: ['pipe', 'pipe', 'pipe'],
  shell: true,
  env: { ...process.env, DEBUG: 'vite:*' }
});

// 处理标准输出
viteProcess.stdout.on('data', (data) => {
  const output = data.toString();
  console.log(`[VITE STDOUT] ${output}`);
  logStream.write(`[STDOUT] ${new Date().toISOString()}: ${output}`);
});

// 处理标准错误
viteProcess.stderr.on('data', (data) => {
  const output = data.toString();
  console.error(`[VITE STDERR] ${output}`);
  logStream.write(`[STDERR] ${new Date().toISOString()}: ${output}`);
});

// 处理进程退出
viteProcess.on('close', (code, signal) => {
  console.log(`\nVite进程退出，退出码: ${code}, 信号: ${signal}`);
  logStream.write(`\n[EXIT] ${new Date().toISOString()}: 退出码 ${code}, 信号 ${signal}\n`);
  
  if (code === 1) {
    console.log('检测到ELIFECYCLE错误，请检查日志文件获取详细信息');
    console.log('常见原因：端口占用、内存不足、依赖冲突');
  }
  
  logStream.end();
});

// 捕获未处理的异常
process.on('uncaughtException', (error) => {
  console.error('未捕获的异常:', error);
  logStream.write(`[UNCAUGHT_EXCEPTION] ${new Date().toISOString()}: ${error.stack}\n`);
});

process.on('unhandledRejection', (reason, promise) => {
  console.error('未处理的Promise拒绝:', reason);
  logStream.write(`[UNHANDLED_REJECTION] ${new Date().toISOString()}: ${reason}\n`);
});

// 优雅退出处理
process.on('SIGINT', () => {
  console.log('\n收到SIGINT信号，正在关闭Vite进程...');
  viteProcess.kill('SIGINT');
});

process.on('SIGTERM', () => {
  console.log('\n收到SIGTERM信号，正在关闭Vite进程...');
  viteProcess.kill('SIGTERM');
});

console.log('Vite调试脚本已启动，按Ctrl+C退出...');