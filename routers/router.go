package routers

import (
	"ConcurrentSearch/controller"
	"github.com/gin-gonic/gin"
)

func SetRouter() *gin.Engine {
	r := gin.Default()
	//加载静态资源
	r.Static("/static/", "view/static/")
	//解析模板
	r.LoadHTMLGlob("view/index/*")
	index := r.Group("/index")
	{
		//主页
		index.GET("/", controller.IndexHandler)
		//反馈建立索引完成信号
		//index.GET("/init", controller.IsOK)
		//搜素返回搜索结果
		index.GET("/search", controller.Search)
	}
	return r
}
