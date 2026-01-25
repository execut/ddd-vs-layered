package domain

type Category struct {
    CategoryID int64
    TypeID     *int64
}

func NewCategory(categoryID int64, typeID *int64) (Category, error) {
    return Category{
        CategoryID: categoryID,
        TypeID:     typeID,
    }, nil
}

func (c Category) Same(category Category) bool {
    if c.TypeID == nil && category.TypeID == nil || c.TypeID != nil && category.TypeID != nil &&
        *c.TypeID == *category.TypeID {
        return c.CategoryID == category.CategoryID
    }

    return false
}
