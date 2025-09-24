package panel

import (
	"craapi/packages/log"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/jordan-wright/email"
	"github.com/valyala/fasthttp"
)

type PanelDefaultConfig_t struct {
	DefaultUserGroup        string
	DefaultPerferedLanguage string
	SmtpAddress             string
	SmtpServer              string
	SmtpAuthCode            string
	SmtpUsername            string
}

var defaultconfig *PanelDefaultConfig_t
var email_send email.Email

func Panelinit(config *PanelDefaultConfig_t) {
	defaultconfig = config
}

func isExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

type D404tpl struct {
	Lasturi string
	Fullurl string
}

func staticfileget(ctx *fasthttp.RequestCtx, filename string) {
	ctx.SetContentType("text/html; charset=utf-8")
	file, err := os.ReadFile("./templates/" + filename)
	if err != nil {
		fmt.Println("Requesting file", string(ctx.URI().LastPathSegment()), "does not exist. err:", err)
		//传输404模板
		tmpl, err := template.ParseFiles("./templates/404.tpl")
		if err != nil {
			fmt.Println("create template failed, err", err)
			ctx.Redirect("/static/403.html", 307)
			ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
			return
		}
		c404 := D404tpl{Lasturi: string(ctx.URI().LastPathSegment()), Fullurl: string(ctx.URI().FullURI())}
		tmpl.Execute(ctx, c404)
		//传输404模板结束
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	ctx.Write(file)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

var D404tmpl *template.Template

var D404t bool

func init() {
	var err error
	D404tmpl, err = template.ParseFiles("./templates/404.html")
	D404t = true
	if err != nil {
		D404t = false
		log.LOGE("create template failed, err", err)
	}
}

func D404(ctx *fasthttp.RequestCtx) {
	log.LOGI(
		"D404",
		" URL: ",
		" ",
		ctx.Request.Header.Method(),
		" ",
		string(ctx.URI().FullURI()),
		" Not found! ",
		"Body: ",
		string(ctx.Request.Body()),
	)
	ctx.SetContentType("text/html; charset=utf-8")
	if D404t {
		c404 := D404tpl{Lasturi: string(ctx.URI().LastPathSegment()), Fullurl: string(ctx.URI().FullURI())}
		D404tmpl.Execute(ctx, c404)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	ctx.Redirect("/static/403.html", 307)
	ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
}

type Logined_t struct {
	Logined  bool
	Username string
}

func Index_page_get(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("X-Authlib-Injector-API-Location", string(ctx.URI().Host())+"/mcauth")
	logined := Logined_t{
		Logined:  false,
		Username: "",
	}
	ctx.SetContentType("text/html; charset=utf-8")
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("create template failed, err", err)
		ctx.Redirect("/static/403.html", 307)
		ctx.SetStatusCode(fasthttp.StatusTemporaryRedirect)
		return
	}
	tmpl.Execute(ctx, logined)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func ismobile(userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}

	isMobile := false
	mobileKeywords := []string{"Mobile", "Android", "Silk/", "Kindle",
		"BlackBerry", "Opera Mini", "Opera Mobi"}

	for i := 0; i < len(mobileKeywords); i++ {
		if strings.Contains(userAgent, mobileKeywords[i]) {
			isMobile = true
			break
		}
	}

	return isMobile
}

func UserpanelFunc_get(ctx *fasthttp.RequestCtx) {
	// ctx.WriteString("信息：")
	// fmt.Println(ctx.Request.Header.String())
	// ctx.WriteString("\n是否是移动端：")
	// ctx.Write(ctx.Request.Header.PeekBytes([]byte("Sec-Ch-Ua-Mobile")))
	// fmt.Println(ctx.Request.Header.PeekBytes([]byte("Sec-Ch-Ua-Mobile")))
	// ctx.WriteString("\n平台：")
	// ctx.Write(ctx.Request.Header.PeekBytes([]byte("Sec-Ch-Ua-Platform")))
	// fmt.Println(string(ctx.Request.Header.UserAgent()))
	mobile := ctx.Request.Header.Peek("Sec-Ch-Ua-Mobile")
	ismo := false
	if mobile != nil {
		if mobile[1] == 49 {
			// 移动端
			ismo = true
		}
	} else {
		if ismobile(string(ctx.Request.Header.UserAgent())) {
			ismo = true
		}
	}
	if ismo {
		ctx.WriteString("移动端")
	} else {
		ctx.WriteString("其他端")
	}
}

func Adminpanel_page_get(ctx *fasthttp.RequestCtx) {
	staticfileget(ctx, "adminpanel.tpl")
}
