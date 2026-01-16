@echo off
REM 检查是否安装 Npcap
chcp 65001 >nul
where wpcap.dll >nul 2>nul
if errorlevel 1 (
    where npcap.exe >nul 2>nul
    if errorlevel 1 (
        echo 正在下载并安装 Npcap...
        powershell -Command "Invoke-WebRequest -Uri https://npcap.com/dist/npcap-1.86.exe -OutFile npcap.exe"
        echo 下载完成，以管理员身份运行安装程序
    )
    start npcap.exe
    echo 安装完成后重新运行此程序
    pause
    exit
)

REM 运行程序
netguard.exe