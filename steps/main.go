package main

import (
    "context"
    "fmt"

    labels2 "effective-architecture/steps/labels"
    "github.com/spf13/cobra"
)

func main() {
    var (
        labelTemplateID string
        name            string
    )

    rootCmd := &cobra.Command{
        Use:   "",
        Short: "",
        Long:  ``,
    }

    ctx := context.Background()

    repository, err := labels2.NewRepository(ctx)
    if err != nil {
        panic(err)
    }

    service := labels2.NewService(repository)

    var createLabelTemplate = &cobra.Command{
        Use:   "labels-create-template",
        Short: "",
        RunE: func(_ *cobra.Command, _ []string) error {
            err = service.CreateLabelTemplate(ctx, labelTemplateID, name)
            if err != nil {
                return err
            }

            fmt.Println("1")

            return nil
        },
    }

    rootCmd.PersistentFlags().StringVarP(&labelTemplateID, "id", "i", "", "id")
    createLabelTemplate.PersistentFlags().StringVarP(&name, "manufacturer-organization-name", "m",
        "", "manufacturer-organization-name")
    rootCmd.AddCommand(createLabelTemplate)

    err = rootCmd.Execute()
    if err != nil {
        panic(err)
    }
}
