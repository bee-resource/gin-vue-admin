package router

import (
	"gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
)

func InitBeeNodesRouter(Router *gin.RouterGroup) {
	BeeNodesRouter := Router.Group("beeNodes").Use(middleware.OperationRecord())
	{
		BeeNodesRouter.POST("createBeeNodes", v1.CreateBeeNodes)   // 新建BeeNodes
		BeeNodesRouter.DELETE("deleteBeeNodes", v1.DeleteBeeNodes) // 删除BeeNodes
		BeeNodesRouter.DELETE("deleteBeeNodesByIds", v1.DeleteBeeNodesByIds) // 批量删除BeeNodes
		BeeNodesRouter.PUT("updateBeeNodes", v1.UpdateBeeNodes)    // 更新BeeNodes
		BeeNodesRouter.GET("findBeeNodes", v1.FindBeeNodes)        // 根据ID获取BeeNodes
		BeeNodesRouter.GET("getBeeNodesList", v1.GetBeeNodesList)  // 获取BeeNodes列表
	}
}
