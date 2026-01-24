package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/contract"
    "github.com/spf13/cobra"
)

func InitLabelsGetTemplate(ctx context.Context, app contract.IApplication) error {
    var getLabelTemplateCmd = &cobra.Command{
        Use:   "labels-get-template",
        Short: "",
        Long:  ``,
        Run: func(_ *cobra.Command, _ []string) {
            result, err := app.Get(ctx, labelTemplateID)
            if err != nil {
                panic(err)
            }

            fmt.Println(result)
        },
    }

    getLabelTemplateCmd.PersistentFlags().StringVarP(&organizationName, "labels-get-template",
        "g", "", "labels-get-template")

    rootCmd.AddCommand(getLabelTemplateCmd)

    return nil
}
