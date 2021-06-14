package service

import (
	"errors"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
	"gin-vue-admin/utils"
	"strconv"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateBeeNodes
//@description: 创建BeeNodes记录
//@param: beeNodes model.BeeNodes
//@return: err error

func CreateBeeNodes(beeNodes model.BeeNodes) (err error) {
	err = global.GVA_DB.Create(&beeNodes).Error
	return err
}

func CreateBeeNodesInBatch(beeNodesList []model.BeeNodes) (err error) {
	err = global.GVA_DB.CreateInBatches(&beeNodesList, 1000).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteBeeNodes
//@description: 删除BeeNodes记录
//@param: beeNodes model.BeeNodes
//@return: err error

func DeleteBeeNodes(beeNodes model.BeeNodes) (err error) {
	err = global.GVA_DB.Delete(&beeNodes).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteBeeNodesByIds
//@description: 批量删除BeeNodes记录
//@param: ids request.IdsReq
//@return: err error

func DeleteBeeNodesByIds(ids request.IdsReq, jwtId uint) (err error) {
	if jwtId == 1 {
		err = global.GVA_DB.Delete(&[]model.BeeNodes{}, "id in ?", ids.Ids).Error
	} else {
		err = global.GVA_DB.Delete(&[]model.BeeNodes{}, "id in ? and user_id=?", ids.Ids, jwtId).Error
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateBeeNodes
//@description: 更新BeeNodes记录
//@param: beeNodes *model.BeeNodes
//@return: err error

func UpdateBeeNodes(beeNodes model.BeeNodes, jwtId uint) (err error) {
	if jwtId == 1 {
		err = global.GVA_DB.Save(&beeNodes).Error
	} else {
		var tmpBeeNodes model.BeeNodes
		err = global.GVA_DB.Where("id = ? and user_id = ?", beeNodes.ID, jwtId).First(&tmpBeeNodes).Error
		if err != nil {
			return err
		}

		if tmpBeeNodes.ID == beeNodes.ID {
			err = global.GVA_DB.Save(&beeNodes).Error
		}
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetBeeNodes
//@description: 根据id获取BeeNodes记录
//@param: id uint
//@return: err error, beeNodes model.BeeNodes

func GetBeeNodes(id uint) (err error, beeNodes model.BeeNodes) {
	err = global.GVA_DB.Where("id = ?", id).First(&beeNodes).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetBeeNodesInfoList
//@description: 分页获取BeeNodes记录
//@param: info request.BeeNodesSearch
//@return: err error, list interface{}, total int64

func GetBeeNodesInfoList(info request.BeeNodesSearch, jwtId uint) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.BeeNodes{})
	var beeNodess []model.BeeNodes
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Count(&total).Error
	if jwtId == 1 {
		err = db.Limit(limit).Offset(offset).Find(&beeNodess).Error
	} else {
		err = db.Where("user_id = ?", jwtId).Limit(limit).Offset(offset).Find(&beeNodess).Error
	}
	return err, beeNodess, total
}

func UpdateBeeNodeStatus(id uint, jwtId uint) (err error, beeNodes model.BeeNodes) {
	err, beeNodes = GetBeeNodes(id)
	if err != nil {
		return
	}
	if jwtId == 1 || beeNodes.ID == id {
		nodeState := utils.GetBeeNodeState(beeNodes.Ip, strconv.Itoa(beeNodes.DebugPort))
		beeNodes.BzzBalance = nodeState.BzzBalance
		beeNodes.EthBalance = nodeState.EthBalance
		beeNodes.PeerCount = nodeState.PeerCount
		beeNodes.UncashedCount = nodeState.CashCount
		beeNodes.UncashedAmount = nodeState.TotalUnCashed
		beeNodes.WalletAddress = nodeState.Address
		beeNodes.Version = nodeState.Version
		err = global.GVA_DB.Save(&beeNodes).Error
	} else {
		err = errors.New("Can not query this bee node")
	}
	return
}

func UpdateBeeNodeStatusInBatch(ids request.IdsReq, jwtId uint) (err error) {
	var beeNodess []model.BeeNodes
	beeNodess = make([]model.BeeNodes, len(ids.Ids))
	if jwtId == 1 {
		err = global.GVA_DB.Where("id in ?", ids.Ids).Find(&beeNodess).Error
	} else {
		err = global.GVA_DB.Where("id in ? and user_id = ?", ids.Ids, jwtId).Find(&beeNodess).Error
	}

	if err != nil {
		return
	}

	var ipPortList []request.IpPort
	ipPortList = make([]request.IpPort, len(beeNodess))
	for _, node := range beeNodess {
		ipPortList = append(ipPortList, request.IpPort{Ip: node.Ip, Port: node.DebugPort})
	}
	utils.GetBeeNodeStateInConcurrently(ipPortList)
	return err
}
