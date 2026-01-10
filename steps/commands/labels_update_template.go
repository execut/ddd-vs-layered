package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/application"
    "github.com/spf13/cobra"
)

func InitLabelsUpdateTemplate(ctx context.Context, app *application.Application) error {
    var updateLabelTemplateCmd = &cobra.Command{
        Use:   "labels-update-template",
        Short: "",
        Long:  ``,
        Run: func(_ *cobra.Command, _ []string) {
            err := app.UpdateLabelTemplate(ctx, labelTemplateID, name)
            if err != nil {
                panic(err)
            }

            fmt.Println("1")
        },
    }

    updateLabelTemplateCmd.PersistentFlags().StringVarP(&name, "manufacturer-organization-name",
        "m", "", "manufacturer-organization-name")
    rootCmd.AddCommand(updateLabelTemplateCmd)

    return nil
}
