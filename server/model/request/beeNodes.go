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
