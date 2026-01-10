package e2e_test

import (
    "bytes"
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
    //cmd.Env = append(os.Environ(), "GOCOVERDIR=.coverdata")
    var (
        stdOut = bytes.NewBufferString("")
        stdErr = bytes.NewBufferString("")
    )

    cmd.Stdout = stdOut
    cmd.Stderr = stdErr

    err := cmd.Run()
    if err != nil {
        return "", err
    }

    stdOutString := strings.Trim(stdOut.String(), "\n")

    return stdOutString, nil
}
