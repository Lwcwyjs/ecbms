#!/bin/bash

echo "========================================"
echo "  边缘计算盒子管理系统 - 启动脚本"
echo "========================================"

echo ""
echo "启动后端服务..."
cd backend
go mod download
nohup go run cmd/server/main.go > backend.log 2>&1 &
BACKEND_PID=$!
echo "后端服务已启动 (PID: $BACKEND_PID)"

echo ""
echo "启动前端服务..."
cd ../frontend
npm install
npm run dev &
FRONTEND_PID=$!
echo "前端服务已启动 (PID: $FRONTEND_PID)"

echo ""
echo "========================================"
echo "  服务启动完成！"
echo "  前端地址: http://localhost:3000"
echo "  后端地址: http://localhost:8080"
echo "  默认账号: admin / admin123"
echo "========================================"
echo ""
echo "按 Ctrl+C 停止所有服务"

trap "echo ''; echo '正在停止服务...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; echo '服务已停止'; exit" INT

wait
