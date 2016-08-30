package models

import (
	"github.com/revel/revel"
	"time"
)

type Node struct {
	BaseModel
	Name        string `sql:"not null"`
	Summary     string `sql:"type:text"`
	ParentId    *int
	Sort        int `sql:"default: 0; not null"`
	Children    []Node `gorm:"ForeignKey:ParentId"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

func FindAllNodeRoots() (roots []*Node) {
	db.Preload("Children").Order("sort desc").Where("parent_id is null or parent_id = 0").Find(&roots)
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
