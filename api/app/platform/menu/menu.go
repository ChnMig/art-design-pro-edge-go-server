package menu

import (
	"github.com/gin-gonic/gin"

	systemmenu "api-server/api/app/system/menu"
)

func GetMenuList(c *gin.Context) {
	systemmenu.GetMenuList(c)
}

func AddMenu(c *gin.Context) {
	systemmenu.AddMenu(c)
}

func UpdateMenu(c *gin.Context) {
	systemmenu.UpdateMenu(c)
}

func DeleteMenu(c *gin.Context) {
	systemmenu.DeleteMenu(c)
}

func GetMenuAuthList(c *gin.Context) {
	systemmenu.GetMenuAuthList(c)
}

func AddMenuAuth(c *gin.Context) {
	systemmenu.AddMenuAuth(c)
}

func UpdateMenuAuth(c *gin.Context) {
	systemmenu.UpdateMenuAuth(c)
}

func DeleteMenuAuth(c *gin.Context) {
	systemmenu.DeleteMenuAuth(c)
}
