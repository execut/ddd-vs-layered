package domain_test

import (
	"testing"

	"effective-architecture/steps/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabelTemplate_Live(t *testing.T) {
	t.Parallel()

	var (
		labelTemplate *domain.LabelTemplate
		err           error
	)

	t.Run("New", func(t *testing.T) {
		labelTemplate, err = domain.NewLabelTemplate(expectedLabelTemplateID)

		require.NoError(t, err)
		assert.Equal(t, expectedLabelTemplateID, labelTemplate.ID)
	})

	t.Run("Create", func(t *testing.T) {
		err := labelTemplate.Create(expectedUserID, expectedManufacturer)

		require.NoError(t, err)
		assert.Equal(t, expectedManufacturer, labelTemplate.Manufacturer)
	})

	t.Run("Чтобы возвращалась уникальная ошибка при попытке создать уже существующий шаблон", func(t *testing.T) {
		err := labelTemplate.Create(expectedUserID, expectedManufacturer)

		require.ErrorIs(t, err, domain.ErrLabelTemplateAlreadyCreated)
	})

	t.Run("20. Чтобы нельзя было работать с чужими записями", func(t *testing.T) {
		err := labelTemplate.Update(expectedOtherUserID, expectedNewManufacturer)

		require.ErrorIs(t, err, domain.ErrLabelTemplateWrongUser)
	})

	t.Run("Обновлять данные шаблона", func(t *testing.T) {
		err := labelTemplate.Update(expectedUserID, expectedNewManufacturer)

		require.NoError(t, err)
		assert.Equal(t, expectedNewManufacturer, labelTemplate.Manufacturer)
	})

	t.Run("Удалять шаблон этикетки товара по UUID", func(t *testing.T) {
		err := labelTemplate.Delete(expectedUserID)

		require.NoError(t, err)
	})

	t.Run("Чтобы возвращалась уникальная ошибка при попытке удалить уже удалённый шаблон", func(t *testing.T) {
		err := labelTemplate.Delete(expectedUserID)

		require.ErrorIs(t, err, domain.ErrLabelTemplateAlreadyDeleted)
	})

	t.Run("CreateAgain", func(t *testing.T) {
		err := labelTemplate.Create(expectedUserID, expectedManufacturer)

		require.NoError(t, err)
		assert.Equal(t, expectedManufacturer, labelTemplate.Manufacturer)
	})

	t.Run("Привязывать шаблон к списку категорий или категорий+типов", func(t *testing.T) {
		err := labelTemplate.AddCategoryList(expectedUserID, []domain.Category{expectedCategory1, expectedCategory2})

		require.NoError(t, err)

		t.Run("и получать ошибку при попытке привязать уже существующую категорию", func(t *testing.T) {
			err := labelTemplate.AddCategoryList(expectedUserID, []domain.Category{expectedCategory1})

			require.ErrorIs(t, err, domain.ErrCategoryAlreadyAdded)
			require.ErrorContains(t, err, " (категория 1, тип 2)")
		})
	})

	t.Run("12. Отвязывать шаблон от списка категорий или категорий+типов", func(t *testing.T) {
		err = labelTemplate.UnlinkCategoryList(expectedUserID, []domain.Category{expectedCategory1, expectedCategory2})

		require.NoError(t, err)

		t.Run("и получать ошибку при попытке отвязать уже отвязанную категорию", func(t *testing.T) {
			err = labelTemplate.UnlinkCategoryList(expectedUserID, []domain.Category{expectedCategory2})

			require.Error(t, err)
			assert.ErrorContains(t, err, "категория уже отвязана от шаблона (категория 3, тип 4)")
		})
	})

	t.Run("19. Возможность", func(t *testing.T) {
		t.Run("деактивировать", func(t *testing.T) {
			err = labelTemplate.Deactivate(expectedUserID)

			require.NoError(t, err)
			assert.Equal(t, domain.LabelTemplateStatusDeactivated, labelTemplate.Status)
		})
		t.Run("и активировать шаблоны", func(t *testing.T) {
			err = labelTemplate.Activate(expectedUserID)

			require.NoError(t, err)
			assert.Equal(t, domain.LabelTemplateStatusCreated, labelTemplate.Status)
		})
	})
}
