package main

import (
    "context"
    "fmt"

    "effective-architecture/steps/application"
    "effective-architecture/steps/infrastructure"
    "github.com/spf13/cobra"
)

var (
    rootCmd = &cobra.Command{
        Use:   "",
        Short: "",
        Long:  ``,
    }
)

func main() {
    repository, err := infrastructure.NewLabelTemplateRepository()
    if err != nil {
        panic(err)
    }

    app, err := application.NewApplication(repository)
    if err != nil {
        panic(err)
    }

    ctx := context.Background()

    var (
        labelTemplateID        string
        name                   string
        createLabelTemplateCmd = &cobra.Command{
            Use:   "labels-create-template",
            Short: "",
            Long:  ``,
            Run: func(_ *cobra.Command, _ []string) {
                err := app.CreateLabelTemplate(ctx, labelTemplateID, name)
                if err != nil {
                    panic(err)
                }

                fmt.Println("1")
            },
        }
        getLabelTemplateCmd = &cobra.Command{
            Use:   "labels-get-template",
            Short: "",
            Long:  ``,
            Run: func(_ *cobra.Command, _ []string) {
                result, err := app.GetLabelTemplate(ctx, labelTemplateID)
                if err != nil {
                    panic(err)
                }

                fmt.Println(result)
            },
        }
    )

    createLabelTemplateCmd.PersistentFlags().StringVarP(&name, "manufacturer-organization-name",
        "m", "", "manufacturer-organization-name")

    getLabelTemplateCmd.PersistentFlags().StringVarP(&name, "labels-get-template",
        "g", "", "labels-get-template")

    rootCmd.PersistentFlags().StringVarP(&labelTemplateID, "id", "i", "", "id")
    rootCmd.AddCommand(createLabelTemplateCmd)
    rootCmd.AddCommand(getLabelTemplateCmd)

    err = rootCmd.Execute()
    if err != nil {
        panic(err)
    }
}
