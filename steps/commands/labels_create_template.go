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
            err := service.CreateLabelTemplate(ctx, labelTemplateID, labels.Manufacturer{
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

    initManufacturerFlags(createLabelTemplate)
    rootCmd.AddCommand(createLabelTemplate)
}
