package charge_detail

import (
	"charge/container"
	"charge/dto"
	"charge/models"
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
)

func List(pageStart, pageSize int) []dto.ChargeDetail {
	var details []dto.ChargeDetail
	var detailRows []models.ChargeDetail
	db := container.GetContainer().GetDb()
	db.Order("id DESC").Limit(pageStart).Offset(pageSize).Preload(clause.Associations).Find(&detailRows)
	fmt.Println(detailRows)
	//for _, detail := range detailRows {
	//	createTm := time.Unix(detail.CreateAt, 0)
	//	updateTm := time.Unix(detail.UpdateAt, 0)
	//	repayTm := time.Unix(detail.RepayAt, 0)
	//
	//	repayDetail := nil
	//
	//	details = append(details, dto.ChargeDetail{
	//		ID:        detail.ID,
	//		AccountId: detail.AccountId,
	//		Account: dto.AccountList{
	//			ID:        detail.Account.ID,
	//			Name:      detail.Account.Name,
	//			HasCredit: detail.Account.HasCredit == 1,
	//			Cash:      float64(detail.Account.Cash) / 1000.0,
	//			Credit:    float64(detail.Account.Cash) / 1000.0,
	//		},
	//		Type:      detail.Type,
	//		CategoryId: detail.CategoryId,
	//		Category:    dto.Category{
	//			ID:   detail.Category.ID,
	//			Type: detail.Category.Type,
	//			Name: detail.Category.Name,
	//		},
	//		Money:      float64(detail.Money) / 1000.0,
	//		Description:    detail.Description,
	//		Repay:      detail.Repay == 1,
	//		ReplayDetail: repayDetail,
	//		CreateAt:  createTm.Format("2006-01-02 15:04:05"),
	//		UpdateAt:  updateTm.Format("2006-01-02 15:04:05"),
	//		RepayAt:  repayTm.Format("2006-01-02 15:04:05"),
	//	})
	//}
	return details
}

func Add(accountId uint, _type uint8, categoryId uint, money int64, description string, repayAt int64, repayAccountId uint, transfer uint8, transferAccountId uint) (*models.ChargeDetail, error) {
	detail := new(models.ChargeDetail)
	detail.AccountId = accountId
	detail.Type = _type
	detail.CategoryId = categoryId
	detail.Money = money
	detail.Description = description
	detail.RepayAt = repayAt
	detail.RepayAccountId = repayAccountId
	detail.Transfer = transfer
	detail.TransferAccountId = transferAccountId
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

// EditDetail 编辑账单详情
func EditDetail(id int) (*dto.ChargeEditDetail, error) {
	db := container.GetContainer().GetDb()
	detailDto := new(dto.ChargeEditDetail)
	detail := new(models.ChargeDetail)
	result := db.Where("id = ?", id).Preload(clause.Associations).First(detail)
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
	detailDto.RepayDetailId = detail.RepayDetailId
	detailDto.RepayAt = detail.RepayAt
	detailDto.Transfer = detail.Transfer == 1
	detailDto.TransferAccountId = detail.TransferAccountId

	return detailDto, nil
}

func GetUnPayList() []dto.ChargeDetail {
	db := container.GetContainer().GetDb()
	unpayDetailDtos := make([]dto.ChargeDetail, 0, 0)
	var unPayDetails []models.ChargeDetail
	db.Where("repay = ?", 0).Preload(clause.Associations).Find(&unPayDetails)
	for _, detail := range unPayDetails {
		unpayDetailDtos = append(unpayDetailDtos, dto.ChargeDetail{
			ID:         detail.ID,
			AccountId:  detail.AccountId,
			Type:       detail.Type,
			CategoryId: detail.CategoryId,
			Category: dto.Category{
				ID:   detail.Category.ID,
				Type: detail.Category.Type,
				Name: detail.Category.Name,
			},
			Money:       float64(detail.Money) / 1000.0,
			Description: detail.Description,
		})
	}

	return unpayDetailDtos
}

// Edit 编辑账单
func Edit(id uint, accountId uint, _type uint8, categoryId uint, money int64, description string, repayAt int64, repayAccountId uint, transfer uint8, transferAccountId uint) (*models.ChargeDetail, error) {
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
	detail.RepayDetailId = repayAccountId
	detail.RepayAt = repayAt
	detail.RepayAccountId = repayAccountId
	detail.Transfer = transfer
	detail.TransferAccountId = transferAccountId

	result := db.Save(detail)

	return detail, result.Error
}

// UpdateRepay 更新还款记录
func UpdateRepay(repayId uint, repayTime int64, chargeIdArr []int) {
	db := container.GetContainer().GetDb()
	db.Model(models.ChargeDetail{}).Where(
		"id in ?",
		chargeIdArr,
	).Updates(
		models.ChargeDetail{RepayDetailId: repayId, RepayAt: repayTime},
	)
}
