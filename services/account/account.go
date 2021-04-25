package account

import (
	"charge/container"
	"charge/dto"
	"charge/models"
	"errors"
	"gorm.io/gorm"
	"time"
)

// List 账户列表，会返回所有的账户
func List() []dto.AccountDetail {
	var accounts []dto.AccountDetail
	var accountRows []models.Account
	db := container.GetContainer().GetDb()
	db.Order("sort DESC").Find(&accountRows)

	for _, acc := range accountRows {
		createTm := time.Unix(acc.CreateAt, 0)
		updateTm := time.Unix(acc.UpdateAt, 0)
		changeTm := time.Unix(acc.ChangeAt, 0)
		accounts = append(accounts, dto.AccountDetail{
			ID:        acc.ID,
			Name:      acc.Name,
			HasCredit: acc.HasCredit == 1,
			Cash:      float64(acc.Cash) / 1000.0,
			Credit:    float64(acc.Credit) / 1000.0,
			Sort:      acc.Sort,
			CreateAt:  createTm.Format("2006-01-02 15:04:05"),
			UpdateAt:  updateTm.Format("2006-01-02 15:04:05"),
			ChangeAt:  changeTm.Format("2006-01-02 15:04:05"),
		})
	}
	return accounts
}

// Add 增加账户
func Add(name string, opts ...Options) (*models.Account, error) {
	account := new(models.Account)
	account.Name = name

	for _, o := range opts {
		o(account)
	}

	db := container.GetContainer().GetDb()
	result := db.Create(account)

	return account, result.Error
}

func Delete(id uint) (*models.Account, error) {
	db := container.GetContainer().GetDb()
	account := new(models.Account)
	db.Where("id = ?", id).First(account)
	if account.ID == 0 {
		return account, errors.New("删除对象不存在")
	}

	result := db.Delete(account)

	return account, result.Error
}

func EditDetail(id int) (*dto.AccountEditDetail, error) {
	db := container.GetContainer().GetDb()
	accountDto := new(dto.AccountEditDetail)
	account := new(models.Account)
	result := db.Where("id = ?", id).First(account)
	if result.Error != nil {
		return accountDto, result.Error
	}

	accountDto.ID = account.ID
	accountDto.Name = account.Name
	accountDto.HasCredit = account.HasCredit == 1
	accountDto.Cash = float64(account.Cash) / 1000.0
	accountDto.Credit = float64(account.Credit) / 1000.0
	accountDto.Sort = account.Sort

	return accountDto, nil
}

// Add 增加账户
func Edit(id int, opts ...Options) (*models.Account, error) {
	db := container.GetContainer().GetDb()
	account := new(models.Account)
	db.Where("id = ?", id).First(account)
	if account.ID == 0 {
		return account, errors.New("更新对象不存在")
	}

	for _, o := range opts {
		o(account)
	}

	result := db.Save(account)

	return account, result.Error
}

func IncreaseCash(accountId uint, money int64) error {
	return changeMoney("cash", accountId, money, "+")
}

func IncreaseCredit(accountId uint, money int64) error {
	return changeMoney("credit", accountId, money, "+")
}

func DecreaseCash(accountId uint, money int64) error {
	return changeMoney("cash", accountId, money, "-")
}

func DecreaseCredit(accountId uint, money int64) error {
	return changeMoney("credit", accountId, money, "-")
}

func changeMoney(fieldName string, accountId uint, money int64, operate string) error {
	db := container.GetContainer().GetDb()

	account := new(models.Account)
	db.Where("id = ?", accountId).First(account)
	if account.ID == 0 {
		return errors.New("账户不存在")
	}

	result := db.Model(account).Updates(
		map[string]interface{}{
			fieldName:  gorm.Expr(fieldName+" "+operate+" ?", money),
			"ChangeAt": time.Now().Unix(),
		},
	)

	return result.Error
}

// ------------ 上面各种方法用的 with 函数 -------------------

type Options func(account *models.Account)

func WithName(name string) Options {
	return func(account *models.Account) {
		account.Name = name
	}
}

func WithHasCredit(hasCredit uint8) Options {
	return func(account *models.Account) {
		account.HasCredit = hasCredit
	}
}

func WithCash(cash int64) Options {
	return func(account *models.Account) {
		account.Cash = cash
	}
}

func WithCredit(credit int64) Options {
	return func(account *models.Account) {
		if credit > 0 {
			credit = credit * -1
		}
		account.Credit = credit
	}
}

func WithSort(sort uint8) Options {
	return func(account *models.Account) {
		account.Sort = sort
	}
}

func WithChangeAt(time int64) Options {
	return func(account *models.Account) {
		account.ChangeAt = time
	}
}
