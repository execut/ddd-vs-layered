package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/labels"
    "github.com/spf13/cobra"
)

func InitCreateLabelTemplate(ctx context.Context, service *labels.Service) {
    var createLabelTemplate = &cobra.Command{
        Use:   "labels-create-template",
        Short: "",
        RunE: func(_ *cobra.Command, _ []string) error {
            err := service.CreateLabelTemplate(ctx, labelTemplateID, name)
            if err != nil {
                return err
            }

            fmt.Println("1")

            return nil
        },
    }
    createLabelTemplate.PersistentFlags().StringVarP(&name, "manufacturer-organization-name", "m",
        "", "manufacturer-organization-name")
    rootCmd.AddCommand(createLabelTemplate)
}
