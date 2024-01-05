# WeChat Service In Docker

微信作为一个服务运行在 Docker 容器中。本项目旨在提供一个简易方式来构建和运行一个微信服务，它允许通过 Docker 容器注入不同版本的微信以及相关的第三方服务。

## 功能

- 支持多版本微信运行环境
- 允许注入不同的第三方服务
- 容易配置和部署

## 目前支持第三方服务


- [https://github.com/cixingguangming55555/wechat-bot](https://github.com/cixingguangming55555/wechat-bot)
wechat-bot 更新慢，非注入方式，安全性高，但功能偏少，最新版本支持3.9.2.23
- [https://github.com/ttttupup/wxhelper](https://github.com/ttttupup/wxhelper)
wxhelper 持续更新，目前已经支持最新版微信
- [https://github.com/ljc545w/ComWeChatRobot](https://github.com/ljc545w/ComWeChatRobot)
ComWeChatRobot 功能最强，但是版本过旧已经很久没更新了，最高版本：3.7.0.30

## 快速开始

1. **设置脚本权限**  
   允许执行 `build-injector-box.sh` 脚本。
```sh
chmod +x build-injector-box.sh
```
2. **构建注入器和微信服务**  
运行脚本以构建和设置所需的 Docker 容器。
```sh
   ./build-injector-box.sh
```

3. **启动服务**  
使用 Docker Compose 启动服务。
```sh
docker-compose up
```

## 配置选项

- **微信版本**  
在 `build-injector-box.sh` 脚本中，可以自定义微信的版本和第三方插件版本。

- **第三方服务版本**  
支持不同版本的第三方服务。详细信息和配置选项在脚本文件中有说明。

- **相关依赖**  
请自行下载 `bin_deps` 目录下相应版本的 dll 文件。

## 脚本说明

### 主脚本：`build-injector-box.sh`

此脚本是项目的核心，用于设置目录、更新 Git 仓库、下载文件、运行 VNC 授权、运行 HTTP 转发、选择注入器和构建 Docker 镜像。

#### 功能说明：

- `setup_directory`: 创建必要的目录。
- `update_git_repo`: 克隆或更新 Git 仓库。
- `download_file`: 下载必要的文件。
- `run_vnc_auth`: 复制 VNC 相关文件。
- `run_http_frward`: 复制 HTTP 转发相关文件。
- `injector_select`: 根据所选注入器进行配置。
- `build_docker_image`: 构建 Docker 镜像。

### 注意事项

- 确保 Docker 环境已正确安装和配置。
- 脚本中的部分步骤可能需要网络连接来下载文件或更新仓库。
- 根据您的系统配置，可能需要对脚本赋予执行权限。
- 在运行 `build-injector-box.sh` 脚本之前，请确保所有必需的依赖文件已经放在 `bin_deps` 目录中
- 根据您的使用需求选择正确的 docker-compose.yaml 文件。
- 如果您需要修改 VNC 或 HTTP 转发的设置，可以编辑 `vnc` 和 `forward` 目录下的相关文件。
- 只开启了http服务，不同插件的websocket请自行配置。
- 相关dll仅在win x86环境下测试正常运行(64位环境请自行编译)


## 目录结构和说明

本项目包含多个关键目录，每个目录都承担特定的功能和作用。以下是项目的主要目录结构及其说明：

### `bin_deps`
此目录包含了所有必要的二进制依赖文件，如特定版本的微信机器人、DLL 文件等。根据不同的注入器选择，您可能需要从此目录中获取相应的文件。

### `forward`
包含 HTTP 转发相关的脚本和可执行文件。这些文件用于设置和管理微信服务与外部通信的转发机制。

### `vnc`
存放 VNC 服务相关的文件，如 HTML 页面、入口脚本和配置文件。VNC 服务用于远程访问和控制 Docker 容器中运行的微信实例。

### `docker-compose`
此目录包含基于三种不同 DLL 钩子的 Docker Compose 配置文件。根据您选择的注入器类型（wechat-bot、wxhelper 或 ComWeChatRobot），您可以使用相应的 docker-compose.yaml 文件来启动和管理容器。

#### 使用示例 docker-compose.yaml

```yaml
version: "3.3"

services:
    wechat-service:
        image: "sayue/wechat-service-comwechatrobot:latest"
        restart: unless-stopped
        container_name: "wechat-service-comwechatrobot"
        environment:
            #vnc访问密码
            VNC_PASSWORD: "vncpass"
            TARGET_AUTO_RESTART: "yes"
            INJMON_LOG_FILE: "/dev/stdout"
            # dll注入状态判断接口
            INJ_CONDITION: "[ \"`sudo netstat -tunlp | grep 19088`\" != '' ] && exit 0 ; sleep 5 ; curl 'http://127.0.0.1:8680/hi' 2>/dev/null | grep -P 'code.:0'"
            HOOK_PROC_NAME: "WeChat"
            TARGET_CMD: "wechat-start"
            HOOK_DLL: "auto.dll"
            #HTTP转发地址设置
            FORWARD_URL: "http://127.0.0.1:19088"
            #http转发端口设置
            LISTEN_PORT: "5999"
            #http访问密码设置
            ACCESS_PASSWORD: "wechat5999"
            #HTTP转发状态判断接口
            FOR_CONDITION: "[ \"`sudo netstat -tunlp | grep 5999`\" != '' ] && exit 0"
            #optional INJMON_LOG_FILE: "/dev/null"
            #optional TARGET_LOG_FILE: "/dev/stdout"
        ports:
            - "8088:8080" # noVNC
            #- "19088:19088" # websocket server
            - "7777:5999" # forward server
            #- "5900:5900" # vnc server
        volumes:
              - "~/bread/.wechat/WeChat Files/:/home/app/WeChat Files/"
              - "~/bread/.wechat/Applcation Data/:/home/app/.wine/drive_c/users/user/Application Data/"
              - "~/bread/external:/home/app/external"
        extra_hosts:
            - "dldir1.qq.com:127.0.0.1"
        tty: true

```


## 引用和参考

- [https://github.com/cixingguangming55555/wechat-bot](https://github.com/cixingguangming55555/wechat-bot)
- [https://github.com/ttttupup/wxhelper](https://github.com/ttttupup/wxhelper)
- [https://github.com/ljc545w/ComWeChatRobot](https://github.com/ljc545w/ComWeChatRobot)
- [go-http-forward](https://github.com/sayue2019/go-http-forward)
- [wechat-box](https://github.com/sayue2019/wechat-box)
- [injector-box](https://github.com/sayue2019/injector-box)
- 参考项目: [ChisBread/wechat-service](https://github.com/ChisBread/wechat-service)

## 许可

请在使用本项目时遵守相应的许可协议。

## 贡献

欢迎通过 Pull Request 或 Issue 来贡献和提出建议。