package app

import (
	// "github.com/cbonello/revel-csrf"
	"github.com/huacnlee/mediom/app/models"
	"github.com/huacnlee/train"
	"github.com/revel/revel"
	"strings"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter, // Recover from panics and display an error page instead.
		AdminFilter,
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		// csrf.CSRFFilter,               // CSRF
		revel.ValidationFilter,  // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,        // Resolve the requested language
		revel.InterceptorFilter, // Run interceptors around the action.
		revel.CompressFilter,    // Compress the result.
		revel.ActionInvoker,     // Invoke the action.
	}

	train.Config.AssetsPath = "app/assets"
	if revel.RunMode == "prod" {
		train.Config.AssetsPath = "src/github.com/huacnlee/mediom/app/assets"
	}
	train.Config.SASS.DebugInfo = false
	train.Config.SASS.LineNumbers = false
	train.Config.Verbose = false
	train.Config.BundleAssets = true

	// csrf.ExemptedGlob("/msg")

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
