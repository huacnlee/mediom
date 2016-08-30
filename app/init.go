package app

import (
	"fmt"
	"github.com/cbonello/revel-csrf"
	"github.com/huacnlee/mediom/app/models"
	"github.com/huacnlee/train"
	"github.com/qor/qor"
	"github.com/qor/admin"
	"github.com/qor/publish"
	"github.com/qor/sorting"
	"github.com/qor/validations"
	"github.com/revel/revel"
	"net/http"
	"strings"
)

var Admin *admin.Admin
var Publish *publish.Publish
var mux *http.ServeMux

func init() {
	// fmt.Println("Start app with dev mode:", revel.Config.BoolDefault("mode.dev", false))
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter, // Recover from panics and display an error page instead.
		AdminFilter,
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		csrf.CSRFFilter,               // CSRF
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	train.Config.AssetsPath = "app/assets"
	train.Config.SASS.DebugInfo = false
	train.Config.SASS.LineNumbers = false
	train.Config.Verbose = false
	train.Config.BundleAssets = true

	csrf.ExemptedGlob("/msg")

	revel.OnAppStart(func() {
		models.InitDatabase()
		initAdmin()

		if revel.DevMode {
			train.ConfigureHttpHandler(nil)
			revel.Filters = append([]revel.Filter{AssetsFilter}, revel.Filters...)
		}
	})

	revel.TemplateFuncs["javascript_include_tag"] = train.JavascriptTag
	revel.TemplateFuncs["stylesheet_link_tag"] = train.StylesheetTag
}

var AssetsFilter = func(c *revel.Controller, fc []revel.Filter) {
	if strings.HasPrefix(c.Request.URL.Path, "/assets") {
		train.ServeRequest(c.Response.Out, c.Request.Request)
	} else {
		fc[0](c, fc[1:])
	}
}

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
	setting.NewAttrs("Key", "Val")
	setting.EditAttrs("Key", "Val")
	setting.IndexAttrs("Id", "Key", "CreatedAt", "UpdatedAt")

	reply := Admin.AddResource(&models.Reply{})
	reply.NewAttrs("TopicId", "UserId", "Body")
	reply.EditAttrs("Body")
	reply.IndexAttrs("Id", "Topic", "User", "Body", "CreatedAt", "UpdatedAt")
	reply.Meta(bodyMeta)

	user := Admin.AddResource(&models.User{})
	user.SearchAttrs("Login", "Email")
	user.EditAttrs("Login", "Email", "Location", "GitHub", "Twitter", "HomePage", "Tagline", "Description")
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
