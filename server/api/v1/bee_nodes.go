package v1

import (
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
	"gin-vue-admin/model/response"
	"gin-vue-admin/service"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

// @Tags BeeNodes
// @Summary 创建BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BeeNodes true "创建BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /beeNodes/createBeeNodes [post]
func CreateBeeNodes(c *gin.Context) {
	var beeNodes model.BeeNodes
	_ = c.ShouldBindJSON(&beeNodes)
	beeNodes.UserId = getUserID(c)
	beeNodes.Uuid = uuid.NewV4()
	if err := service.CreateBeeNodes(beeNodes); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags BeeNodes
// @Summary 删除BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BeeNodes true "删除BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /beeNodes/deleteBeeNodes [delete]
func DeleteBeeNodes(c *gin.Context) {
	var beeNodes model.BeeNodes
	_ = c.ShouldBindJSON(&beeNodes)
	jwtId := getUserID(c)
	if jwtId != uint(beeNodes.UserId) && jwtId != 1 {
		response.FailWithMessage("删除失败", c)
		return
	}
	if err := service.DeleteBeeNodes(beeNodes); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags BeeNodes
// @Summary 批量删除BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /beeNodes/deleteBeeNodesByIds [delete]
func DeleteBeeNodesByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	jwtId := getUserID(c)
	if err := service.DeleteBeeNodesByIds(IDS, jwtId); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags BeeNodes
// @Summary 更新BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BeeNodes true "更新BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /beeNodes/updateBeeNodes [put]
func UpdateBeeNodes(c *gin.Context) {
	var beeNodes model.BeeNodes
	_ = c.ShouldBindJSON(&beeNodes)
	jwtId := getUserID(c)

	if err := service.UpdateBeeNodes(beeNodes, jwtId); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags BeeNodes
// @Summary 用id查询BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BeeNodes true "用id查询BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /beeNodes/findBeeNodes [get]
func FindBeeNodes(c *gin.Context) {
	var beeNodes model.BeeNodes
	_ = c.ShouldBindQuery(&beeNodes)
	if err, rebeeNodes := service.GetBeeNodes(beeNodes.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rebeeNodes": rebeeNodes}, c)
	}
}

// @Tags BeeNodes
// @Summary 分页获取BeeNodes列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.BeeNodesSearch true "分页获取BeeNodes列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /beeNodes/getBeeNodesList [get]
func GetBeeNodesList(c *gin.Context) {
	var pageInfo request.BeeNodesSearch
	_ = c.ShouldBindQuery(&pageInfo)
	jwtId := getUserID(c)
	if err, list, total := service.GetBeeNodesInfoList(pageInfo, jwtId); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func UpdateBeeNodeStatus(c *gin.Context) {
	var beeNodes model.BeeNodes
	_ = c.ShouldBindJSON(&beeNodes)
	jwtId := getUserID(c)
	if err, rebeeNodes := service.UpdateBeeNodeStatus(beeNodes.ID, jwtId); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithData(gin.H{"rebeeNodes": rebeeNodes}, c)
	}
}

func ImportBeeNodes(c *gin.Context) {
	var ipPortListReq request.IpPortListReq
	_ = c.ShouldBindJSON(&ipPortListReq)

	var beeNodesList []model.BeeNodes
	userId := getUserID(c)
	batchSize := len(ipPortListReq.IpPortList)
	if batchSize == 0 {
		response.FailWithMessage("批量导入失败", c)
		return
	}
	beeNodesList = make([]model.BeeNodes, batchSize)
	for _, ipPort := range ipPortListReq.IpPortList {
		beeNodesList = append(beeNodesList, model.BeeNodes{UserId: userId, Uuid: uuid.NewV4(), Ip: strings.Trim(ipPort.Ip, " "), DebugPort: ipPort.Port})
	}
	if err := service.CreateBeeNodesInBatch(beeNodesList); err != nil {
		global.GVA_LOG.Error("批量导入失败!", zap.Any("err", err))
		response.FailWithMessage("批量导入失败", c)
	} else {
		response.OkWithMessage("批量导入成功", c)
	}
}

func UpdateBeeNodeStatusInBatch(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	jwtId := getUserID(c)
	if err := service.UpdateBeeNodeStatusInBatch(IDS, jwtId); err != nil {
		global.GVA_LOG.Error("批量更新失败!", zap.Any("err", err))
		response.FailWithMessage("批量更新失败", c)
	} else {
		response.OkWithMessage("批量更新成功", c)
	}
}

func CashoutBeeNodesInBatch(c *gin.Context) {
	var cashOutInBatchReq request.CashOutInBatchReq
	_ = c.ShouldBindJSON(&cashOutInBatchReq)
	jwtId := getUserID(c)
	if err, ipPeerCashOutInfoMap := service.CashoutBeeNodesInBatch(cashOutInBatchReq, jwtId); err != nil {
		global.GVA_LOG.Error("批量取票失败!", zap.Any("err", err))
		response.FailWithMessage("批量取票失败", c)
	} else {
		response.OkWithData(gin.H{"ipPeerCashOutInfoMap": ipPeerCashOutInfoMap}, c)
	}
}
