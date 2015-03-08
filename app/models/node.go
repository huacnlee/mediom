package models

import (
	"github.com/huacnlee/revel"
	"time"
)

type Node struct {
	BaseModel
	Name        string `sql:"not null"`
	Summary     string `sql:"type:text"`
	NodeGroupId int
	Sort        int `sql:"default: 0; not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NodeGroup struct {
	Id    int32
	Name  string
	Sort  int `sql:"default: 0; not null"`
	Nodes []Node
}

func (n *Node) validate() (v revel.Validation) {
	v = revel.Validation{}
	switch n.NewRecord() {
	case false:
	default:
		v.Required(n.Name).Key("name").Message("不能为空")
	}
	return v
}

func CreateNode(n *Node) revel.Validation {
	v := n.validate()
	if v.HasErrors() {
		return v
	}

	err := db.Save(n).Error
	if err != nil {
		v.Error("服务器异常创建失败")
	}
	return v
}

func UpdateNode(n *Node) revel.Validation {
	v := n.validate()
	if v.HasErrors() {
		return v
	}

	err := db.Save(&n).Error
	if err != nil {
		v.Error("服务器异常更新失败")
	}
	return v
}

func FindAllNodeGroups() (groups []*NodeGroup) {
	db.Preload("Nodes").Order("sort desc").Find(&groups)
	return
}

func FindAllNodes() (nodes []*Node) {
	db.Order("name asc").Find(&nodes)
	return
}

func FindNodesBySort(limit int) (nodes []*Node) {
	db.Order("sort desc, name asc").Limit(limit).Find(&nodes)
	return
}
