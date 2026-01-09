package main

import (
    "context"
    "fmt"

    "effective-architecture/steps/1-layered/labels"
    "github.com/spf13/cobra"
)

func main() {

    var (
        id   string
        name string
    )
    rootCmd := &cobra.Command{
        Use:   "",
        Short: "",
        Long:  ``,
    }

    ctx := context.Background()
    repository, err := labels.NewRepository(ctx)
    if err != nil {
        panic(err)
    }

    service := labels.NewService(repository)
    var createLabelTemplate = &cobra.Command{
        Use:   "labels-create-template",
        Short: "",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Println("1")
            return service.CreateLabelTemplate(ctx, id, name)
        },
    }

    rootCmd.PersistentFlags().StringVarP(&id, "id", "i", "", "id")
    createLabelTemplate.PersistentFlags().StringVarP(&name, "manufacturer-organization-name", "m", "", "manufacturer-organization-name")
    rootCmd.AddCommand(createLabelTemplate)
    err = rootCmd.Execute()
    if err != nil {
        panic(err)
    }
}
