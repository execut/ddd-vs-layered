package main

import "effective-architecture/steps/commands"

func main() {
    err := commands.Execute()
    if err != nil {
        panic(err)
    }
}
