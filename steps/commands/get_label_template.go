package commands

import (
    "context"
    "fmt"

    "github.com/spf13/cobra"
)

var (
    getLabelTemplate = &cobra.Command{
        Use:   "labels-get-template",
        Short: "",
        RunE: func(_ *cobra.Command, _ []string) error {
            ctx := context.Background()
            service, err := NewService(ctx)
            if err != nil {
                return err
            }

            result, err := service.GetLabelTemplate(ctx, labelTemplateID)
            if err != nil {
                return err
            }

            fmt.Println(result)

            return nil
        },
    }
)
