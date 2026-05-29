@echo off
chcp 65001 >nul
echo ========================================
echo   边缘计算盒子管理系统 - 启动脚本
echo ========================================
echo.

echo 启动后端服务...
cd backend
start "ECBMS Backend" cmd /k "go mod download && go run cmd/server/main.go"
cd ..

timeout /t 3 /nobreak >nul

echo.
echo 启动前端服务...
cd frontend
start "ECBMS Frontend" cmd /k "npm install && npm run dev"
cd ..

echo.
echo ========================================
echo   服务启动完成！
echo   前端地址: http://localhost:3000
echo   后端地址: http://localhost:8080
echo   默认账号: admin / admin123
echo ========================================
echo.
pause
