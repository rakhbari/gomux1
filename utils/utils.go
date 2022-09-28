package utils

import (
    "fmt"
    "bytes"
    "os"

    "github.com/rakhbari/gomux1/config"
)

func ProcessError(err error) {
    fmt.Printf("ERROR: %v", err)
    os.Exit(1)
}

func GetTlsCertBundleFile(cfg *config.Config) string {
    if len(cfg.Server.TlsCaPaths) == 0 {
        fmt.Println("---> No tlsCaPaths - returning cfg.Server.TlsCertPath ...")
        return cfg.Server.TlsCertPath
    }
    caCertPaths := []string{cfg.Server.TlsCertPath} // Initialize with value of the "leaf" cert
    caCertPaths = append(caCertPaths, cfg.Server.TlsCaPaths...) // Append the tlsCaPaths to it
    fmt.Printf("---> Processing caCertPaths: %+v\n", caCertPaths)

    // Loop through caCertPaths and concat all their content into bundleData
    var bundleData bytes.Buffer
    for _, filePath := range caCertPaths {
        fmt.Printf("------> Reading filePath: \"%+v\"\n", filePath)
        data, err := os.ReadFile(filePath)
        if err != nil {
            ProcessError(err)
        }

        bundleData.Write(data)
    }

    tlsCertBundle := "tlsCertBundle"
    err := os.WriteFile(tlsCertBundle, bundleData.Bytes(), 0644)
    if err != nil {
        ProcessError(err)
    }

    return tlsCertBundle
}

