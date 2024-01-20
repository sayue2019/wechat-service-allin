#!/usr/bin/env bash

# Variables
wechat_version="3.9.2.23" # wechat版本
injector_name="wechat-bot" # 需要注入的wechat服务，可选：comwechatrobot(3.7.0.30)，wechat-bot(3.9.2.23), wxhelper(3.9.2.23)
injector_box_dir="docker_buiding/injector-box"
wechat_box_dir="${injector_box_dir}/wechat-box"

# Functions
setup_directory() {
    mkdir docker_buiding || true
}

update_git_repo() {
    local repo_url=$1
    local clone_dir=$2

    if [ ! -d "$clone_dir" ]; then
        git clone "$repo_url" "$clone_dir"
    else
        (cd "$clone_dir" && git pull)
    fi
}

download_file() {
    local file_url=$1
    local destination=$2

    if [ ! -f "$destination" ]; then
        wget -O "$destination" "$file_url"
    fi
}

run_vnc_auth() {
    cp vnc/index.html "${injector_box_dir}/root/"
    cp vnc/inj-entrypoint.sh "${injector_box_dir}/root/"
    cp vnc/x11vnc.conf "${wechat_box_dir}/root/wechat-etc/supervisord.d"
    cp vnc/websockify.conf "${wechat_box_dir}/root/wechat-etc/supervisord.d"
}

run_http_frward() {
    cp forward/http-forward "${injector_box_dir}/root/bin/"
    cp forward/forward-monitor "${injector_box_dir}/root/bin/"
    cp forward/inj-entrypoint.sh "${injector_box_dir}/root/"
}

injector_wechat_bot() {
    cp bin_deps/wechat-bot/funtool_wx=3.9.2.23.exe "${injector_box_dir}/root/bin/"
    cp bin_deps/wechat-bot/inject-dll "${injector_box_dir}/root/bin"
    cp bin_deps/wechat-bot/inject-monitor "${injector_box_dir}/root/bin"
}

injector_wxhelper() {
    cp bin_deps/wxhelper/wxhelper.dll "${injector_box_dir}/root/drive_c/injector/auto.dll"
}

injector_comwechatrobot() {
    cp bin_deps/comwechatrobot/inject-dll "${injector_box_dir}/root/bin"
    cp bin_deps/comwechatrobot/inject-monitor "${injector_box_dir}/root/bin"
    cp bin_deps/comwechatrobot/wxDriver.exe "${injector_box_dir}/root/drive_c/injector/"
    cp bin_deps/comwechatrobot/http/SWeChatRobot.dll "${injector_box_dir}/root/drive_c/injector/"
    cp bin_deps/comwechatrobot/http/wxDriver.dll "${injector_box_dir}/root/drive_c/injector/"
    cp bin_deps/comwechatrobot/http/wxDriver.lib "${injector_box_dir}/root/drive_c/injector/"
    cp bin_deps/comwechatrobot/http/wxDriver64.dll "${injector_box_dir}/root/drive_c/injector/"
    cp bin_deps/comwechatrobot/http/wxDriver64.lib "${injector_box_dir}/root/drive_c/injector/"
}

# Function to select the injector based on injname variable
injector_select() {
    local injector_name=$1

    case "$injector_name" in
        "wechat_bot")
            injector_wechat_bot
            ;;
        "wxhelper")
            injector_wxhelper
            ;;
        "comwechatrobot")
            injector_comwechatrobot
            ;;
        *)
            echo "Invalid injector name: $injector_name"
            ;;
    esac
}

build_docker_image() {
    (cd "$injector_box_dir" && sudo docker build -t "sayue/wechat-service-${injector_name}:latest" . --progress=plain)
}

# Script execution
setup_directory
update_git_repo "https://github.com/sayue2019/injector-box" "$injector_box_dir"
update_git_repo "https://github.com/sayue2019/wechat-box" "$wechat_box_dir"
download_file "https://github.com/tom-snow/wechat-windows-versions/releases/download/v${wechat_version}/WeChatSetup-${wechat_version}.exe" "${wechat_box_dir}/root/WeChatSetup.exe"
download_file "https://github.com/sayue2019/go-http-forward/releases/download/win86/http_forwarder.exe" "${injector_box_dir}/root/bin/http_forwarder.exe"
run_vnc_auth
run_http_frward
injector_select "$injector_name"
build_docker_image
