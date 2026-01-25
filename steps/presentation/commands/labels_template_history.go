package commands

import (
    "context"
    "encoding/json"
    "fmt"

    "effective-architecture/steps/contract"
    "github.com/spf13/cobra"
)

func InitLabelsTemplateHistory(ctx context.Context, app contract.IApplication) error {
    var labelTemplateHistoryCmd = &cobra.Command{
        Use:   "labels-template-history",
        Short: "",
        Long:  ``,
        RunE: func(_ *cobra.Command, _ []string) error {
            result, err := app.HistoryList(ctx, labelTemplateID)
            if err != nil {
                return err
            }

            resultJSON, err := json.Marshal(result)
            if err != nil {
                return err
            }

            fmt.Println(string(resultJSON))

            return nil
        },
    }

    initManufacturerFlags(labelTemplateHistoryCmd)
    rootCmd.AddCommand(labelTemplateHistoryCmd)

    return nil
}
