@echo off
setlocal enabledelayedexpansion

REM Windows 本地开发测试启动脚本
REM 此脚本从 config.yaml 读取配置并设置环境变量，用于安全的本地测试
REM 使用方法: scripts\start-local.bat [--migrate|--help]

REM 脚本目录和项目根目录
set "SCRIPT_DIR=%~dp0"
set "PROJECT_ROOT=%SCRIPT_DIR%.."
set "CONFIG_FILE=%PROJECT_ROOT%\config.yaml"

REM 彩色输出函数 (Windows 不支持颜色，使用纯文本)
goto :main

:print_info
echo [信息] %~1
goto :eof

:print_success
echo [成功] %~1
goto :eof

:print_warning
echo [警告] %~1
goto :eof

:print_error
echo [错误] %~1
goto :eof

:show_help
echo Windows 本地开发测试启动脚本
echo.
echo 此脚本从 config.yaml 读取配置并设置环境变量，用于安全的本地测试，然后启动应用程序。
echo.
echo 使用方法:
echo     %~nx0 [选项]
echo.
echo 选项:
echo     --migrate       在启动前运行数据库迁移
echo     --env-only      只设置环境变量，不启动应用
echo     --help          显示此帮助信息
echo.
echo 设置的环境变量:
echo     SERVER_PORT         服务器监听端口
echo     JWT_KEY            JWT 签名密钥
echo     JWT_EXPIRATION     JWT 令牌过期时间
echo     REDIS_HOST         Redis 服务器地址
echo     REDIS_PASSWORD     Redis 密码
echo     POSTGRES_HOST      PostgreSQL 主机
echo     POSTGRES_USER      PostgreSQL 用户名
echo     POSTGRES_PASSWORD  PostgreSQL 密码
echo     POSTGRES_DBNAME    PostgreSQL 数据库名
echo     POSTGRES_PORT      PostgreSQL 端口
echo     POSTGRES_SSLMODE   PostgreSQL SSL 模式
echo     POSTGRES_TIMEZONE  PostgreSQL 时区
echo     ADMIN_PASSWORD     管理员密码
echo     PASSWORD_SALT      密码哈希盐值
echo.
echo 示例:
echo     scripts\start-local.bat              # 使用当前配置启动
echo     scripts\start-local.bat --migrate    # 运行迁移后启动
echo     scripts\start-local.bat --env-only   # 只设置环境变量
echo.
goto :eof

:load_config_to_env
call :print_info "正在从 %CONFIG_FILE% 加载配置"

if not exist "%CONFIG_FILE%" (
    call :print_error "找不到配置文件: %CONFIG_FILE%"
    call :print_error "请创建 config.yaml 文件或从项目根目录运行。"
    exit /b 1
)

REM 从 config.yaml 中提取值（简化的提取方式）
REM 注意：这是一个基本实现。在生产环境中，考虑使用合适的 YAML 解析器

REM 设置一些默认值并从配置中提取
set "SERVER_PORT=8080"
set "JWT_EXPIRATION=12h"
set "REDIS_HOST=127.0.0.1:6379"
set "POSTGRES_HOST=127.0.0.1"
set "POSTGRES_USER=postgres"
set "POSTGRES_DBNAME=server"
set "POSTGRES_PORT=5432"
set "POSTGRES_SSLMODE=disable"
set "POSTGRES_TIMEZONE=Asia/Shanghai"

REM 使用 findstr 从 config.yaml 提取值（基本提取）
for /f "tokens=2 delims=: " %%i in ('findstr /c:"port:" "%CONFIG_FILE%" 2^>nul ^| findstr /v /c:"#" ^| head -1') do set "SERVER_PORT=%%i"

for /f "tokens=2* delims=:" %%i in ('findstr /c:"key:" "%CONFIG_FILE%" 2^>nul ^| findstr /v /c:"#"') do (
    set "temp=%%j"
    set "temp=!temp: =!"
    set "temp=!temp:"=!"
    set "JWT_KEY=!temp!"
)

REM 提取过期时间
for /f "tokens=2* delims=:" %%i in ('findstr /c:"expiration:" "%CONFIG_FILE%" 2^>nul ^| findstr /v /c:"#"') do (
    set "temp=%%j"
    set "temp=!temp: =!"
    set "temp=!temp:"=!"
    set "JWT_EXPIRATION=!temp!"
)

REM 提取 Redis 主机
for /f "tokens=2* delims=:" %%i in ('findstr /A /c:"redis:" "%CONFIG_FILE%" 2^>nul') do (
    for /f "skip=1 tokens=2* delims=:" %%j in ('findstr /c:"host:" "%CONFIG_FILE%" 2^>nul') do (
        set "temp=%%k"
        set "temp=!temp: =!"
        set "temp=!temp:"=!"
        set "REDIS_HOST=!temp!"
        goto :redis_host_done
    )
)
:redis_host_done

REM 提取密码（Redis、PostgreSQL、Admin）
set "password_count=0"
for /f "tokens=2* delims=:" %%i in ('findstr /c:"password:" "%CONFIG_FILE%" 2^>nul ^| findstr /v /c:"#"') do (
    set "temp=%%j"
    set "temp=!temp: =!"
    set "temp=!temp:"=!"
    set /a password_count+=1
    
    if !password_count! equ 1 set "REDIS_PASSWORD=!temp!"
    if !password_count! equ 2 set "POSTGRES_PASSWORD=!temp!"
    if !password_count! equ 3 set "ADMIN_PASSWORD=!temp!"
)

REM 提取 PostgreSQL 用户名
for /f "tokens=2* delims=:" %%i in ('findstr /A /c:"postgres:" "%CONFIG_FILE%" 2^>nul') do (
    for /f "tokens=2* delims=:" %%j in ('findstr /c:"user:" "%CONFIG_FILE%" 2^>nul') do (
        set "temp=%%k"
        set "temp=!temp: =!"
        set "temp=!temp:"=!"
        set "POSTGRES_USER=!temp!"
        goto :pg_user_done
    )
)
:pg_user_done

REM 提取数据库名
for /f "tokens=2* delims=:" %%i in ('findstr /c:"dbname:" "%CONFIG_FILE%" 2^>nul') do (
    set "temp=%%j"
    set "temp=!temp: =!"
    set "temp=!temp:"=!"
    set "POSTGRES_DBNAME=!temp!"
)

REM 提取盐值
for /f "tokens=2* delims=:" %%i in ('findstr /c:"salt:" "%CONFIG_FILE%" 2^>nul ^| findstr /v /c:"#"') do (
    set "temp=%%j"
    set "temp=!temp: =!"
    set "temp=!temp:"=!"
    set "PASSWORD_SALT=!temp!"
)

call :print_success "从 config.yaml 设置环境变量成功"

REM 验证关键变量
set "missing_vars="
if "%JWT_KEY%"=="" set "missing_vars=%missing_vars% JWT_KEY"
if "%POSTGRES_PASSWORD%"=="" set "missing_vars=%missing_vars% POSTGRES_PASSWORD"
if "%ADMIN_PASSWORD%"=="" set "missing_vars=%missing_vars% ADMIN_PASSWORD"
if "%PASSWORD_SALT%"=="" set "missing_vars=%missing_vars% PASSWORD_SALT"

if not "%missing_vars%"=="" (
    call :print_error "缺少必需的配置值:"
    for %%i in (%missing_vars%) do call :print_error "  - %%i"
    exit /b 1
)

REM 显示加载的配置（敏感值加掩码）
call :print_info "已加载的配置:"
echo   SERVER_PORT=%SERVER_PORT%
echo   JWT_KEY=%JWT_KEY:~0,8%*** (已掩码)
echo   JWT_EXPIRATION=%JWT_EXPIRATION%
echo   REDIS_HOST=%REDIS_HOST%
if defined REDIS_PASSWORD (echo   REDIS_PASSWORD=***已掩码***) else (echo   REDIS_PASSWORD=^(未设置^))
echo   POSTGRES_HOST=%POSTGRES_HOST%
echo   POSTGRES_USER=%POSTGRES_USER%
echo   POSTGRES_PASSWORD=***已掩码***
echo   POSTGRES_DBNAME=%POSTGRES_DBNAME%
echo   POSTGRES_PORT=%POSTGRES_PORT%
echo   POSTGRES_SSLMODE=%POSTGRES_SSLMODE%
echo   POSTGRES_TIMEZONE=%POSTGRES_TIMEZONE%
echo   ADMIN_PASSWORD=***已掩码***
echo   PASSWORD_SALT=%PASSWORD_SALT:~0,8%*** (已掩码)

goto :eof

:run_migration
call :print_info "正在运行数据库迁移..."
cd /d "%PROJECT_ROOT%"

if not exist "main.go" (
    call :print_error "找不到 main.go 文件。请从项目根目录运行。"
    exit /b 1
)

go run main.go --migrate
call :print_success "数据库迁移完成"
goto :eof

:start_application
call :print_info "正在启动应用程序..."
cd /d "%PROJECT_ROOT%"

if not exist "main.go" (
    call :print_error "找不到 main.go 文件。请从项目根目录运行。"
    exit /b 1
)

call :print_success "正在使用环境变量启动服务器 (config.yaml 值将被忽略)"
call :print_info "按 Ctrl+C 停止服务器"
go run main.go
goto :eof

:main
set "run_migrate=false"
set "env_only=false"

REM 解析命令行参数
:parse_args
if "%~1"=="" goto :args_done
if "%~1"=="--migrate" (
    set "run_migrate=true"
    shift
    goto :parse_args
)
if "%~1"=="--env-only" (
    set "env_only=true"
    shift
    goto :parse_args
)
if "%~1"=="--help" (
    call :show_help
    exit /b 0
)
if "%~1"=="-h" (
    call :show_help
    exit /b 0
)
call :print_error "未知选项: %~1"
call :show_help
exit /b 1

:args_done
call :print_info "艺术设计专业版 Go 服务器 - 本地测试脚本 (Windows)"
call :print_warning "此脚本仅用于本地测试"
call :print_warning "请勿在生产环境中使用此脚本"
echo.

REM 加载配置
call :load_config_to_env
echo.

REM 如果需要，运行迁移
if "%run_migrate%"=="true" (
    call :run_migration
    echo.
)

REM 除非是仅环境变量模式，否则启动应用程序
if "%env_only%"=="false" (
    call :start_application
) else (
    call :print_success "环境变量已设置。现在可以运行您的应用程序:"
    call :print_info "cd %PROJECT_ROOT% && go run main.go"
)