package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/labels"
    "github.com/spf13/cobra"
)

func InitTemplateHistory(ctx context.Context, service *labels.Service) {
    var templateHistory = &cobra.Command{
        Use:   "labels-template-history",
        Short: "",
        RunE: func(_ *cobra.Command, _ []string) error {
            result, err := service.GetLabelHistory(ctx, labelTemplateID)
            if err != nil {
                return err
            }

            fmt.Println(result)

            return nil
        },
    }

    rootCmd.AddCommand(templateHistory)
}
