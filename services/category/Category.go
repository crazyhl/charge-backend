package category

import (
	"charge/container"
	"charge/dto"
	"charge/models"
	"errors"
)

func ListGroup() dto.CategoryGroup {
	categoryGroup := make(dto.CategoryGroup)
	var categoryRows []models.Category
	db := container.GetContainer().GetDb()
	db.Order("type ASC").Order("sort DESC").Order("id ASC").Find(&categoryRows)
	for _, category := range categoryRows {
		_, found := categoryGroup[category.Type]
		if found == false {
			categoryGroup[category.Type] = make([]dto.Category, 0, 0)
		}
		categoryGroup[category.Type] = append(categoryGroup[category.Type], dto.Category(category))
	}

	return categoryGroup
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

func Delete(id uint) (*models.Category, error) {
	db := container.GetContainer().GetDb()
	category := new(models.Category)
	db.Where("id = ?", id).First(category)
	if category.ID == 0 {
		return nil, errors.New("删除对象不存在")
	}

	result := db.Delete(category)

	return category, result.Error
}
