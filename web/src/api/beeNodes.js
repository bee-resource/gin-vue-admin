import service from "@/utils/request";

// @Tags BeeNodes
// @Summary 创建BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BeeNodes true "创建BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /beeNodes/createBeeNodes [post]
export const createBeeNodes = (data) => {
  return service({
    url: "/beeNodes/createBeeNodes",
    method: "post",
    data,
  });
};

// @Tags BeeNodes
// @Summary 删除BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BeeNodes true "删除BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /beeNodes/deleteBeeNodes [delete]
export const deleteBeeNodes = (data) => {
  return service({
    url: "/beeNodes/deleteBeeNodes",
    method: "delete",
    data,
  });
};

// @Tags BeeNodes
// @Summary 删除BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /beeNodes/deleteBeeNodes [delete]
export const deleteBeeNodesByIds = (data) => {
  return service({
    url: "/beeNodes/deleteBeeNodesByIds",
    method: "delete",
    data,
  });
};

// @Tags BeeNodes
// @Summary 更新BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BeeNodes true "更新BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /beeNodes/updateBeeNodes [put]
export const updateBeeNodes = (data) => {
  return service({
    url: "/beeNodes/updateBeeNodes",
    method: "put",
    data,
  });
};

// @Tags BeeNodes
// @Summary 用id查询BeeNodes
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BeeNodes true "用id查询BeeNodes"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /beeNodes/findBeeNodes [get]
export const findBeeNodes = (params) => {
  return service({
    url: "/beeNodes/findBeeNodes",
    method: "get",
    params,
  });
};

// @Tags BeeNodes
// @Summary 分页获取BeeNodes列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取BeeNodes列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /beeNodes/getBeeNodesList [get]
export const getBeeNodesList = (params) => {
  return service({
    url: "/beeNodes/getBeeNodesList",
    method: "get",
    params,
  });
};

export const updateBeeNodesStatus = (data) => {
  return service({
    url: "/beeNodes/updateBeeNodeStatus",
    method: "put",
    data,
  });
};

export const importBeeNodes = (data) => {
  return service({
    url: "/beeNodes/importBeeNodes",
    method: "post",
    data,
  });
};

export const updateBeeNodeStatusInBatch = (data) => {
  return service({
    url: "/beeNodes/updateBeeNodeStatusInBatch",
    method: "put",
    data,
  });
};
