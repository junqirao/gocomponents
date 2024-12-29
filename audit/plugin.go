package audit

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/junqirao/gocomponents/gfutil"
	"github.com/junqirao/gocomponents/response"
)

func Plugin(prefix string, middleware ...ghttp.HandlerFunc) ghttp.Plugin {
	if strings.HasSuffix(prefix, "/") {
		prefix = prefix[:len(prefix)-1]
	}
	return gfutil.NewPlugin(
		gfutil.WithPrefix(fmt.Sprintf("%s%s", prefix, "/audit")),
		gfutil.WithMiddleware(middleware...),
		gfutil.WithName("audit-record-helper"),
		gfutil.WithDescription("audit record helper"),
		gfutil.WithVersion("v1.0.0"),
		gfutil.WithAuthor("junqirao"),
		gfutil.WithInstallHandler(func(group *ghttp.RouterGroup) {
			group.POST("/record", query)
			group.GET("/supported-modules", supportedModules)
		}),
	)
}

func query(r *ghttp.Request) {
	req := new(RecordQueryParams)
	err := r.Parse(req)
	if err != nil {
		err = response.CodeInvalidParameter.WithDetail(err.Error())
		return
	}

	res, err := Logger.adaptor.Load(r.Context(), req)
	if err != nil {
		response.WriteJSON(r, response.CodeFromError(err))
		return
	}
	response.WriteData(r, response.DefaultSuccess(), res)
}

func supportedModules(r *ghttp.Request) {
	response.WriteData(r, response.DefaultSuccess(), SupportedModules())
}
