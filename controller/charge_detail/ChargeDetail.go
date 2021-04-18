package category

import (
	"charge/container"
	"charge/models"
	"charge/services/category"
	"charge/services/charge_detail"
	"charge/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type addDetail struct {
	AccountId         uint    `json:"account_id" form:"account_id" validate:"required,number"`
	Type              *uint8  `json:"type" form:"type"  validate:"required,number,gte=0,lte=4"`
	CategoryId        uint    `json:"category_id" form:"category_id" validate:"required,number"`
	Money             float64 `json:"money" form:"money" validate:"required,number"`
	Description       string  `json:"description" form:"description"`
	Repay             bool    `json:"repay"  form:"description" validate:"required"`
	RepayDetailIds    []uint  `json:"repay_detail_ids" form:"repay_detail_ids"`
	RepayAccountId    uint    `json:"repay_account_id" form:"repay_account_id"`
	TransferAccountId uint    `json:"transfer_account_id" form:"transfer_account_id"`
	RepayAt           int64   `json:"repay_at" form:"repay_at"`
}

type editDetail struct {
	Id                uint    `json:"id" form:"id" validate:"required"  validate:"required"`
	AccountId         uint    `json:"account_id" form:"account_id" validate:"required,number"`
	Type              *uint8  `json:"type" form:"type"  validate:"required,number,gte=0,lte=4"`
	CategoryId        uint    `json:"category_id" form:"category_id" validate:"required,number"`
	Money             float64 `json:"money" form:"money" validate:"required,number"`
	Description       string  `json:"description" form:"description"`
	Repay             bool    `json:"repay"  form:"description" validate:"required"`
	RepayDetailIds    []uint  `json:"repay_detail_ids" form:"repay_detail_ids"`
	RepayAccountId    uint    `json:"repay_account_id" form:"repay_account_id"`
	Transfer          bool    `json:"transfer" form:"transfer"`
	TransferAccountId uint    `json:"transfer_account_id" form:"transfer_account_id"`
	RepayAt           int64   `json:"repay_at" form:"repay_at"`
}

//func List(ctx *fiber.Ctx) error {
//	listGroup := charge_detail.List()
//	return ctx.JSON(fiber.Map{
//		"status": 0,
//		"data":   listGroup,
//	})
//}

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
		db.Where("id =?", detail.RepayAccountId).First(transferAccount)
		if transferAccount.ID == 0 {
			return ctx.JSON(fiber.Map{
				"status":  -8,
				"message": "转账账户不存在",
			})
		}
	}
	cate := new(models.Category)
	db.Where("id =?", detail.CategoryId).Where("type =?", detail.Type).First(cate)
	if cate.ID == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -8,
			"message": "所选分类不存在",
		})
	}

	money := detail.Money * 1000
	transfer := 0
	if *detail.Type == 4 {
		transfer = 1
	}

	charge_detail.Add(
		detail.AccountId,
		*detail.Type,
		detail.CategoryId,
		int64(money),
		detail.Description,
		detail.RepayAt,
		detail.RepayAccountId,
		uint8(transfer),
		detail.TransferAccountId,
	)

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
	_, err = charge_detail.Delete(uint(id))
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -2,
			"message": err.Error(),
		})
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
	cate := new(models.Category)
	db.Where("id =?", detail.CategoryId).Where("type =?", detail.Type).First(cate)
	if cate.ID == 0 {
		return ctx.JSON(fiber.Map{
			"status":  -8,
			"message": "所选分类不存在",
		})
	}

	money := detail.Money * 1000
	transfer := 0
	if *detail.Type == 4 {
		transfer = 1
	}

	newAccount, err := charge_detail.Edit(
		detail.Id,
		detail.AccountId,
		*detail.Type,
		detail.CategoryId,
		int64(money),
		detail.Description,
		detail.RepayAt,
		detail.RepayAccountId,
		uint8(transfer),
		detail.TransferAccountId,
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
