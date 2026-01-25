package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/contract"
    "github.com/spf13/cobra"
)

func InitLabelsDeleteTemplate(ctx context.Context, app contract.IApplication) error {
    var deleteLabelTemplateCmd = &cobra.Command{
        Use:   "labels-delete-template",
        Short: "",
        Long:  ``,
        RunE: func(_ *cobra.Command, _ []string) error {
            err := app.Delete(ctx, labelTemplateID)
            if err != nil {
                return err
            }

            fmt.Println("1")

            return nil
        },
    }

    deleteLabelTemplateCmd.PersistentFlags().StringVarP(&organizationName, "manufacturer-organization-name",
        "m", "", "manufacturer-organization-name")

    rootCmd.AddCommand(deleteLabelTemplateCmd)

    return nil
}
