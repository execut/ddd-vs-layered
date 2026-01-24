package commands

import (
    "context"
    "fmt"
    "strings"

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
            appCategoryList := make([]application.Category, 0, len(categoryList))
            for _, categoryID := range categoryList {
                categoryIDParts := strings.Split(categoryID, "-")
                appCategoryList = append(appCategoryList, application.Category{
                    CategoryID: categoryIDParts[0],
                    TypeID:     &categoryIDParts[1],
                })
            }
            err := app.AddCategoryList(ctx, labelTemplateID, appCategoryList)
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
