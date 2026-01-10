package commands

import (
    "context"

    labels2 "effective-architecture/steps/labels"
    "github.com/spf13/cobra"
)

var (
    rootCmd = &cobra.Command{
        Use:   "",
        Short: "",
        Long:  ``,
    }
    labelTemplateID string
    name            string
)

func Execute() {
    InitRootCmd()
    InitCreateLabelTemplate()

    err := rootCmd.Execute()
    if err != nil {
        panic(err)
    }
}

func InitRootCmd() {
    rootCmd.PersistentFlags().StringVarP(&labelTemplateID, "id", "i", "", "id")
    rootCmd.AddCommand(createLabelTemplate)
    rootCmd.AddCommand(getLabelTemplate)
}

func NewService(ctx context.Context) (*labels2.Service, error) {
    repository, err := labels2.NewRepository(ctx)
    if err != nil {
        return nil, err
    }

    service := labels2.NewService(repository)

    return service, nil
}
