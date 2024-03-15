package utils

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CheckPermission Todo: 放在缓存
// CheckPermission 检查权限
func CheckPermission(c *gin.Context, db *gorm.DB) error {
	return nil
	//u, err := url.Parse(c.Request.RequestURI)
	//if err != nil {
	//	panic(err)
	//}
	//slug := strings.Replace(u.Path, "/", "_", -1)[1:]
	//
	//// 用户名
	//userName, _ := c.Get(types.LoginUserID)
	//
	//var manager model.AdminUser
	//db.Table("chatgpt_admin_users").Select("chatgpt_admin_users.id").Where("username = ?", userName).First(&manager)
	//
	//// 超级管理员不判断
	//if manager.Id == 1 {
	//	return nil
	//}
	//var roleIds []int
	//var count int64
	//db.Raw("SELECT `chatgpt_admin_user_roles`.role_id FROM `chatgpt_admin_users` LEFT JOIN `chatgpt_admin_user_roles` ON ( `chatgpt_admin_users`.id = `chatgpt_admin_user_roles`.admin_id ) WHERE `chatgpt_admin_users`.id = ?", manager.Id).Find(&roleIds)
	//db.Raw("SELECT `chatgpt_admin_permissions`.slug FROM `chatgpt_admin_permissions` LEFT JOIN `chatgpt_admin_role_permissions` ON (`chatgpt_admin_permissions`.id = `chatgpt_admin_role_permissions`.permission_id) WHERE `chatgpt_admin_role_permissions`.role_id IN ? and `chatgpt_admin_permissions`.slug = ? ", roleIds, slug).Count(&count)
	//if count > 0 {
	//	return nil
	//}
	//return fmt.Errorf("没有权限")
}
