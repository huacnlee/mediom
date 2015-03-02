package models

import (
	"github.com/revel/revel"
	"time"
)

type Node struct {
	BaseModel
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
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
