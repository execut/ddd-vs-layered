package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/application"
    "github.com/spf13/cobra"
)

func InitLabelsCreateTemplate(ctx context.Context, app *application.Application) error {
    var createLabelTemplateCmd = &cobra.Command{
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

    createLabelTemplateCmd.PersistentFlags().StringVarP(&name, "manufacturer-organization-name",
        "m", "", "manufacturer-organization-name")
    rootCmd.AddCommand(createLabelTemplateCmd)

    return nil
}
