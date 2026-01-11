package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/application"
    "github.com/spf13/cobra"
)

func InitLabelsTemplateHistory(ctx context.Context, app *application.Application) error {
    var labelTemplateHistoryCmd = &cobra.Command{
        Use:   "labels-template-history",
        Short: "",
        Long:  ``,
        RunE: func(_ *cobra.Command, _ []string) error {
            result, err := app.LabelTemplateHistoryList(ctx, labelTemplateID)
            if err != nil {
                return err
            }

            fmt.Println(result)

            return nil
        },
    }

    initManufacturerFlags(labelTemplateHistoryCmd)
    rootCmd.AddCommand(labelTemplateHistoryCmd)

    return nil
}
