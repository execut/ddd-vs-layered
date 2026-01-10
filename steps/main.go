package main

import (
    "fmt"
    "os"

    "effective-architecture/steps/commands"
)

func main() {
    err := commands.Execute()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
