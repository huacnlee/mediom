package app

import (
	"fmt"
	"github.com/revel/revel"
	"github.com/huacnlee/mediom/app/models"
	"github.com/qor/qor"
	"github.com/qor/admin"
	"github.com/qor/publish"
	"github.com/qor/sorting"
	"github.com/qor/validations"
	"net/http"
	"strings"
)

var Admin *admin.Admin
var Publish *publish.Publish
var mux *http.ServeMux

var AdminFilter = func(c *revel.Controller, fc []revel.Filter) {
	if strings.HasPrefix(c.Request.URL.Path, "/admin") {
		mux.ServeHTTP(c.Response.Out, c.Request.Request)
	} else {
		fc[0](c, fc[1:])
	}
}

func initAdmin() {
	Admin = admin.New(&qor.Config{DB: models.DB})
	Publish = publish.New(models.DB)
	sorting.RegisterCallbacks(models.DB)
	validations.RegisterCallbacks(models.DB)

	nodeSelectMeta := &admin.Meta{Name: "NodeId", Type: "select_one", Collection: nodeCollection}
	bodyMeta := &admin.Meta{Name: "Body", Type: "text"}

	topic := Admin.AddResource(&models.Topic{})
	topic.SearchAttrs("Title")
	topic.NewAttrs("Title", "NodeId", "Body", "UserId")
	topic.EditAttrs("Title", "NodeId", "Body", "UserId")
	topic.IndexAttrs("Id", "UserId", "Title", "NodeId", "RepliesCount", "CreatedAt", "UpdatedAt")
	topic.Meta(bodyMeta)
	topic.Meta(nodeSelectMeta)

	setting := Admin.AddResource(&models.Setting{})
	setting.EditAttrs("Key", "Val")
  setting.Meta(&admin.Meta{Name: "Val", Type: "text"})
	setting.IndexAttrs("Id", "Key", "CreatedAt", "UpdatedAt")

	reply := Admin.AddResource(&models.Reply{})
	reply.NewAttrs("TopicId", "UserId", "Body")
	reply.EditAttrs("Body")
	reply.IndexAttrs("Id", "Topic", "User", "Body", "CreatedAt", "UpdatedAt")
	reply.Meta(bodyMeta)

	user := Admin.AddResource(&models.User{})
	user.SearchAttrs("Login", "Email")
	user.EditAttrs("Login", "Email", "Location", "GitHub", "Twitter", "HomePage", "Tagline", "Description")
  user.Meta(&admin.Meta{Name: "Description", Type: "text"})
	user.IndexAttrs("Id", "Login", "Email", "Location", "CreatedAt", "UpdatedAt")

	node := Admin.AddResource(&models.Node{})
	node.IndexAttrs("Id", "ParentId", "Name", "Summary", "Sort")
	node.NewAttrs("ParentId", "Name", "Summary", "Sort")
	node.EditAttrs("ParentId", "Name", "Summary", "Sort")


	notification := Admin.AddResource(&models.Notification{})
	notification.EditAttrs("Id")
	notification.IndexAttrs("Id", "UserId", "ActorId", "NotifyType", "Read", "NotifyableType", "NotifyableId", "CreatedAt", "UpdatedAt")

	mux = http.NewServeMux()
	mux.Handle("/system/", http.FileServer(http.Dir("public")))
	Admin.MountTo("/admin", mux)
}

func nodeCollection(resource interface{}, context *qor.Context) (results [][]string) {
	nodes := models.FindAllNodes()
	for _, node := range nodes {
		results = append(results, []string{fmt.Sprintf("%v", node.Id), node.Name})
	}
	return
}

func nodeRootCollection(resource interface{}, context *qor.Context) (results [][]string) {
	roots := models.FindAllNodeRoots()
	for _, node := range roots {
		results = append(results, []string{fmt.Sprintf("%v", node.Id), node.Name})
	}
	return
}
