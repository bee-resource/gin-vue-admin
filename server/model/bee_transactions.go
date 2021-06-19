// 自动生成模板BeeNodes
package model

import (
	"gin-vue-admin/global"
)

// 如果含有time.Time 请自行import time包
type BeeTransactions struct {
	global.GVA_MODEL
	BeeNodeId    uint    `json:"beeNodeId" form:"beeNodeId" gorm:"index,column:bee_node_id;comment:;type:bigint;size:19;"`
	Peer         string  `json:"per" form:"peer" gorm:"column:peer;comment:;type:varchar(90);size:90;"`
	CashedAmount float64 `json:"cashedAmount" form:"cashedAmount" gorm:"column:cashed_amount;comment:;type:double"`
	GasPrice     float64 `json:"gasPrice" form:"gasPrice" gorm:"column:gas_price;comment:;type:double"`
	Nonce        int64   `json:"nonce" form:"nonce" gorm:"column:nonce;comment:;type:bigint;size:19;"`
	TxId         string  `json:"txid" form:"txid" gorm:"column:txid;comment:;type:varchar(191);size:191;"`
}

func (BeeTransactions) TableName() string {
	return "bee_transactions"
}
