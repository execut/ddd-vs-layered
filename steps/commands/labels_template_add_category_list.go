package commands

import (
    "context"
    "fmt"

    "effective-architecture/steps/application"
    "github.com/spf13/cobra"
)

func InitLabelsTemplateAddCategoryList(ctx context.Context, app *application.Application) error {
    var categoryList []string

    var labelsTemplateAddCategoryListCmd = &cobra.Command{
        Use:   "labels-template-add-category-list",
        Short: "",
        Long:  ``,
        Run: func(_ *cobra.Command, _ []string) {
            err := app.LabelTemplateAddCategoryList(ctx, labelTemplateID, categoryList)
            if err != nil {
                panic(err)
            }

            fmt.Println("1")
        },
    }

    labelsTemplateAddCategoryListCmd.Flags().StringSliceVar(&categoryList, "category-id-list", []string{},
        "category list")

    rootCmd.AddCommand(labelsTemplateAddCategoryListCmd)

    return nil
}
