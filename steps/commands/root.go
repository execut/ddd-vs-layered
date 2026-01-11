package commands

import (
    "context"

    "effective-architecture/steps/labels"
    "github.com/spf13/cobra"
)

var (
    rootCmd = &cobra.Command{
        Use:   "",
        Short: "",
        Long:  ``,
    }
    labelTemplateID                 string
    manufacturerOrganizationName    string
    manufacturerOrganizationAddress string
    manufacturerEmail               string
    manufacturerSite                string
)

func Execute() error {
    ctx := context.Background()

    service, err := NewService(ctx)
    if err != nil {
        return err
    }

    InitRootCmd()
    InitCreateLabelTemplate(ctx, service)
    InitUpdateLabelTemplate(ctx, service)
    InitDeleteLabelTemplate(ctx, service)
    InitGetLabelTemplate(ctx, service)

    err = rootCmd.Execute()
    if err != nil {
        return err
    }

    return nil
}

func InitRootCmd() {
    rootCmd.PersistentFlags().StringVarP(&labelTemplateID, "id", "i", "", "id")
}

func NewService(ctx context.Context) (*labels.Service, error) {
    repository, err := labels.NewRepository(ctx)
    if err != nil {
        return nil, err
    }

    service := labels.NewService(repository)

    return service, nil
}

func initManufacturerFlags(cmd *cobra.Command) {
    cmd.PersistentFlags().StringVarP(&manufacturerOrganizationName, "manufacturer-organization-name", "m",
        "", "manufacturer-organization-name")
    cmd.Flags().StringVarP(&manufacturerOrganizationAddress, "manufacturer-organization-address", "a",
        "", "manufacturer-organization-address")
    cmd.Flags().StringVarP(&manufacturerEmail, "manufacturer-email", "e",
        "", "manufacturer-email")
    cmd.Flags().StringVarP(&manufacturerSite, "manufacturer-site", "s",
        "", "manufacturer-site")
}
