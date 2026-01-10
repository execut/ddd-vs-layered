package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/labels"
    "github.com/spf13/cobra"
)

func InitDeleteLabelTemplate(ctx context.Context, service *labels.Service) {
    var deleteLabelTemplate = &cobra.Command{
        Use:   "labels-delete-template",
        Short: "",
        RunE: func(_ *cobra.Command, _ []string) error {
            err := service.DeleteLabelTemplate(ctx, labelTemplateID)
            if err != nil {
                return err
            }

            fmt.Println("1")

            return nil
        },
    }

    rootCmd.AddCommand(deleteLabelTemplate)
}
