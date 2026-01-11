package commands

import (
    "context"

    "effective-architecture/steps/application"
    "effective-architecture/steps/infrastructure"
    "effective-architecture/steps/infrastructure/history"
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
)

func Execute() error {
    ctx := context.Background()

    repository, err := infrastructure.NewEventsRepository()
    if err != nil {
        panic(err)
    }

    historyRepository, err := history.NewRepository()
    if err != nil {
        panic(err)
    }

    app, err := application.NewApplication(repository, historyRepository)
    if err != nil {
        panic(err)
    }

    err = InitLabelsCreateTemplate(ctx, app)
    if err != nil {
        return err
    }

    err = InitLabelsUpdateTemplate(ctx, app)
    if err != nil {
        return err
    }

    err = InitLabelsDeleteTemplate(ctx, app)
    if err != nil {
        return err
    }

    err = InitLabelsGetTemplate(ctx, app)
    if err != nil {
        return err
    }

    err = InitLabelsTemplateHistory(ctx, app)
    if err != nil {
        return err
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
