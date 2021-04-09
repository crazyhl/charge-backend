package account

import (
	"charge/container"
	"charge/models"
)

// AddAccount 增加账户
func AddAccount(name string, opts ...AccountOptions) (uint, error) {
	account := new(models.Account)
	account.Name = name

	for _, o := range opts {
		o(account)
	}

	db := container.GetContainer().GetDb()
	result := db.Create(account)

	return account.ID, result.Error
}

// ------------ 上面各种方法用的 with 函数 -------------------

type AccountOptions func(account *models.Account)

func WithHasCredit(hasCredit uint8) AccountOptions {
	return func(account *models.Account) {
		account.HasCredit = hasCredit
	}
}

func WithCash(cash int64) AccountOptions {
	return func(account *models.Account) {
		account.Cash = cash
	}
}

func WithCredit(credit int64) AccountOptions {
	return func(account *models.Account) {
		if credit > 0 {
			credit = credit * -1
		}
		account.Credit = credit
	}
}

func WithSort(sort uint8) AccountOptions {
	return func(account *models.Account) {
		account.Sort = sort
	}
}

func WithChangeAt(time int) AccountOptions {
	return func(account *models.Account) {
		account.ChangeAt = time
	}
}
