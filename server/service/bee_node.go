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
	err = db.Where("user_id = ?", jwtId).Count(&total).Error
	if jwtId == 888 {
		err = db.Limit(limit).Offset(offset).Find(&beeNodess).Error
	} else {
		err = db.Where("user_id = ?", jwtId).Limit(limit).Offset(offset).Find(&beeNodess).Error
	}
	return err, beeNodess, total
}

type StatisticResult struct {
	TotalNum            int     `gorm:"column:totalNum"`
	TotalUnCashedCount  int     `gorm:"column:totalUnCashedCount"`
	TotalUnCashedAmount float64 `gorm:"column:totalUnCashedAmount"`
}

func StatisticInfo(jwtId uint) (err error, info interface{}) {
	// 创建db
	var result StatisticResult
	var beeNodes []model.BeeNodes
	db := global.GVA_DB.Model(&model.BeeNodes{})
	if jwtId == 888 {
		err = db.Find(&beeNodes).Error
	} else {
		err = db.Select("count(id) as totalNum, sum(uncashed_count) as totalUnCashedCount, sum(uncashed_amount) as totalUnCashedAmount").Where("user_id = ?", jwtId).Find(&result).Error
	}
	return err, result
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
	beeNodess = make([]model.BeeNodes, 0)
	if jwtId == 1 {
		err = global.GVA_DB.Where("id in ?", ids.Ids).Find(&beeNodess).Error
	} else {
		err = global.GVA_DB.Where("id in ? and user_id = ?", ids.Ids, jwtId).Find(&beeNodess).Error
	}

	if err != nil {
		return
	}
	utils.GetBeeNodeStateInConcurrently(beeNodess)
	err = global.GVA_DB.Save(&beeNodess).Error
	return err
}

func CashoutBeeNodesInBatch(cashoutBeeNodesInBatchReq request.CashOutInBatchReq, jwtId uint) (err error, transactions []model.BeeTransactions) {
	var beeNodess []model.BeeNodes
	var ids []int
	ids = make([]int, 0)
	beeNodess = make([]model.BeeNodes, 0)

	for _, cashoutReq := range cashoutBeeNodesInBatchReq.CashoutList {
		ids = append(ids, cashoutReq.Id)
	}
	if jwtId == 1 {
		err = global.GVA_DB.Where("id in ?", ids).Find(&beeNodess).Error
	} else {
		err = global.GVA_DB.Where("id in ? and user_id = ?", ids, jwtId).Find(&beeNodess).Error
	}
	if err != nil {
		return
	}

	ipPeerCashOutInfoMap := utils.CashoutBeeNodesInConcurrently(cashoutBeeNodesInBatchReq, beeNodess)
	err = global.GVA_DB.Save(&beeNodess).Error
	if err != nil {
		return
	}
	transactions = make([]model.BeeTransactions, 0)
	// save transactions
	for _, peerCashOutInfoMap := range ipPeerCashOutInfoMap {
		for peer, cashOutInfo := range peerCashOutInfoMap {
			transactions = append(transactions, model.BeeTransactions{BeeNodeId: cashOutInfo.NodeId,
				Peer: peer, CashedAmount: float64(cashOutInfo.Amount),
				GasPrice: cashOutInfo.GasPrice, Nonce: cashOutInfo.Nonce,
				TxId: cashOutInfo.TxId,
			})
		}
	}
	err = global.GVA_DB.Save(&transactions).Error
	return
}