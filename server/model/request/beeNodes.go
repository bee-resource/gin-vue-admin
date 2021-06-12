package request

import "gin-vue-admin/model"

type BeeNodesSearch struct{
    model.BeeNodes
    PageInfo
}