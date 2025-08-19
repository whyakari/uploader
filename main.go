package main

import (
    "bytes"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Uso: go run main.go <device>")
        os.Exit(1)
    }

    device := os.Args[1]
    outDir := filepath.Join("out", "target", "product", device)

    files := map[string]string{
        "Boot Image":        filepath.Join(outDir, "boot.img"),
        "Vendor Boot Image": filepath.Join(outDir, "vendor_boot.img"),
        "DTBO Image":        filepath.Join(outDir, "dtbo.img"),
    }

    zips, err := filepath.Glob(filepath.Join(outDir, "*.zip"))
    if err != nil {
        log.Fatal(err)
    }
    if len(zips) > 0 {
        files["ROM Zip"] = zips[0]
    }

    if _, err := os.Stat("./go-up"); os.IsNotExist(err) {
        fmt.Println("[INFO] go-up not found, downloading...")
        cmd := exec.Command("wget", "https://raw.githubusercontent.com/GustavoMends/go-up/master/go-up", "-O", "go-up")
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
            log.Fatal(err)
        }
        if err := exec.Command("chmod", "+x", "./go-up").Run(); err != nil {
            log.Fatal(err)
        }
    }

    fmt.Println("\n=== Upload files ===\n")

    links := make(map[string]string)

    for label, path := range files {
        if _, err := os.Stat(path); os.IsNotExist(err) {
            fmt.Printf("[ERROR] %s not found: %s\n", label, path)
            continue
        }

        fmt.Printf("[UPLOAD] %s -> %s\n", label, path)

        cmd := exec.Command("./go-up", path)
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = os.Stderr
        if err := cmd.Run(); err != nil {
            fmt.Printf("[ERRO] Failed to send %s: %v\n", label, err)
            continue
        }
        link := strings.TrimSpace(out.String())
        links[label] = link
    }

    fmt.Println("\n=== Generated Links ===\n")
    for label, link := range links {
        fmt.Printf("%s:\n%s\n\n", label, link)
    }
}
