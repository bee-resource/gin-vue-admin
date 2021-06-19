package router

import (
	v1 "gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"

	"github.com/gin-gonic/gin"
)

func InitBeeTransactionsRouter(Router *gin.RouterGroup) {
	BeeTransactionsRouter := Router.Group("beeTransactions").Use(middleware.OperationRecord())
	{
		BeeTransactionsRouter.GET("getBeeTransactionList", v1.GetBeeTransactionList) // 获取BeeTransactions列表
	}
}
