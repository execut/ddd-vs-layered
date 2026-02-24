package commands

import (
	"context"
	"fmt"

	"effective-architecture/steps/contract"

	"github.com/spf13/cobra"
)

func InitLabelsCreateTemplate(ctx context.Context, app contract.IApplication) error {
	var createLabelTemplateCmd = &cobra.Command{
		Use:   "labels-create-template",
		Short: "",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			err := app.Create(ctx, userID, labelTemplateID, contract.Manufacturer{
				OrganizationName:    organizationName,
				OrganizationAddress: organizationAddress,
				Email:               email,
				Site:                site,
			})
			if err != nil {
				panic(err)
			}

			fmt.Println("1")
		},
	}

	initManufacturerFlags(createLabelTemplateCmd)
	rootCmd.AddCommand(createLabelTemplateCmd)

	return nil
}
