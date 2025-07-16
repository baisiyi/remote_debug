package model

// RunTpl CLE配置模板
var RunTpl = `#!/bin/bash

cd {{ .DestPath }}

echo "[run.sh] chmod +x {{ .App }}"
chmod +x {{ .App }}

# 优雅地杀死占用调试端口的进程
debugPid=$(lsof -t -i:{{ .ServerDebugPort }})
if [ -n "$debugPid" ]; then
    echo "[run.sh] kill $debugPid"
    kill $debugPid
    sleep 2
    # 如果还在，强制杀死
    if kill -0 $debugPid 2>/dev/null; then
        echo "[run.sh] kill -9 $debugPid"
        kill -9 $debugPid
    fi
fi

# 优雅地杀死服务进程
appPid=$(pgrep -f "{{ .App }}")
if [ -n "$appPid" ]; then
    echo "[run.sh] kill $appPid"
    kill $appPid
    sleep 2
    if kill -0 $appPid 2>/dev/null; then
        echo "[run.sh] kill -9 $appPid"
        kill -9 $appPid
    fi
fi

sleep 2

echo "[run.sh] 启动 dlv..."
/go/bin/dlv \
  --log \
  --log-output=debugger \
  --listen=:{{ .ServerDebugPort }} \
  --headless=true \
  --api-version=2 \
  --accept-multiclient \
  --check-go-version=false \
  exec {{ .App }} -- {{ .RunCmdArgs}} > /logs/app.log 2>&1 &

dlv_pid=$!
echo "[run.sh] dlv started with PID $dlv_pid"

echo "[run.sh] done."`
