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
            err := service.UpdateLabelTemplate(ctx, labelTemplateID, labels.Manufacturer{
                OrganizationName:    manufacturerOrganizationName,
                OrganizationAddress: manufacturerOrganizationAddress,
                Email:               manufacturerEmail,
                Site:                manufacturerSite,
            })
            if err != nil {
                return err
            }

            fmt.Println("1")

            return nil
        },
    }

    initManufacturerFlags(updateLabelTemplate)
    rootCmd.AddCommand(updateLabelTemplate)
}
