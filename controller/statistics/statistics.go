package statistics

import (
	"charge/container"
	"github.com/gofiber/fiber/v2"
)

// SummaryMonthList 获取月汇总的月份信息列表
func SummaryMonthList(ctx *fiber.Ctx) error {
	db := container.GetContainer().GetDb()
	monthList := make([]string, 0, 0)
	db.Table("charge_summary_months").
		Distinct("date").
		Order("date DESC").
		Select("date").Find(&monthList)
	return ctx.JSON(fiber.Map{
		"status": 0,
		"data":   monthList,
	})
}
