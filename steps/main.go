package main

import (
    "effective-architecture/steps/commands"
    "golang.org/x/text/cases"
    "golang.org/x/text/language"
)

func main() {
    err := commands.Execute()
    if err != nil {
        panic(cases.Title(language.Russian, cases.Compact).String(err.Error()))
    }
}
