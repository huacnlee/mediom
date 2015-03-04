package controllers

import (
	"github.com/revel/revel"
	. "mediom/app/models"
)

type Nodes struct {
	App
}

func init() {
	revel.InterceptMethod((*Nodes).Before, revel.BEFORE)
}

func (c *Nodes) Before() revel.Result {
	if r := c.requireAdmin(); r != nil {
		return r
	}

	return nil
}

func (c Nodes) Index() revel.Result {
	nodes := FindAllNodes()
	c.RenderArgs["nodes"] = nodes
	return c.Render("nodes/index.html")
}

func (c Nodes) Create() revel.Result {
	n := Node{Name: c.Params.Get("name")}

	v := CreateNode(&n)
	if v.HasErrors() {
		c.RenderArgs["node"] = n
		return c.renderValidation("nodes/index.html", v)
	}
	c.Flash.Success("节点创建成功")
	return c.Redirect("/nodes")
}

func (c Nodes) Edit() revel.Result {
	node := Node{}
	err := DB.First(&node, c.Params.Get("id")).Error
	if err != nil {
		return c.RenderError(err)
	}

	c.RenderArgs["node"] = node
	return c.Render("nodes/edit.html")
}

func (c Nodes) Update() revel.Result {
	node := Node{}
	err := DB.First(&node, c.Params.Get("id")).Error
	if err != nil {
		return c.RenderError(err)
	}
	node.Name = c.Params.Get("name")
	node.Summary = c.Params.Get("summary")
	v := UpdateNode(&node)

	c.RenderArgs["node"] = node
	if v.HasErrors() {
		return c.renderValidation("nodes/edit.html", v)
	}
	c.Flash.Success("节点更新成功")
	return c.Redirect("/nodes")
}

func (c Nodes) Delete() revel.Result {
	node := Node{}
	err := DB.First(&node, c.Params.Get("id")).Error
	if err != nil {
		return c.RenderError(err)
	}

	DB.Delete(&node)
	return c.Redirect("/nodes")
}
