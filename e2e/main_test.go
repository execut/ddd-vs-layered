package e2e_test

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "testing"
)

//var binaryName = "scenario_of_transaction"

var currentPath = ""

func TestMain(m *testing.M) {
    err := os.Chdir("..")
    if err != nil {
        fmt.Printf("could not change dir: %v", err)
        os.Exit(1)
    }

    currentPath, err = os.Getwd()
    if err != nil {
        fmt.Printf("could not get current dir: %v", err)
    }

    os.Exit(m.Run())
}

func runBinary(args []string) (string, error) {
    cmd := exec.Command(filepath.Join(currentPath, "bin/main"), args...)

    output, err := cmd.CombinedOutput()
    stdOutString := strings.Trim(string(output), "\n")
    if err != nil {
        return stdOutString, err
    }

    return stdOutString, nil
}
