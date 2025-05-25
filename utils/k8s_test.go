package utils

import (
    "log"
    "os"
    "path"
    "testing"
)

func TestGetSvcAcctToken(t *testing.T) {
    home, err := os.UserHomeDir()
    if err != nil {
        panic(err)
    }

    svcAcctToken, err := GetSvcAcctToken(path.Join(home, ".kube/config"), "app1", "user1")
    if err != nil {
        log.Fatalf("Error from GetSvcAcctToken: %v", err)
    }
    log.Printf("svcAcctToken: %s", *svcAcctToken)
}
