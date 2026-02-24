package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"effective-architecture/steps/contract"

	"github.com/spf13/cobra"
)

func InitLabelsGetTemplate(ctx context.Context, app contract.IApplication) error {
	var getLabelTemplateCmd = &cobra.Command{
		Use:   "labels-get-template",
		Short: "",
		Long:  ``,
		RunE: func(_ *cobra.Command, _ []string) error {
			result, err := app.Get(ctx, userID, labelTemplateID)
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

	getLabelTemplateCmd.PersistentFlags().StringVarP(&organizationName, "labels-get-template",
		"g", "", "labels-get-template")

	rootCmd.AddCommand(getLabelTemplateCmd)

	return nil
}
