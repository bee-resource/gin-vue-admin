package service

import (
	"fmt"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
	"time"
)

type TransactionResult struct {
	CreatedAt     time.Time // 创建时间
	Ip            string    `json:"ip"`
	Peer          string    `json:"peer"`
	WalletAddress string    `json:"walletAddress"`
	Txid          string    `json:"txid"`
	GasPrice      string    `json:"gasPrice"`
	CashedAmount  float64   `json:"cashedAmount"`
	Nonce         int64     `json:"nonce"`
}

func GetBeeTransactionList(info request.BeeTransactionsSearch, jwtId uint) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&model.BeeTransactions{})
	var beeTransactions []TransactionResult
	// 如果有条件搜索 下方会自动创建搜索语句
	err = db.Debug().Select("bee_transactions.created_at, bee_nodes.ip, bee_nodes.wallet_address, bee_transactions.peer, bee_transactions.txid, bee_transactions.gas_price, bee_transactions.cashed_amount, bee_transactions.nonce").Joins("left join bee_nodes on bee_nodes.ID = bee_transactions.bee_node_id").Where("bee_nodes.user_id = ?", jwtId).Count(&total).Error
	if jwtId == 1 {
		err = db.Limit(limit).Offset(offset).Find(&beeTransactions).Error
	} else {
		err = db.Debug().Select("bee_transactions.created_at, bee_nodes.ip, bee_nodes.wallet_address, bee_transactions.peer, bee_transactions.txid, bee_transactions.gas_price, bee_transactions.cashed_amount, bee_transactions.nonce").Joins("left join bee_nodes on bee_nodes.ID = bee_transactions.bee_node_id").Where("bee_nodes.user_id = ?", jwtId).Order("bee_transactions.created_at desc").Limit(limit).Offset(offset).Scan(&beeTransactions).Error
	}
	fmt.Printf("%v\n", beeTransactions)
	return err, beeTransactions, total
}
