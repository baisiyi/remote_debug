package model

// RunTpl CLE配置模板
var RunTpl = `#!/bin/bash

cd {{ .DestPath }}

echo "[run.sh] chmod +x {{ .App }}"
chmod +x {{ .App }}

debugPid=$(lsof -t -i:{{ .ServerDebugPort }})
if [ -n "$debugPid" ]; then
    echo "[run.sh] kill $debugPid"
    kill $debugPid
fi

echo "[run.sh] 启动 dlv..."
/go/bin/dlv exec {{ .App }}\
  --log \
  --log-output=debugger \
  --listen=:{{ .ServerDebugPort }} \
  --headless=true \
  --api-version=2 \
  --accept-multiclient \
  --check-go-version=false \
  --continue \
  -- {{ .RunCmdArgs}} > /logs/app.log 2>&1 &

dlv_pid=$!
echo "[run.sh] dlv started with PID $dlv_pid"

echo "[run.sh] done."`
