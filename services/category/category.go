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

// EditDetail 编辑详情
func EditDetail(id int) (*dto.Category, error) {
	db := container.GetContainer().GetDb()
	categoryDto := new(dto.Category)
	category := new(models.Category)
	result := db.Where("id = ?", id).First(category)
	if result.Error != nil {
		return categoryDto, result.Error
	}

	categoryDto.ID = category.ID
	categoryDto.Name = category.Name
	categoryDto.Type = category.Type
	categoryDto.Sort = category.Sort

	return categoryDto, nil
}

// Edit 编辑分类
func Edit(id int, _type uint8, name string, sort uint8) (*models.Category, error) {
	db := container.GetContainer().GetDb()
	category := new(models.Category)
	db.Where("id = ?", id).First(category)
	if category.ID == 0 {
		return category, errors.New("更新对象不存在")
	}

	category.Type = _type
	category.Name = name
	category.Sort = sort

	result := db.Save(category)

	return category, result.Error
}
