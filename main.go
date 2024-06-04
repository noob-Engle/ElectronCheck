package main

import (
    "bufio"
    "bytes"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "sort"
    "strings"
)

func getAppName(appPath string) string {
    return filepath.Base(appPath)
}

func getElectronVersion(filename string) string {
    cmd := exec.Command("sh", "-c", fmt.Sprintf("strings '%s' | grep 'Chrome/' | grep -i Electron | grep -v '%%s' | sort -u | cut -f 3 -d '/'", filename))
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return ""
    }
    return strings.TrimSpace(out.String())
}

func generateGitHubLink(version string) string {
    return fmt.Sprintf("https://github.com/electron/electron/releases/tag/v%s", version)
}

func formatRow(appName, electronVersion, filename string) string {
    return fmt.Sprintf("%-30s %-20s %s", appName, electronVersion, filename)
}

func findElectronApps() {
    cmd := exec.Command("mdfind", "kind:app")
    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    scanner := bufio.NewScanner(&out)
    var apps []string
    for scanner.Scan() {
        apps = append(apps, scanner.Text())
    }

    sort.Strings(apps)

    fmt.Println(strings.Repeat("_", 100))
    fmt.Println(formatRow("App Name", "Electron Version", "File Name"))
    fmt.Println(strings.Repeat("=", 100))

    for _, app := range apps {
        filename := filepath.Join(app, "Contents/Frameworks/Electron Framework.framework/Electron Framework")
        if _, err := os.Stat(filename); err == nil {
            appName := getAppName(app)
            electronVersion := getElectronVersion(filename)
            fmt.Println(formatRow(appName, electronVersion, filename))
        }
    }

    fmt.Println("\n")
    fmt.Println(strings.Repeat("=", 60))
    fmt.Println(fmt.Sprintf("%-30s %s", "App Name", "GitHub Link"))
    fmt.Println(strings.Repeat("=", 60))

    for _, app := range apps {
        filename := filepath.Join(app, "Contents/Frameworks/Electron Framework.framework/Electron Framework")
        if _, err := os.Stat(filename); err == nil {
            appName := getAppName(app)
            electronVersion := getElectronVersion(filename)
            githubLink := "Not Found"
            if electronVersion != "" {
                githubLink = generateGitHubLink(electronVersion)
            }
            fmt.Printf("%-30s %s\n", appName, githubLink)
        }
    }
}

func main() {
    findElectronApps()
}
