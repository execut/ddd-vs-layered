package commands

import (
    "context"

    "effective-architecture/steps/application"
    "github.com/spf13/cobra"
)

var (
    rootCmd = &cobra.Command{
        Use:   "",
        Short: "",
        Long:  ``,
    }
    labelTemplateID     string
    organizationName    string
    organizationAddress string
    site                string
    email               string
    initiators          = []func(ctx context.Context, app *application.Application) error{
        InitLabelsCreateTemplate,
        InitLabelsDeleteTemplate,
        InitLabelsGetTemplate,
        InitLabelsTemplateAddCategoryList,
        InitLabelsTemplateHistory,
        InitLabelsUpdateTemplate,
    }
)

func Execute() error {
    ctx := context.Background()

    app, err := application.NewApplication()
    if err != nil {
        panic(err)
    }

    for _, initiator := range initiators {
        err := initiator(ctx, app)
        if err != nil {
            return err
        }
    }

    rootCmd.PersistentFlags().StringVarP(&labelTemplateID, "id", "i", "", "id")

    err = rootCmd.Execute()
    if err != nil {
        return err
    }

    return nil
}

func initManufacturerFlags(cmd *cobra.Command) {
    cmd.PersistentFlags().StringVarP(&organizationName, "manufacturer-organization-name",
        "m", "", "manufacturer-organization-name")
    cmd.Flags().StringVarP(&organizationAddress, "manufacturer-organization-address",
        "a", "", "manufacturer-organization-address")
    cmd.Flags().StringVarP(&email, "manufacturer-email",
        "e", "", "manufacturer-email")
    cmd.Flags().StringVarP(&site, "manufacturer-site",
        "s", "", "manufacturer-site")
}
