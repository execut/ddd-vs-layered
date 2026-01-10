package commands

import (
    "context"

    "effective-architecture/steps/application"
    "effective-architecture/steps/infrastructure"
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

func Execute() error {
    ctx := context.Background()

    repository, err := infrastructure.NewEventsRepository()
    if err != nil {
        panic(err)
    }

    app, err := application.NewApplication(repository)
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

    rootCmd.PersistentFlags().StringVarP(&labelTemplateID, "id", "i", "", "id")

    err = rootCmd.Execute()
    if err != nil {
        return err
    }

    return nil
}
