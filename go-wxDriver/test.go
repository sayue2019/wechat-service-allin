package main

import (
    "bytes"
    "log"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "fmt"
    "strings"
)

type GitHubRelease []struct {
    TagName string `json:"tag_name"`
}


func main() {
    port := 5999 // 假设端口号为 19088
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

    return // 手工断点，程序到这里将终止

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
    url := fmt.Sprintf("http://127.0.0.1:%d/api/?type=%d&password=wechat5999", port, api)

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
