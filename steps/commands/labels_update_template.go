package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/labels"
    "github.com/spf13/cobra"
)

func InitUpdateLabelTemplate(ctx context.Context, service *labels.Service) {
    var updateLabelTemplate = &cobra.Command{
        Use:   "labels-update-template",
        Short: "",
        RunE: func(_ *cobra.Command, _ []string) error {
            err := service.UpdateLabelTemplate(ctx, labelTemplateID, name)
            if err != nil {
                return err
            }

            fmt.Println("1")

            return nil
        },
    }
    updateLabelTemplate.PersistentFlags().StringVarP(&name, "manufacturer-organization-name", "m",
        "", "manufacturer-organization-name")
    rootCmd.AddCommand(updateLabelTemplate)
}
