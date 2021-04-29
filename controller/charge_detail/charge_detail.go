package charge_detail

import (
	"charge/container"
	"charge/models"
	accountService "charge/services/account"
	"charge/services/category"
	"charge/services/charge_detail"
	"charge/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
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

	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(size)
	if pageSizeInt <= 0 {
		pageSizeInt = 20
	}
	if pageInt < 1 {
		pageInt = 1
	}
	pageStart := (pageInt - 1) * pageSizeInt
	listData := charge_detail.List(pageStart, pageSizeInt)

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
	// 验证
	validateError := utils.Validate(detail)
	if validateError != nil {
		fmt.Println(validateError)
		return ctx.JSON(fiber.Map{
			"status":  -3,
			"message": validateError.Error(),
		})
	}

	if *detail.Type == uint8(3) && detail.RepayAccountId == 0 {
		// 是还款的时候,
		return ctx.JSON(fiber.Map{
			"status":  -4,
			"message": "请选择还款账户",
		})
	}

	if *detail.Type == uint8(4) && detail.TransferAccountId == 0 {
		// 是转账的时候,
		return ctx.JSON(fiber.Map{
			"status":  -5,
			"message": "请选择转账账户",
		})
	}

	// 如果是借款需要判定改账户是否支持借款
	account := new(models.Account)
	db.Where("id =?", detail.AccountId).First(account)
	if account.ID == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -6,
			"message": "账户不存在",
		})
	}
	if *detail.Type == uint8(2) && account.HasCredit == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -7,
			"message": "该账户不支持借款",
		})
	}

	if *detail.Type == uint8(3) {
		repayAccount := new(models.Account)
		db.Where("id =?", detail.RepayAccountId).First(repayAccount)
		if repayAccount.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "还款账户不存在",
			})
		}
	}

	if *detail.Type == uint8(4) {
		transferAccount := new(models.Account)
		db.Where("id =?", detail.TransferAccountId).First(transferAccount)
		if transferAccount.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "转账账户不存在",
			})
		}
	}

	if detail.CategoryId > 0 {
		cate := new(models.Category)
		db.Where("id =?", detail.CategoryId).Where("type =?", detail.Type).First(cate)
		if cate.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "所选分类不存在",
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
		// 改钱采用按月统计合并的方式
		switch *detail.Type {
		case 0:
			// 收入 改变账户 cashIn
			accountService.SummaryMoney("cashIn", detail.AccountId, time.Now())
		case 1:
			// 支出 改变账户 cashOut
			accountService.SummaryMoney("cashOut", detail.AccountId, time.Now())
		case 2:
			// 借 改变账户 creditOut
			accountService.SummaryMoney("creditOut", detail.AccountId, time.Now())
		case 3:
			// 还 改变账户 cashOut 改变 还款账户 creditIn
			accountService.SummaryMoney("cashOut", detail.AccountId, time.Now())
			accountService.SummaryMoney("creditIn", detail.RepayAccountId, time.Now())
			// 更新借款账目的还款id
			charge_detail.UpdateRepay(chargeDetail.ID, detail.RepayDetailIds)
		case 4:
			// 转 改变账户 transferOut  改变 转账账户 transferIn
			accountService.SummaryMoney("transferOut", detail.AccountId, time.Now())
			accountService.SummaryMoney("transferIn", detail.TransferAccountId, time.Now())
		}
	}

	return ctx.JSON(fiber.Map{
		"status":  0,
		"message": "添加完成",
	})
}

// Delete 删除账户
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
		// 收入 改变账户 cashIn
		accountService.SummaryMoney("cashIn", detail.AccountId, time.Now())
	case 1:
		// 支出 改变账户 cashOut
		accountService.SummaryMoney("cashOut", detail.AccountId, time.Now())
	case 2:
		// 借 改变账户 creditOut
		accountService.SummaryMoney("creditOut", detail.AccountId, time.Now())
	case 3:
		// 还 改变账户 cashOut 改变 还款账户 creditIn
		accountService.SummaryMoney("cashOut", detail.AccountId, time.Now())
		accountService.SummaryMoney("creditIn", detail.RepayAccountId, time.Now())
		// 更新借款账目的还款id
		charge_detail.ClearRepay(detail.RepaidDetailId)
	case 4:
		// 转 改变账户 transferOut  改变 转账账户 transferIn
		accountService.SummaryMoney("transferOut", detail.AccountId, time.Now())
		accountService.SummaryMoney("transferIn", detail.TransferAccountId, time.Now())
	}

	return ctx.JSON(fiber.Map{
		"status":  0,
		"message": "删除成功",
	})
}

// EditDetail 编辑前获取详情
func EditDetail(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -1,
			"message": err.Error(),
		})
	}

	categoryDto, err := category.EditDetail(id)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -2,
			"message": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"status": 0,
		"data":   categoryDto,
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
	// 验证
	validateError := utils.Validate(detail)
	if validateError != nil {
		fmt.Println(validateError)
		return ctx.JSON(fiber.Map{
			"status":  -3,
			"message": validateError.Error(),
		})
	}

	if *detail.Type == uint8(3) && detail.RepayAccountId == 0 {
		// 是还款的时候,
		return ctx.JSON(fiber.Map{
			"status":  -4,
			"message": "请选择还款账户",
		})
	}

	if *detail.Type == uint8(4) && detail.TransferAccountId == 0 {
		// 是转账的时候,
		return ctx.JSON(fiber.Map{
			"status":  -5,
			"message": "请选择转账账户",
		})
	}

	// 如果是借款需要判定改账户是否支持借款
	account := new(models.Account)
	db.Where("id =?", detail.AccountId).First(account)
	if account.ID == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -6,
			"message": "账户不存在",
		})
	}
	if *detail.Type == uint8(2) && account.HasCredit == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -7,
			"message": "该账户不支持借款",
		})
	}

	if *detail.Type == uint8(3) {
		repayAccount := new(models.Account)
		db.Where("id =?", detail.RepayAccountId).First(repayAccount)
		if repayAccount.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "还款账户不存在",
			})
		}
	}

	if *detail.Type == uint8(4) {
		transferAccount := new(models.Account)
		db.Where("id =?", detail.RepayAccountId).First(transferAccount)
		if transferAccount.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "转账账户不存在",
			})
		}
	}
	if detail.CategoryId > 0 {
		cate := new(models.Category)
		db.Where("id =?", detail.CategoryId).Where("type =?", detail.Type).First(cate)
		if cate.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "所选分类不存在",
			})
		}
	}

	money := detail.Money * 1000

	newAccount, err := charge_detail.Edit(
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

	return ctx.JSON(fiber.Map{
		"status":  0,
		"data":    newAccount,
		"message": "修改成功",
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
