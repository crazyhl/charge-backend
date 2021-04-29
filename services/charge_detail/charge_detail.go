package charge_detail

import (
	"charge/container"
	"charge/dto"
	"charge/models"
	"errors"
	"gorm.io/gorm/clause"
	"time"
)

func List(pageStart, pageSize int) dto.ListData {
	var listData dto.ListData
	var details []dto.ChargeDetail
	var detailRows []models.ChargeDetail
	var totalCount int64
	db := container.GetContainer().GetDb()
	db.Limit(pageSize).Offset(pageStart).Order("create_at DESC").Preload(clause.Associations).Find(&detailRows)
	db.Model(&models.ChargeDetail{}).Count(&totalCount)
	for _, detail := range detailRows {
		createTm := time.Unix(detail.CreateAt, 0)
		updateTm := time.Unix(detail.UpdateAt, 0)

		chargeDetailDto := dto.ChargeDetail{
			ID:        detail.ID,
			AccountId: detail.AccountId,
			Account: dto.AccountDetail{
				ID:        detail.Account.ID,
				Name:      detail.Account.Name,
				HasCredit: detail.Account.HasCredit == 1,
				Cash:      float64(detail.Account.Cash) / 1000.0,
				Credit:    float64(detail.Account.Credit) / 1000.0,
			},
			Type:        detail.Type,
			Money:       float64(detail.Money) / 1000.0,
			Description: detail.Description,
			CreateAt:    createTm.Format("2006-01-02"),
			UpdateAt:    updateTm.Format("2006-01-02"),
		}

		if detail.CategoryId != 0 {
			chargeDetailDto.Category = &dto.Category{
				ID:   detail.Category.ID,
				Type: detail.Category.Type,
				Name: detail.Category.Name,
			}
		}

		if detail.RepaidDetailId != 0 {
			repaidDetailCreateTm := time.Unix(detail.RepaidDetail.CreateAt, 0)
			chargeDetailDto.RepaidDetail = &dto.RepaidDetail{
				ID:       detail.RepaidDetail.RepaidDetailId,
				CreateAt: repaidDetailCreateTm.Format("2006-01-02"),
			}
		}
		if detail.RepayAccountId != 0 {
			chargeDetailDto.RepayAccount = &dto.AccountDetail{
				ID:        detail.RepayAccount.ID,
				Name:      detail.RepayAccount.Name,
				HasCredit: detail.RepayAccount.HasCredit == 1,
				Cash:      float64(detail.RepayAccount.Cash) / 1000.0,
				Credit:    float64(detail.RepayAccount.Credit) / 1000.0,
			}
		}
		if detail.TransferAccountId != 0 {
			chargeDetailDto.TransferAccount = &dto.AccountDetail{
				ID:        detail.TransferAccount.ID,
				Name:      detail.TransferAccount.Name,
				HasCredit: detail.TransferAccount.HasCredit == 1,
				Cash:      float64(detail.TransferAccount.Cash) / 1000.0,
				Credit:    float64(detail.TransferAccount.Credit) / 1000.0,
			}
		}
		details = append(details, chargeDetailDto)
	}
	listData.Total = totalCount
	listData.Data = details
	return listData
}

func Add(accountId uint, _type uint8, categoryId uint, money int64, description string, repayAccountId uint, transferAccountId uint, createAt int64) (*models.ChargeDetail, error) {
	detail := new(models.ChargeDetail)
	detail.AccountId = accountId
	detail.Type = _type
	detail.CategoryId = categoryId
	detail.Money = money
	detail.Description = description
	detail.RepayAccountId = repayAccountId
	detail.TransferAccountId = transferAccountId
	detail.CreateAt = createAt
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
	detailDto.RepayAccountId = detail.RepayAccountId
	detailDto.TransferAccountId = detail.TransferAccountId
	if detail.Type == 3 {
		detailDto.RepaidDetails = GetRepaidList(detail.ID)
	}
	return detailDto, nil
}

func GetUnPaidList(accountId uint) []dto.UnpaidDetail {
	db := container.GetContainer().GetDb()
	unPaidDetailDtoList := make([]dto.UnpaidDetail, 0, 0)
	var unPayDetails []models.ChargeDetail
	db.Where("account_id = ?", accountId).
		Where("type = ?", 2).
		Where("repaid_detail_id = ?", 0).
		Preload(clause.Associations).Find(&unPayDetails)
	for _, detail := range unPayDetails {
		createTm := time.Unix(detail.CreateAt, 0)
		unpaidDetail := new(dto.UnpaidDetail)
		unpaidDetail.ID = detail.ID
		unpaidDetail.Category = &dto.Category{
			ID:   detail.Category.ID,
			Type: detail.Category.Type,
			Name: detail.Category.Name,
		}
		unpaidDetail.Money = float64(detail.Money) / 1000.0
		unpaidDetail.Description = detail.Description
		unpaidDetail.CreateAt = createTm.Format("2006-01-02")
		unPaidDetailDtoList = append(unPaidDetailDtoList, *unpaidDetail)
	}

	return unPaidDetailDtoList
}

// Edit 编辑账单
func Edit(id uint, accountId uint, _type uint8, categoryId uint, money int64, description string, repayAccountId uint, transferAccountId uint, createAt int64) (*models.ChargeDetail, error) {
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
	detail.RepayAccountId = repayAccountId
	detail.TransferAccountId = transferAccountId
	detail.CreateAt = createAt

	result := db.Save(detail)

	return detail, result.Error
}

func GetRepaidList(id uint) []dto.RepaidDetail {
	db := container.GetContainer().GetDb()
	paidDtoList := make([]dto.RepaidDetail, 0, 0)
	var list []models.ChargeDetail
	db.Model(models.ChargeDetail{}).Where(
		"id = ?",
		id,
	).Find(&list)

	for _, detail := range list {
		createTm := time.Unix(detail.CreateAt, 0)
		paidDtoList = append(paidDtoList, dto.RepaidDetail{
			ID:    detail.RepaidDetailId,
			Money: float64(detail.Money) / 1000.0,
			Category: &dto.Category{
				ID:   detail.Category.ID,
				Type: detail.Category.Type,
				Name: detail.Category.Name,
			},
			Description: detail.Description,
			CreateAt:    createTm.Format("2006-01-02"),
		})
	}

	return paidDtoList
}

// UpdateRepay 更新还款记录
func UpdateRepay(repaidId uint, chargeIdArr []uint) {
	db := container.GetContainer().GetDb()
	db.Model(models.ChargeDetail{}).Where(
		"id in ?",
		chargeIdArr,
	).Updates(
		models.ChargeDetail{RepaidDetailId: repaidId},
	)
}

// ClearRepay 清除借款的还款记录
func ClearRepay(id uint) {
	db := container.GetContainer().GetDb()
	db.Model(models.ChargeDetail{}).Where(
		"repay_detail_id in ?",
		id,
	).Updates(
		models.ChargeDetail{RepaidDetailId: 0},
	)
}
