package request

import "gin-vue-admin/model"

type BeeNodesSearch struct {
	model.BeeNodes
	PageInfo
}

type IpPort struct {
	Ip   string `json:"ip" form:"ip"`
	Port int    `json:"port" form:"port"`
}

type IpPortListReq struct {
	IpPortList []IpPort `json:"ipPortList" form:"ipPortList"`
}

type CashOutInBatchReq struct {
	CashoutList []CashOutReq `json:"cashoutList" form:"cashoutList"`
}

type CashOutReq struct {
	Id       int    `json:"id" form:"id"`
	Nonce    int64  `json:"nonce" form:"nonce"`
	Count    int    `json:"count" form:"count"`
	GasPrice string `json:"gasPrice" form:"gasPrice"`
}
