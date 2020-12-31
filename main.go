package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type test struct {
	ID       int    `gorm:"column:id" json:"id"`
	Testnull string `gorm:"column:testnull" json:"testnull"`
}

func main() {
	DB, err := gorm.Open("mysql", "root:123456@/go_shop?charset=utf8&parseTime=True&loc=Local") //连接MYSQL
	if err != nil {
		log.Fatal(err.Error())
	}
	defer DB.Close()
	router := gin.Default()
	router.Use(Cors())
	router.LoadHTMLGlob("html/*")
	pan := router.Group("/pan")
	{
		pan.StaticFS("/static", http.Dir("./static"))
		pan.GET("/list", func(context *gin.Context) {
			context.HTML(http.StatusOK, "list.html", gin.H{})
		})
		pan.GET("/upload", func(context *gin.Context) {
			context.HTML(http.StatusOK, "upload.html", gin.H{
				"title": "Main website",
			})
		})
		pan.POST("/uploadfile", func(context *gin.Context) {
			// slicename := uuid.New()
			file, err := context.FormFile("file")
			index, _ := context.FormFile("index")
			fmt.Println(index)
			if err != nil {
				log.Fatal(err.Error())
			}
			os.Mkdir("./static/files/"+file.Filename, 0666)
			dst := fmt.Sprintf("./static/files/"+file.Filename+"/%s", file.Filename)
			context.SaveUploadedFile(file, dst)
			if strings.Contains(" XpartX0", file.Filename) {
				//开始

			}
		})
	}
	router.Run(":8080")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("origin")
		var headerKeys []string
		for key, _ := range c.Request.Header {
			headerKeys = append(headerKeys, key)
		}
		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin，acceSs-control-allow-headers，%s", headerStr)
		} else {
			headerStr = "access-control-allow-origin,access-control-allow-headers"
		}
		//请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()
		c.Next()
	}
}
