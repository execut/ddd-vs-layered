package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/application"
    "github.com/spf13/cobra"
)

func InitLabelsCreateTemplate(ctx context.Context, app *application.Application) error {
    var createLabelTemplateCmd = &cobra.Command{
        Use:   "labels-create-template",
        Short: "",
        Long:  ``,
        Run: func(_ *cobra.Command, _ []string) {
            err := app.Create(ctx, labelTemplateID, application.Manufacturer{
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
