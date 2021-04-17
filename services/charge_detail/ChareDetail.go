package charge_detail

import (
	"charge/container"
	"charge/dto"
	"charge/models"
	"errors"
)

func Add(accountId uint, _type uint8, categoryId uint, money int64, description string, repay uint8, repayId uint, repayAt int) (*models.ChargeDetail, error) {
	detail := new(models.ChargeDetail)
	detail.AccountId = accountId
	detail.Type = _type
	detail.CategoryId = categoryId
	detail.Money = money
	detail.Description = description
	detail.Repay = repay
	detail.RepayId = repayId
	detail.RepayAt = repayAt
	db := container.GetContainer().GetDb()
	result := db.Create(detail)

	return detail, result.Error
}

func Delete(id uint) (*models.ChargeDetail, error) {
	db := container.GetContainer().GetDb()
	detail := new(models.ChargeDetail)
	db.Where("id = ?", id).First(detail)
	if detail.ID == 0 {
		return nil, errors.New("删除对象不存在")
	}

	result := db.Delete(detail)

	return detail, result.Error
}

// EditDetail 编辑详情
func EditDetail(id int) (*dto.ChargeEditDetail, error) {
	db := container.GetContainer().GetDb()
	detailDto := new(dto.ChargeEditDetail)
	detail := new(models.ChargeDetail)
	result := db.Where("id = ?", id).First(detail)
	if result.Error != nil {
		return detailDto, result.Error
	}

	detailDto.ID = detail.ID
	detailDto.AccountId = detail.AccountId
	detailDto.Type = detail.Type
	detailDto.CategoryId = detail.CategoryId
	detailDto.Money = float64(detail.Money) / 1000.0
	detailDto.Description = detail.Description
	detailDto.Repay = detail.Repay == 1
	detailDto.RepayId = detail.RepayId
	detailDto.RepayAt = detail.RepayAt

	return detailDto, nil
}

// Edit 编辑分类
func Edit(id int, accountId uint, _type uint8, categoryId uint, money int64, description string, repay uint8, repayId uint, repayAt int) (*models.ChargeDetail, error) {
	db := container.GetContainer().GetDb()
	detail := new(models.ChargeDetail)
	db.Where("id = ?", id).First(detail)
	if detail.ID == 0 {
		return detail, errors.New("更新对象不存在")
	}

	detail.AccountId = accountId
	detail.Type = _type
	detail.CategoryId = categoryId
	detail.Money = money
	detail.Description = description
	detail.Repay = repay
	detail.RepayId = repayId
	detail.RepayAt = repayAt

	result := db.Save(detail)

	return detail, result.Error
}
