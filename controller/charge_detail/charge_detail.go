package charge_detail

import (
	"charge/container"
	"charge/models"
	accountService "charge/services/account"
	"charge/services/charge_detail"
	"charge/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/now"
	"strconv"
	"time"
)

type addDetail struct {
	AccountId         uint    `json:"account_id" form:"account_id" validate:"required,number"`
	Type              *uint8  `json:"type" form:"type"  validate:"required,number,gte=0,lte=4"`
	CategoryId        uint    `json:"category_id" form:"category_id" validate:"number,gte=0"`
	Money             float64 `json:"money" form:"money" validate:"required,number"`
	Description       string  `json:"description" form:"description"`
	RepayDetailIds    []uint  `json:"repay_detail_ids" form:"repay_detail_ids"`
	RepayAccountId    uint    `json:"repay_account_id" form:"repay_account_id"`
	TransferAccountId uint    `json:"transfer_account_id" form:"transfer_account_id"`
	CreateAt          int64   `json:"date" form:"date"`
}

type editDetail struct {
	Id                uint    `json:"id" form:"id" validate:"required"  validate:"required"`
	AccountId         uint    `json:"account_id" form:"account_id" validate:"required,number"`
	Type              *uint8  `json:"type" form:"type"  validate:"required,number,gte=0,lte=4"`
	CategoryId        uint    `json:"category_id" form:"category_id" validate:"number,gte=0"`
	Money             float64 `json:"money" form:"money" validate:"required,number"`
	Description       string  `json:"description" form:"description"`
	RepayDetailIds    []uint  `json:"repay_detail_ids" form:"repay_detail_ids"`
	RepayAccountId    uint    `json:"repay_account_id" form:"repay_account_id"`
	TransferAccountId uint    `json:"transfer_account_id" form:"transfer_account_id"`
	CreateAt          int64   `json:"date" form:"date"`
}

func List(ctx *fiber.Ctx) error {
	page := ctx.Query("page", "1")
	size := ctx.Query("pageSize", "20")
	month := ctx.Params("month", "")
	category, _ := ctx.ParamsInt("category")

	var beginningOfMonth time.Time
	var endOfMonth time.Time
	if month != "" {
		location, _ := time.LoadLocation("Asia/Shanghai")
		monthTime, _ := time.Parse("20060102 -0700 MST", month+"01 +0800 CST")

		myConfig := &now.Config{
			TimeLocation: location,
		}

		beginningOfMonth = myConfig.With(monthTime).BeginningOfMonth()
		endOfMonth = myConfig.With(monthTime).EndOfMonth()
	}

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(size)
	if pageSizeInt <= 0 {
		pageSizeInt = 20
	}
	if pageInt < 1 {
		pageInt = 1
	}
	pageStart := (pageInt - 1) * pageSizeInt
	listData := charge_detail.List(pageStart, pageSizeInt, category, beginningOfMonth.Unix(), endOfMonth.Unix())

	return ctx.JSON(fiber.Map{
		"status": 0,
		"data":   listData,
	})
}

func Add(ctx *fiber.Ctx) error {
	db := container.GetContainer().GetDb()

	detail := new(addDetail)
	if err := ctx.BodyParser(detail); err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -1,
			"message": err.Error(),
		})
	}
	// ??????
	validateError := utils.Validate(detail)
	if validateError != nil {
		fmt.Println(validateError)
		return ctx.JSON(fiber.Map{
			"status":  -3,
			"message": validateError.Error(),
		})
	}

	if *detail.Type == uint8(3) && detail.RepayAccountId == 0 {
		// ??????????????????,
		return ctx.JSON(fiber.Map{
			"status":  -4,
			"message": "?????????????????????",
		})
	}

	if *detail.Type == uint8(4) && detail.TransferAccountId == 0 {
		// ??????????????????,
		return ctx.JSON(fiber.Map{
			"status":  -5,
			"message": "?????????????????????",
		})
	}

	// ??????????????????????????????????????????????????????
	account := new(models.Account)
	db.Where("id =?", detail.AccountId).First(account)
	if account.ID == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -6,
			"message": "???????????????",
		})
	}
	if *detail.Type == uint8(2) && account.HasCredit == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -7,
			"message": "????????????????????????",
		})
	}

	if *detail.Type == uint8(3) {
		repayAccount := new(models.Account)
		db.Where("id =?", detail.RepayAccountId).First(repayAccount)
		if repayAccount.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "?????????????????????",
			})
		}
	}

	if *detail.Type == uint8(4) {
		transferAccount := new(models.Account)
		db.Where("id =?", detail.TransferAccountId).First(transferAccount)
		if transferAccount.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "?????????????????????",
			})
		}
	}

	if detail.CategoryId > 0 {
		cate := new(models.Category)
		db.Where("id =?", detail.CategoryId).Where("type =?", detail.Type).First(cate)
		if cate.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -9,
				"message": "?????????????????????",
			})
		}
	}

	money := int64(detail.Money * 1000)

	chargeDetail, err := charge_detail.Add(
		detail.AccountId,
		*detail.Type,
		detail.CategoryId,
		money,
		detail.Description,
		detail.RepayAccountId,
		detail.TransferAccountId,
		detail.CreateAt,
	)

	if err == nil {
		// ???????????????????????????????????????
		switch *detail.Type {
		case 0:
			// ?????? ???????????? cashIn
			accountService.SummaryMoney("cashIn", detail.AccountId, time.Unix(detail.CreateAt, 0))
		case 1:
			// ?????? ???????????? cashOut
			accountService.SummaryMoney("cashOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
		case 2:
			// ??? ???????????? creditOut
			accountService.SummaryMoney("creditOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
		case 3:
			// ??? ???????????? cashOut ?????? ???????????? creditIn
			accountService.SummaryMoney("cashOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
			accountService.SummaryMoney("creditIn", detail.RepayAccountId, time.Unix(detail.CreateAt, 0))
			// ???????????????????????????id
			charge_detail.UpdateRepay(chargeDetail.ID, detail.RepayDetailIds)
		case 4:
			// ??? ???????????? transferOut  ?????? ???????????? transferIn
			accountService.SummaryMoney("transferOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
			accountService.SummaryMoney("transferIn", detail.TransferAccountId, time.Unix(detail.CreateAt, 0))
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  0,
		"message": "????????????",
	})
}

// Delete ????????????
func Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -1,
			"message": err.Error(),
		})
	}
	detail, err := charge_detail.Delete(uint(id))
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -2,
			"message": err.Error(),
		})
	}

	switch detail.Type {
	case 0:
		// ?????? ???????????? cashIn
		accountService.SummaryMoney("cashIn", detail.AccountId, time.Unix(detail.CreateAt, 0))
	case 1:
		// ?????? ???????????? cashOut
		accountService.SummaryMoney("cashOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
	case 2:
		// ??? ???????????? creditOut
		accountService.SummaryMoney("creditOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
	case 3:
		// ??? ???????????? cashOut ?????? ???????????? creditIn
		accountService.SummaryMoney("cashOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
		accountService.SummaryMoney("creditIn", detail.RepayAccountId, time.Unix(detail.CreateAt, 0))
		// ???????????????????????????id
		charge_detail.ClearRepay(detail.ID)
	case 4:
		// ??? ???????????? transferOut  ?????? ???????????? transferIn
		accountService.SummaryMoney("transferOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
		accountService.SummaryMoney("transferIn", detail.TransferAccountId, time.Unix(detail.CreateAt, 0))
	}

	return ctx.JSON(fiber.Map{
		"status":  0,
		"message": "????????????",
	})
}

// EditDetail ?????????????????????
func EditDetail(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -1,
			"message": err.Error(),
		})
	}

	detail, err := charge_detail.EditDetail(id)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -2,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status": 0,
		"data":   detail,
	})
}

func Edit(ctx *fiber.Ctx) error {
	db := container.GetContainer().GetDb()

	detail := new(editDetail)
	if err := ctx.BodyParser(detail); err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -1,
			"message": err.Error(),
		})
	}

	oldDetail := new(models.ChargeDetail)
	db.Where("id = ?", detail.Id).First(oldDetail)
	if oldDetail.ID == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -10,
			"message": "?????????????????????",
		})
	}

	// ??????
	validateError := utils.Validate(detail)
	if validateError != nil {
		return ctx.JSON(fiber.Map{
			"status":  -3,
			"message": validateError.Error(),
		})
	}

	if *detail.Type == uint8(3) && detail.RepayAccountId == 0 {
		// ??????????????????,
		return ctx.JSON(fiber.Map{
			"status":  -4,
			"message": "?????????????????????",
		})
	}

	if *detail.Type == uint8(4) && detail.TransferAccountId == 0 {
		// ??????????????????,
		return ctx.JSON(fiber.Map{
			"status":  -5,
			"message": "?????????????????????",
		})
	}

	// ??????????????????????????????????????????????????????
	account := new(models.Account)
	db.Where("id =?", detail.AccountId).First(account)
	if account.ID == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -6,
			"message": "???????????????",
		})
	}
	if *detail.Type == uint8(2) && account.HasCredit == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -7,
			"message": "????????????????????????",
		})
	}

	if *detail.Type == uint8(3) {
		repayAccount := new(models.Account)
		db.Where("id =?", detail.RepayAccountId).First(repayAccount)
		if repayAccount.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "?????????????????????",
			})
		}
	}

	if *detail.Type == uint8(4) {
		transferAccount := new(models.Account)
		db.Where("id =?", detail.RepayAccountId).First(transferAccount)
		if transferAccount.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "?????????????????????",
			})
		}
	}
	if detail.CategoryId > 0 {
		cate := new(models.Category)
		db.Where("id =?", detail.CategoryId).Where("type =?", detail.Type).First(cate)
		if cate.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -9,
				"message": "?????????????????????",
			})
		}
	}

	money := detail.Money * 1000

	newDetail, err := charge_detail.Edit(
		detail.Id,
		detail.AccountId,
		*detail.Type,
		detail.CategoryId,
		int64(money),
		detail.Description,
		detail.RepayAccountId,
		detail.TransferAccountId,
		detail.CreateAt,
	)

	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -2,
			"message": err.Error(),
		})
	}

	switch oldDetail.Type {
	case 0:
		// ?????? ???????????? cashIn
		accountService.SummaryMoney("cashIn", oldDetail.AccountId, time.Unix(oldDetail.CreateAt, 0))
	case 1:
		// ?????? ???????????? cashOut
		accountService.SummaryMoney("cashOut", oldDetail.AccountId, time.Unix(oldDetail.CreateAt, 0))
	case 2:
		// ??? ???????????? creditOut
		accountService.SummaryMoney("creditOut", oldDetail.AccountId, time.Unix(oldDetail.CreateAt, 0))
	case 3:
		// ??? ???????????? cashOut ?????? ???????????? creditIn
		accountService.SummaryMoney("cashOut", oldDetail.AccountId, time.Unix(oldDetail.CreateAt, 0))
		accountService.SummaryMoney("creditIn", oldDetail.RepayAccountId, time.Unix(oldDetail.CreateAt, 0))
		// ???????????????????????????id
		charge_detail.ClearRepay(oldDetail.ID)
	case 4:
		// ??? ???????????? transferOut  ?????? ???????????? transferIn
		accountService.SummaryMoney("transferOut", oldDetail.AccountId, time.Unix(oldDetail.CreateAt, 0))
		accountService.SummaryMoney("transferIn", oldDetail.TransferAccountId, time.Unix(oldDetail.CreateAt, 0))
	}

	switch *detail.Type {
	case 0:
		// ?????? ???????????? cashIn
		accountService.SummaryMoney("cashIn", detail.AccountId, time.Unix(detail.CreateAt, 0))
	case 1:
		// ?????? ???????????? cashOut
		accountService.SummaryMoney("cashOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
	case 2:
		// ??? ???????????? creditOut
		accountService.SummaryMoney("creditOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
	case 3:
		// ??? ???????????? cashOut ?????? ???????????? creditIn
		accountService.SummaryMoney("cashOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
		accountService.SummaryMoney("creditIn", detail.RepayAccountId, time.Unix(detail.CreateAt, 0))
		// ???????????????????????????id
		charge_detail.UpdateRepay(newDetail.ID, detail.RepayDetailIds)
	case 4:
		// ??? ???????????? transferOut  ?????? ???????????? transferIn
		accountService.SummaryMoney("transferOut", detail.AccountId, time.Unix(detail.CreateAt, 0))
		accountService.SummaryMoney("transferIn", detail.TransferAccountId, time.Unix(detail.CreateAt, 0))
	}

	return ctx.JSON(fiber.Map{
		"status":  0,
		"data":    newDetail,
		"message": "????????????",
	})
}

func UnRepayList(ctx *fiber.Ctx) error {
	accountId, _ := ctx.ParamsInt("accountId")
	unPayList := charge_detail.GetUnPaidList(uint(accountId))

	return ctx.JSON(fiber.Map{
		"status": 0,
		"data":   unPayList,
	})
}
