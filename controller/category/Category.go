package category

import (
	"charge/services/category"
	"charge/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type addCategory struct {
	Name string   `json:"name" form:"name" validate:"required"`
	Sort uint8    `json:"sort" form:"sort"`
	Type []*uint8 `json:"type" form:"type"  validate:"required,dive,required,number,gte=0,lte=4"`
}

type editCategory struct {
	Id   int    `json:"id" form:"id" validate:"required"  validate:"required"`
	Name string `json:"name" form:"name" validate:"required"`
	Sort uint8  `json:"sort" form:"sort"`
	Type *uint8 `json:"type" form:"type"  validate:"required,number,gte=0,lte=4"`
}

func List(ctx *fiber.Ctx) error {
	listGroup := category.ListGroup()
	return ctx.JSON(fiber.Map{
		"status": 0,
		"data":   listGroup,
	})
}

func Add(ctx *fiber.Ctx) error {
	addCategory := new(addCategory)
	if err := ctx.BodyParser(addCategory); err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -1,
			"message": err.Error(),
		})
	}
	// 验证
	validateError := utils.Validate(addCategory)
	if validateError != nil {
		fmt.Println(validateError)
		return ctx.JSON(fiber.Map{
			"status":  -3,
			"message": validateError.Error(),
		})
	}

	for _, _type := range addCategory.Type {
		_, err := category.Add(*_type, addCategory.Name, addCategory.Sort)
		if err != nil {
			return ctx.JSON(fiber.Map{
				"status":  -4,
				"message": err.Error(),
			})
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
	_, err = category.Delete(uint(id))
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
	editCat := new(editCategory)
	if err := ctx.BodyParser(editCat); err != nil {
		return ctx.JSON(fiber.Map{
			"status":  -1,
			"message": err.Error(),
		})
	}
	// 验证
	validateError := utils.Validate(editCat)
	if validateError != nil {
		fmt.Println(validateError)
		return ctx.JSON(fiber.Map{
			"status":  -3,
			"message": validateError.Error(),
		})
	}

	newAccount, err := category.Edit(
		editCat.Id,
		*editCat.Type,
		editCat.Name,
		editCat.Sort,
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
