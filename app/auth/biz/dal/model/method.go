package model

import (
	"context"

	"gorm.io/gorm"
)

// HasUserPermission 判断指定用户（user_id）是否具有指定服务(serviceName)下的某个资源权限(resourceName)
func HasUserPermission(db *gorm.DB, ctx context.Context, userID int64, serviceName, resourceName string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).
		Table("user_role").
		Joins("JOIN role_permission ON user_role.role_code = role_permission.role_code").
		Joins("JOIN permission ON permission.permission_code = role_permission.permission_code").
		Where("user_role.user_id = ? AND permission.service_name = ? AND permission.resource_name = ?", userID, serviceName, resourceName).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetUserPermissions 获取指定用户所有拥有的权限
func GetUserPermissions(db *gorm.DB, ctx context.Context, userID int64) ([]Permission, error) {
	var permissions []Permission
	err := db.WithContext(ctx).
		Table("user_role").
		Joins("JOIN role_permission ON user_role.role_code = role_permission.role_code").
		Joins("JOIN permission ON permission.permission_code = role_permission.permission_code").
		Where("user_role.user_id = ?", userID).
		Select("permission.*").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetPermissionsByRole 查询指定角色（roleCode）的所有权限
func GetPermissionsByRole(db *gorm.DB, ctx context.Context, roleCode int) (permissions []Permission, err error) {
	err = db.WithContext(ctx).
		Table("role_permission").
		Joins("JOIN permission ON permission.permission_code = role_permission.permission_code").
		Where("role_permission.role_code = ?", roleCode).
		Select("permission.*").
		Find(&permissions).Error
	return
}

// IsPublicPermission 判断给定资源（resourceName）是否为公共权限
func IsPublicPermission(db *gorm.DB, ctx context.Context, resourceName string) (bool, error) {
	var count int64
	err := db.WithContext(ctx).
		Model(&Permission{}).
		Where("resource_name = ? AND is_public = ?", resourceName, true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetRolesByUser 查询指定 userID 的所有角色信息
func GetRolesByUser(db *gorm.DB, ctx context.Context, userID int64) (roles []Role, err error) {
	err = db.WithContext(ctx).
		Table("user_role").
		Joins("JOIN role ON role_permission.role_code = role.role_code").
		Where("user_role.user_id = ?", userID).
		Select("role.*").
		Find(&roles).Error
	return
}
