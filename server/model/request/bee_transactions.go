package request

import "gin-vue-admin/model"

type BeeTransactionsSearch struct {
	model.BeeTransactions
	PageInfo
}
