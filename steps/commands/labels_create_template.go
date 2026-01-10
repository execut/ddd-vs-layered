package commands

import (
    "context"
    "fmt"

    "github.com/spf13/cobra"
)

func InitCreateLabelTemplate() {
    createLabelTemplate.PersistentFlags().StringVarP(&name, "manufacturer-organization-name", "m",
        "", "manufacturer-organization-name")
    rootCmd.AddCommand(createLabelTemplate)
}

var (
    createLabelTemplate = &cobra.Command{
        Use:   "labels-create-template",
        Short: "",
        RunE: func(_ *cobra.Command, _ []string) error {
            ctx := context.Background()
            service, err := NewService(ctx)
            if err != nil {
                return err
            }

            err = service.CreateLabelTemplate(ctx, labelTemplateID, name)
            if err != nil {
                return err
            }

            fmt.Println("1")

            return nil
        },
    }
)
