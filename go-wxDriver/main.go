package main

import (
    "bytes"
    "log"
    "os"
    "os/signal"
    "runtime"
    "syscall"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "strings"
    "fmt"

    ps "github.com/mitchellh/go-ps"
)

type GitHubRelease []struct {
    TagName string `json:"tag_name"`
}


func main() {

    port := 19088 // 假设端口号为 19088

    log.Println("Starting the application...")

    // 1. 加载 DLL
    driver, err := loadDriver()
    if err != nil {
        log.Fatalf("Failed to load driver: %v", err)
    }
    defer syscall.FreeLibrary(driver)

    // 2. 获取微信 PID
    pidList, err := getWechatPIDList()
    if err != nil {
        log.Fatalf("Failed to get WeChat PID list: %v", err)
    }
    if len(pidList) == 0 {
        log.Fatal("No WeChat processes found")
    }

    log.Printf("Found WeChat process with PID: %d", pidList[0])

    // 3. 启动监听
    err = startListen(driver, uintptr(pidList[0]), port) 
    if err != nil {
        log.Fatalf("Failed to start listen: %v", err)
    }

    log.Println("Listening started on port %d", port)

    // 在此处调用postWechatHttpApi
    latestVersion, err := getLatestWechatVersion()
    if err != nil {
        log.Fatalf("Failed to get latest WeChat version: %v", err)
    }
    log.Printf("Latest WeChat version: %s", latestVersion)

    apiResponse, err := postWechatHttpApi(35, port, map[string]string{"version": latestVersion})
    if err != nil {
        log.Fatalf("Failed to set WeChat version: %v", err)
    }

    log.Printf("Response from API: %v", apiResponse)

    // 4. 保持程序运行
    keepRunning()
}

func loadDriver() (syscall.Handle, error) {
    var dllName string
    if runtime.GOARCH == "amd64" {
        dllName = "wxDriver64.dll"
    } else {
        dllName = "wxDriver.dll"
    }

    log.Printf("Loading driver: %s", dllName)

    driver, err := syscall.LoadLibrary(dllName)
    if err != nil {
        return 0, err
    }

    return driver, nil
}

func getWechatPIDList() ([]int, error) {
    var pidList []int
    processList, err := ps.Processes()
    if err != nil {
        return nil, err
    }

    for _, process := range processList {
        if process.Executable() == "WeChat.exe" {
            pidList = append(pidList, process.Pid())
        }
    }

    return pidList, nil
}

func startListen(driver syscall.Handle, pid uintptr, port int) error {
    log.Printf("Starting to listen on PID %d, port %d", pid, port)

    startListen, err := syscall.GetProcAddress(driver, "start_listen")
    if err != nil {
        return err
    }

    _, _, callErr := syscall.Syscall(uintptr(startListen), 2, pid, uintptr(port), 0)
    if callErr != 0 {
        return syscall.Errno(callErr)
    }

    return nil
}

func keepRunning() {
    log.Println("Service is now running. Press Ctrl+C to exit.")

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c

    log.Println("Shutting down...")
    // 在这里添加任何必要的清理操作
}

func getLatestWechatVersion() (string, error) {
    url := "https://api.github.com/repos/tom-snow/wechat-windows-versions/releases?per_page=1"
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var releases GitHubRelease
    err = json.Unmarshal(body, &releases)
    if err != nil {
        return "", err
    }

    if len(releases) > 0 {
        // 移除 "v" 前缀并返回版本号
        return strings.TrimPrefix(releases[0].TagName, "v"), nil
    }

    return "", fmt.Errorf("no releases found")
}

func postWechatHttpApi(api int, port int, data map[string]string) (map[string]interface{}, error) {
    url := fmt.Sprintf("http://127.0.0.1:%d/api/?type=%d", port, api)

    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    // 打印响应内容
    //log.Println("API Response:", string(body))

    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, err
    }

    return result, nil
}
