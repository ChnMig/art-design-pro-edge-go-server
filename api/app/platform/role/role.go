package role

import (
	"github.com/gin-gonic/gin"

	systemrole "api-server/api/app/system/role"
)

func GetRoleList(c *gin.Context) {
	systemrole.GetRoleList(c)
}

func AddRole(c *gin.Context) {
	systemrole.AddRole(c)
}

func UpdateRole(c *gin.Context) {
	systemrole.UpdateRole(c)
}

func DeleteRole(c *gin.Context) {
	systemrole.DeleteRole(c)
}
