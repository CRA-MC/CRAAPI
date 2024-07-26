package cmd

import (
	"craapi/cmd/packages/mongodb"
	"craapi/cmd/packages/panel"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fasthttp/router"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var cfgFile string
var v bool

func runServer() {
	if v {
		versionPrint()
		os.Exit(0)
	}

	fmt.Println("config file:" + cfgFile)
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("toml")
	} else {
		viper.SetConfigFile("./craapi.config")
		viper.SetConfigType("toml")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
	mongodb.Mongodb_INIT(viper.GetString("mongoDB.User"), viper.GetString("mongoDB.Password"), viper.GetString("mongoDB.Host"), viper.GetInt("mongoDB.Port"), viper.GetString("mongoDB.Database"))

	fmt.Println("API Server is tring to run at:", viper.GetString("Servername")+":"+strconv.Itoa(viper.GetInt("Port")))
	router := router.New()
	// 不同的路由执行不同的处理函数
	// 用户面板
	if viper.GetBool("userpanel.Enable") {
		panelconfig := panel.PanelDefaultConfig_t{
			DefaultUserGroup:        viper.GetString("DefaultUserGroup"),
			DefaultPerferedLanguage: viper.GetString("DefaultPerferedLanguage"),
			SmtpAddress:             viper.GetString("Smtp.Address"),
			SmtpServer:              viper.GetString("Smtp.Server"),
			SmtpAuthCode:            viper.GetString("Smtp.AuthCode"),
			SmtpUsername:            viper.GetString("Smtp.Username"),
		}
		panel.Panelinit(&panelconfig)
		fmt.Println("User Panel Enabled")
		router.GET("/", panel.Index_page_get)
		router.GET("/login", panel.Login_page_get)
		router.POST("/login", panel.Login_page_post)
		router.GET("/register", panel.Register_page_get)
		router.POST("/register", panel.Register_page_get)
		router.GET("/panel", panel.UserpanelFunc_get)
		router.POST("/panel", panel.UserpanelFunc_get)
		router.GET("/adminpanel", panel.Adminpanel_page_get)
		router.POST("/api/newuseremailcheck", panel.NewUserEmailCheckPost)
		router.POST("/api/newusernamecheck", panel.NewUserNameCheckPost)
		router.POST("/api/login", panel.Api_Login)
		router.POST("/api/emailauth", panel.Api_Email_Auth)
	}
	// Yggdrasil api
	if viper.GetBool("yggdrasilapi.Enable") {
		router.GET("/mcauth/", panel.Index_page_get)

	}

	// 启动web服务器，监听 servername:port
	server := &fasthttp.Server{
		Handler:      router.Handler,
		ReadTimeout:  time.Duration(viper.GetInt("ReadTimeout")) * time.Second,
		WriteTimeout: time.Duration(viper.GetInt("WriteTimeout")) * time.Second,
		Concurrency:  256 * 1024,
	}
	log.Fatal(server.ListenAndServe(viper.GetString("Servername") + ":" + strconv.Itoa(viper.GetInt("Port"))))
	err := http.ListenAndServe(viper.GetString("Servername")+":"+strconv.Itoa(viper.GetInt("Port")), nil)
	if err != nil {
		fmt.Println("HTTP SERVER failed,err:", err)
		return
	}
	fmt.Println("Closing API server ...")
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
		fmt.Println(err)
		os.Exit(1)
	}
}
