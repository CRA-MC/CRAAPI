package cmd

import (
	"craapi/packages/api"
	"craapi/packages/cos"
	"craapi/packages/encryption"
	"craapi/packages/log"
	"craapi/packages/mongodb"
	"craapi/packages/panel"
	"craapi/packages/register"
	"craapi/packages/yggdrasilapi"
	"fmt"
	"mime"
	"net/smtp"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fasthttp/router"
	"github.com/jordan-wright/email"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var cfgFile string
var v bool

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func runServer() {
	versionPrint()
	if v {
		os.Exit(0)
	}

	go log.LOG()

	var err error

	log.LOGI("config file:" + cfgFile)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("toml")
	} else {
		viper.SetConfigFile("./craapi.config")
		viper.SetConfigType("toml")
	}
	err = viper.ReadInConfig()

	if err != nil {
		log.LOGI("Can't read config:", err)
		os.Exit(1)
	}

	if fileExists("ID_RSA") && fileExists("ID_RSA.pub") {
		encryption.Privatekey, err = os.ReadFile("ID_RSA")
		if err != nil {
			log.LOGI("Can't read ID_RSA:", err)
			os.Exit(1)
		}
		encryption.Publickey, err = os.ReadFile("ID_RSA.pub")
		if err != nil {
			log.LOGI("Can't read ID_RSA.pub:", err)
			os.Exit(1)
		}
	} else {
		GenRsaKey(4096)
	}

	mongodb.Mongodb_INIT(viper.GetString("mongoDB.User"), viper.GetString("mongoDB.Password"), viper.GetString("mongoDB.Host"), viper.GetInt("mongoDB.Port"), viper.GetString("mongoDB.Database"))
	mongodb.Mongodb_GetCollections()

	log.LOGI("API Server is tring to run at: ", viper.GetString("Servername")+":"+strconv.Itoa(viper.GetInt("Port")))

	if viper.GetBool("tencentCOS.Enable") {
		cos.COS_INIT()
		// fmt.Println(cos.GetCOSURL("craregister.exe"))
	}

	router := router.New()

	Email_ch := make(chan *email.Email, 10)
	p, err := email.NewPool(
		viper.GetString("smtp.Host")+":"+strconv.Itoa(viper.GetInt("smtp.Port")),
		viper.GetInt("smtp.Connection"),
		smtp.PlainAuth("", viper.GetString("smtp.Username"), viper.GetString("smtp.AuthCode"), viper.GetString("smtp.Host")),
	)
	if err != nil {
		log.LOGE("failed to create pool:", err)
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(viper.GetInt("smtp.Connection"))
	for i := 0; i < viper.GetInt("smtp.Connection"); i++ {
		go func() {
			defer wg.Done()
			for e := range Email_ch {
				err := p.Send(e, 10*time.Second)
				if err != nil {
					fmt.Fprintf(os.Stderr, "email to %v sent error:%v\n", e.To, err)
				}
			}
		}()
	}
	if viper.GetBool("smtp.Debug") {
		e := email.NewEmail()
		e.From = mime.QEncoding.Encode("UTF-8", viper.GetString("smtp.EmailName")+" <"+viper.GetString("smtp.Address")+">")
		e.To = []string{viper.GetString("smtp.DebugSendAddress")}
		e.Subject = "Craapi Test"
		e.Text = []byte("Craapi email test message")
		Email_ch <- e
	}

	register.Init("user", &Email_ch)

	// 不同的路由执行不同的处理函数
	// 用户面板
	router.NotFound = panel.D404
	router.HandleMethodNotAllowed = false
	if viper.GetBool("userpanel.Enable") {
		panelconfig := panel.PanelDefaultConfig_t{
			DefaultUserGroup:        viper.GetString("DefaultUserGroup"),
			DefaultPerferedLanguage: viper.GetString("DefaultPerferedLanguage"),
			SmtpAddress:             viper.GetString("smtp.Address"),
			SmtpServer:              viper.GetString("smtp.Server"),
			SmtpAuthCode:            viper.GetString("smtp.AuthCode"),
			SmtpUsername:            viper.GetString("smtp.Username"),
		}
		panel.Panelinit(&panelconfig)
		log.LOGI("User Panel Enabled")
		router.GET("/", panel.Index_page_get)
		router.POST("/", panel.Index_page_get)
		router.GET("/login", panel.Login_page_get)
		router.POST("/login", panel.Login_page_post)
		router.GET("/register", panel.Register_page_get)
		router.POST("/register", panel.Register_page_post)
		router.GET("/panel", panel.UserpanelFunc_get)
		router.POST("/panel", panel.UserpanelFunc_get)
		// router.GET("/adminpanel", panel.Adminpanel_page_get)
	}

	// craapi
	// router.POST("/api/newuseremailcheck", panel.NewUserEmailCheckPost)
	// router.POST("/api/newusernamecheck", panel.NewUserNameCheckPost)
	router.POST("/api/login", api.Api_Login)
	router.POST("/api/emailauth", api.Api_EmailAuth)
	router.POST("/api/register", api.Api_Register)
	router.POST("/api/getuserinfo", api.Api_GetUserinfo)

	// Yggdrasil api
	if viper.GetBool("yggdrasilapi.Enable") {
		yggdrasilapi.Yggdrasilapiinit(viper.GetStringSlice("yggdrasilapi.SkinDomains"), viper.GetString("Domain"))
		router.GET("/mcauth", yggdrasilapi.Yggdrasil)
		router.POST("/mcauth/authserver/authenticate", yggdrasilapi.UserAuth)
		router.POST("/mcauth/authserver/refresh", yggdrasilapi.TokenRefresh)
		router.POST("/mcauth/authserver/validate", yggdrasilapi.TokenVail)
		router.POST("/mcauth/authserver/invalidate", yggdrasilapi.TokeninValid)
		router.POST("/mcauth/authserver/signout", yggdrasilapi.Signout)
		router.POST("/mcauth/sessionserver/session/minecraft/join", yggdrasilapi.ServerJoin)
		router.GET("/mcauth/sessionserver/session/minecraft/hasJoined", yggdrasilapi.ServerhasJoined)
		router.GET("/mcauth/sessionserver/session/minecraft/profile/{uuid}", yggdrasilapi.ProfileSearch)
		router.POST("/mcauth/api/profiles/minecraft", yggdrasilapi.ProfileMutiSearch)
		router.POST("/mcauth/minecraftservices/publickeys", yggdrasilapi.Yggdrasilpubkey)
		router.GET("/mcauth/minecraftservices/publickeys", yggdrasilapi.Yggdrasilpubkey)
	}

	// 启动web服务器，监听 servername:port
	server := &fasthttp.Server{
		Handler:      router.Handler,
		ReadTimeout:  time.Duration(viper.GetInt("ReadTimeout")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("WriteTimeout")) * time.Second,
		Concurrency:  256 * 1024,
	}
	err = server.ListenAndServe(viper.GetString("Servername") + ":" + strconv.Itoa(viper.GetInt("Port")))
	if err != nil {
		log.LOGE("HTTP SERVER failed,err:", err)
		panic(err)
	}
	log.LOGI("Closing API server ...")
}

var rootCmd = &cobra.Command{
	Use:   "craapi",
	Short: "CRA-MC API is a minecraft servers API",
	Long: `CRA-MC API is a minecraft servers API which can manage mutiple minecraft servers' mods version or files
	`,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "./craapi.config", "config file")
	rootCmd.Flags().BoolVarP(&v, "version", "v", false, "version")
	viper.SetDefault("Servername", "")
	viper.SetDefault("Port", "19999")
}

func initConfig() {
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.LOGI(err)
		os.Exit(1)
	}
}
