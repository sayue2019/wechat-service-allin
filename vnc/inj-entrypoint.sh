#!/usr/bin/env bash
sudo rm /tmp/.X0-lock
#为vnc添加访问权限
sudo cp /index.html /usr/share/novnc/
mkdir -p /home/app/.vnc
x11vnc -storepasswd $VNC_PASSWORD /home/app/.vnc/passwd

TARGET_AUTO_RESTART=${TARGET_AUTO_RESTART:-no}
TARGET_LOG_FILE=${TARGET_LOG_FILE:-/dev/null}
function run-target() {
    while :
    do
        $TARGET_CMD >${TARGET_LOG_FILE} 2>&1
        case ${TARGET_AUTO_RESTART} in
        false|no|n|0)
            exit 0
            ;;
        esac
    done
}

/entrypoint.sh &
sleep 5
inject-monitor &
run-target &
wait