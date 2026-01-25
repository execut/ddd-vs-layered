package commands

import (
    "context"
    "errors"
    "fmt"
    "strings"

    "effective-architecture/steps/contract"

    "github.com/spf13/cobra"
)

var ErrWrongCategory = errors.New("wrong category")

const (
    categoryWithTypeLen = 2
)

func InitLabelsTemplateAddCategoryList(ctx context.Context, app contract.IApplication) error {
    var categoryList []string

    var labelsTemplateAddCategoryListCmd = &cobra.Command{
        Use:   "labels-template-add-category-list",
        Short: "",
        Long:  ``,
        RunE: func(_ *cobra.Command, _ []string) error {
            appCategoryList := make([]contract.Category, 0, len(categoryList))
            for _, categoryID := range categoryList {
                categoryIDParts := strings.Split(categoryID, "-")

                var typeID *string

                if len(categoryIDParts) == 0 {
                    return ErrWrongCategory
                }

                if len(categoryIDParts) == categoryWithTypeLen {
                    typeID = &categoryIDParts[1]
                }

                appCategoryList = append(appCategoryList, contract.Category{
                    CategoryID: categoryIDParts[0],
                    TypeID:     typeID,
                })
            }

            err := app.AddCategoryList(ctx, labelTemplateID, appCategoryList)
            if err != nil {
                return err
            }

            fmt.Println("1")

            return nil
        },
    }

    labelsTemplateAddCategoryListCmd.Flags().StringSliceVar(&categoryList, "category-id-list", []string{},
        "category list")

    rootCmd.AddCommand(labelsTemplateAddCategoryListCmd)

    return nil
}
