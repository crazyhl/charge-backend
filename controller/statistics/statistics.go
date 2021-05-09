package statistics

import (
	"charge/container"
	"charge/dto"
	"charge/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/now"
	"gorm.io/gorm/clause"
	"time"
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

func SummaryMonthData(ctx *fiber.Ctx) error {
	db := container.GetContainer().GetDb()
	month := ctx.Params("month")
	//
	totalAccountMonthDataDtoList := make([]dto.ChargeSummaryMonth, 0, 0)
	// 先获取总信息
	totalAccountMonthData := new(models.ChargeSummaryMonth)
	db.Model(&models.ChargeSummaryMonth{}).
		Select("date, year, sum(cash_in) as cash_in, sum(cash_out) as cash_out").
		Where("date = ?", month).
		First(totalAccountMonthData)
	totalAccountMonthDataDtoList = append(totalAccountMonthDataDtoList, dto.ChargeSummaryMonth{
		Account: dto.AccountDetail{ID: 0, Name: "全部"},
		Date:    totalAccountMonthData.Date,
		Year:    totalAccountMonthData.Year,
		CashIn:  float64(totalAccountMonthData.CashIn) / 1000,
		CashOut: float64(totalAccountMonthData.CashOut / 1000),
	})
	// 获取每月数据
	everyAccountMonthData := make([]models.ChargeSummaryMonth, 0, 0)
	db.Model(&models.ChargeSummaryMonth{}).
		Preload(clause.Associations).
		Select("account_id, date, year, cash_in, cash_out").
		Where("date = ?", month).
		Find(&everyAccountMonthData)
	fmt.Println("everyAccountMonthData", everyAccountMonthData)
	for _, monthData := range everyAccountMonthData {
		totalAccountMonthDataDtoList = append(totalAccountMonthDataDtoList, dto.ChargeSummaryMonth{
			Account: dto.AccountDetail{
				ID:   monthData.Account.ID,
				Name: monthData.Account.Name,
			},
			Date:    monthData.Date,
			Year:    monthData.Year,
			CashIn:  float64(monthData.CashIn) / 1000,
			CashOut: float64(monthData.CashOut) / 1000,
		})
	}

	return ctx.JSON(fiber.Map{
		"status": 0,
		"data":   totalAccountMonthDataDtoList,
	})
}

// ExpensesCategory 支出分类统计
func ExpensesCategory(ctx *fiber.Ctx) error {
	//db := container.GetContainer().GetDb()
	month := ctx.Params("month")
	location, _ := time.LoadLocation("Asia/Shanghai")
	monthTime, _ := time.Parse("20060102 -0700 MST", month+"01 +0800 CST")

	myConfig := &now.Config{
		TimeLocation: location,
	}

	beginningOfMonth := myConfig.With(monthTime).BeginningOfMonth()
	endOfMonth := myConfig.With(monthTime).EndOfMonth()

	fmt.Println(month + "01")
	fmt.Println(beginningOfMonth)
	fmt.Println(endOfMonth)
	fmt.Println(beginningOfMonth.Unix())
	fmt.Println(endOfMonth.Unix())

	return ctx.JSON(fiber.Map{
		"status": 0,
		"data":   string(beginningOfMonth.Unix()) + " -- " + string(endOfMonth.Unix()),
	})
}
