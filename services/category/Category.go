package category

import (
	"charge/container"
	"charge/dto"
	"charge/models"
)

func ListGroup() []dto.CategoryGroup {
	var categoryGroupRows []dto.CategoryGroup
	var categoryRows []models.Category
	db := container.GetContainer().GetDb()
	db.Order("type ASC").Order("sort DESC").Find(&categoryRows)
	for _, category := range categoryRows {
		intType := int(category.Type)
		if len(categoryGroupRows) < intType+1 {
			var groupCategoryRows []models.Category
			groupCategoryRows = append(groupCategoryRows, category)
			categoryGroupRows = append(categoryGroupRows, dto.CategoryGroup{
				Type:       category.Type,
				Categories: groupCategoryRows,
			})
		} else {
			categoryGroupRows[intType].Categories = append(categoryGroupRows[intType].Categories, category)
		}
	}

	return categoryGroupRows
}

func Add(_type uint8, name string, sort uint8) (*models.Category, error) {
	category := new(models.Category)
	category.Type = _type
	category.Name = name
	category.Sort = sort
	db := container.GetContainer().GetDb()
	result := db.Create(category)

	return category, result.Error
}
