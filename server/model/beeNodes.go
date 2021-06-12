// 自动生成模板BeeNodes
package model

import (
	"gin-vue-admin/global"

	uuid "github.com/satori/go.uuid"
)

// 如果含有time.Time 请自行import time包
type BeeNodes struct {
	global.GVA_MODEL
	Uuid          uuid.UUID `json:"uuid" form:"uuid" gorm:"column:uuid;comment:;type:varchar(191);size:191;"`
	Name          string    `json:"name" form:"name" gorm:"column:name;comment:;type:varchar(191);size:191;"`
	Ip            string    `json:"ip" form:"ip" gorm:"column:ip;comment:;type:varchar(90);size:90;"`
	InternalIp    string    `json:"internalIp" form:"ip" gorm:"column:ip;comment:;type:varchar(90);size:90;"`
	DebugPort     int       `json:"debugPort" form:"debugPort" gorm:"column:debug_port;comment:;type:bigint;size:19;"`
	WalletAddress string    `json:"walletAddress" form:"walletAddress" gorm:"column:wallet_address;comment:;type:varchar(191);size:191;"`
	UncashedCount int       `json:"uncashedCount" form:"uncashedCount" gorm:"column:uncashed_count;comment:;type:bigint;size:19;"`
	PeerCount     int       `json:"peerCount" form:"peerCount" gorm:"column:peer_count;comment:;type:bigint;size:19;"`
	EthBalance    float64   `json:"ethBalance" form:"ethBalance" gorm:"column:eth_balance;comment:;type:decimal;"`
	BzzBalance    float64   `json:"bzzBalance" form:"bzzBalance" gorm:"column:bzz_balance;comment:;type:decimal;"`
	UserId        uint      `json:"userId" form:"userId" gorm:"column:user_id;comment:;type:bigint;size:19;"`
	Status        string    `json:"status" form:"status" gorm:"column:status;comment:;type:varchar(191);size:20;"`
	Version       string    `json:"version" form:"version" gorm:"column:version;comment:;type:varchar(191);size:20;"`
}

func (BeeNodes) TableName() string {
	return "bee_nodes"
}
