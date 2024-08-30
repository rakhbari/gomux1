package utils

import (
    "bytes"
    "encoding/json"
    "errors"
    "log"
    "os"

    "github.com/rakhbari/gomux1/config"
)

func ProcessError(err error) {
    log.Printf("ERROR: %v", err)
    os.Exit(1)
}

func GetTlsCertFile(cfg *config.Config) (tlsCertFile *string) {
    if len(cfg.Server.TlsCaPaths) == 0 {
        log.Println("---> No tlsCaPaths provided - Checking TlsCertPath ...")
        if _, err := os.Stat(cfg.Server.TlsCertPath); errors.Is(err, os.ErrNotExist) {
            log.Printf("!!!!!!> ERROR: File \"%+v\" doesn't exist! Returning ...\n", cfg.Server.TlsCertPath)
            return nil
        }
        return &cfg.Server.TlsCertPath
    }
    caCertPaths := []string{cfg.Server.TlsCertPath}             // Initialize with value of the "leaf" cert
    caCertPaths = append(caCertPaths, cfg.Server.TlsCaPaths...) // Append the tlsCaPaths to it
    log.Printf("---> Processing caCertPaths: %+v\n", caCertPaths)

    // Loop through caCertPaths and concat all their content into bundleData
    var bundleData bytes.Buffer
    for _, filePath := range caCertPaths {
        log.Printf("------> Reading filePath: \"%+v\"\n", filePath)
        if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
            log.Printf("!!!!!!> ERROR: File \"%+v\" doesn't exist! Returning ...\n", filePath)
            return nil
        }
        data, err := os.ReadFile(filePath)
        if err != nil {
            log.Printf("!!!!!!> ERROR: %v", err)
            return nil
        }
        bundleData.Write(data)
    }

    tlsCertBundle := cfg.Server.TempDir + "/tlsCertBundle"
    err := os.WriteFile(tlsCertBundle, bundleData.Bytes(), 0644)
    if err != nil {
        log.Printf("!!!!!!> ERROR: %v", err)
        return nil
    }

    return &tlsCertBundle
}

func LoadVersion(version *Version) {
    versionFile := "version.json"
    f, err := os.Open(versionFile)
    if err != nil {
        log.Printf("---> Version file %s not found. Returning ...", versionFile)
        return
    }
    defer f.Close()

    err = json.NewDecoder(f).Decode(version)
    if err != nil {
        log.Printf("!!!> ERROR: Unable to parse version file %s. Error: %v", versionFile, err)
    }
}
