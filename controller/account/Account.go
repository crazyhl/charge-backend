package account

import (
	"charge/services/account"
	"charge/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type addAccount struct {
	Name      string  `json:"name" form:"name" validate:"required"`
	HasCredit bool    `json:"hasCredit" form:"hasCredit"`
	Cash      float64 `json:"cash" form:"cash"`
	Credit    float64 `json:"credit" form:"credit"`
	Sort      uint8   `json:"sort" form:"sort"`
}

func Add(ctx *fiber.Ctx) error {
	acc := new(addAccount)
	if err := ctx.BodyParser(acc); err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -1,
			"message": err.Error(),
		})
	}
	// 验证
	validateError := utils.Validate(acc)
	if validateError != nil {
		fmt.Println(validateError)
		return ctx.JSON(fiber.Map{
			"status":  -3,
			"message": validateError.Error(),
		})
	}

	cash := int64(acc.Cash * 1000)
	hasCredit := 0
	if acc.HasCredit {
		hasCredit = 1
	}
	credit := int64(acc.Credit * 1000)

	newAccount, err := account.AddAccount(
		acc.Name,
		account.WithCash(cash),
		account.WithHasCredit(uint8(hasCredit)),
		account.WithCredit(credit),
		account.WithSort(uint8(acc.Sort)),
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
		"message": "添加成功",
	})
}