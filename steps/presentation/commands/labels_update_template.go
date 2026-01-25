package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/contract"
    "github.com/spf13/cobra"
)

func InitLabelsUpdateTemplate(ctx context.Context, app contract.IApplication) error {
    var updateLabelTemplateCmd = &cobra.Command{
        Use:   "labels-update-template",
        Short: "",
        Long:  ``,
        Run: func(_ *cobra.Command, _ []string) {
            err := app.Update(ctx, labelTemplateID, contract.Manufacturer{
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

    initManufacturerFlags(updateLabelTemplateCmd)
    rootCmd.AddCommand(updateLabelTemplateCmd)

    return nil
}
