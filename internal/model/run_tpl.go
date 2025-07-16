package model

// RunTpl CLE配置模板
var RunTpl = `#!/bin/bash

cd {{ .DestPath }}

# 赋予可执行权限
echo "[run.sh] chmod +x {{ .App }}"
chmod +x {{ .App }}

# kill 已经存在的进程
debugPid=$(lsof -t -i:{{ .ServerDebugPort }})
if [  -n  "$debugPid"  ];  then
    echo "[run.sh] kill -9 $debugPid"
    kill -9 $debugPid;
fi

echo "[run.sh] 启动 dlv..."
nohup dlv \
  --log \
  --log-output=debugger \
  --listen=:{{ .ServerDebugPort }} \
  --headless=true \
  --api-version=2 \
  --accept-multiclient \
  --check-go-version=false \
  exec {{ .App }} -- {{ .RunCmdArgs}} 2>&1 &
`
