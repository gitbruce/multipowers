#!/bin/bash

while true; do
    echo "[$(date)] Trying git push..."

    # 10 秒超时执行 git push
    if timeout 20 git push; then
        echo "[$(date)] Push succeeded."
        exit 0
    else
        echo "[$(date)] Push failed or timed out. Retry in 120 seconds..."
    fi

    sleep 120
done