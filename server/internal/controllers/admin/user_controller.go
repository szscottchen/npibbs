package admin

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/errs"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bbs-go/internal/models"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/mlogclub/simple/sqls"
	"github.com/mlogclub/simple/web"
	"github.com/mlogclub/simple/web/params"

	"bbs-go/internal/services"
)

type UserController struct {
	Ctx iris.Context
}

// BeforeActivation is called before the controller's dependencies binding to the fields
// but before server ran. It's used to customize controller routes manually.
func (c *UserController) BeforeActivation(b mvc.BeforeActivation) {
	// Manually register the route for /download-template which maps to DownloadImportTemplate method
	// This is needed because Iris MVC doesn't automatically convert hyphens in URLs to camelCase method names
	b.Handle("GET", "/download-template", "DownloadImportTemplate")
	
	// 手动注册批量导入用户的路由，包括POST和OPTIONS方法
	b.Handle("POST", "/batch-import", "PostBatch_import")
	b.Handle("OPTIONS", "/batch-import", "OptionsBatch_import")
}

func (c *UserController) GetSynccount() *web.JsonResult {
	go func() {
		services.UserService.SyncUserCount()
	}()
	return web.JsonSuccess()
}

func (c *UserController) GetBy(id int64) *web.JsonResult {
	t := services.UserService.Get(id)
	if t == nil {
		return web.JsonErrorMsg("Not found, id=" + strconv.FormatInt(id, 10))
	}
	return web.JsonData(c.buildUserItem(t, true))
}

func (c *UserController) AnyList() *web.JsonResult {
	queryParams := params.NewQueryParams(c.Ctx)
	
	// 单独处理username查询，避免GORM对数字字符串的类型转换问题
	username := params.FormValue(c.Ctx, "username")
	if len(username) > 0 {
		queryParams.Cnd.Eq("username", username) // 直接使用字符串值，确保类型正确
	}
	
	list, paging := services.UserService.FindPageByParams(queryParams.
		EqByReq("id").
		LikeByReq("nickname").
		EqByReq("type").
		PageByReq().Desc("id"))
	var itemList []map[string]interface{}
	for _, user := range list {
		itemList = append(itemList, c.buildUserItem(&user, false))
	}
	return web.JsonData(&web.PageResult{Results: itemList, Page: paging})
}

func (c *UserController) PostCreate() *web.JsonResult {
	username := params.FormValue(c.Ctx, "username")
	email := params.FormValue(c.Ctx, "email")
	nickname := params.FormValue(c.Ctx, "nickname")
	password := params.FormValue(c.Ctx, "password")

	user, err := services.UserService.SignUp(username, email, nickname, password, password)
	if err != nil {
		return web.JsonError(err)
	}
	return web.JsonData(c.buildUserItem(user, true))
}

func (c *UserController) PostUpdate() *web.JsonResult {
	var (
		id, _       = params.GetInt64(c.Ctx, "id")
		_type, _    = params.GetInt(c.Ctx, "type")
		username    = params.FormValue(c.Ctx, "username")
		email       = params.FormValue(c.Ctx, "email")
		nickname    = params.FormValue(c.Ctx, "nickname")
		avatar      = params.FormValue(c.Ctx, "avatar")
		gender      = params.FormValue(c.Ctx, "gender")
		homePage    = params.FormValue(c.Ctx, "homePage")
		description = params.FormValue(c.Ctx, "description")
		roleIds     = params.FormValueInt64Array(c.Ctx, "roleIds")
		status      = params.FormValueIntDefault(c.Ctx, "status", 0)
	)

	user := services.UserService.Get(id)
	if user == nil {
		return web.JsonErrorMsg("entity not found")
	}

	user.Type = _type
	user.Username = sqls.SqlNullString(username)
	user.Email = sqls.SqlNullString(email)
	user.Nickname = nickname
	user.Avatar = avatar
	user.Gender = constants.Gender(gender)
	user.HomePage = homePage
	user.Description = description
	user.Status = status

	if err := services.UserService.Update(user); err != nil {
		return web.JsonError(err)
	}
	if err := services.UserRoleService.UpdateUserRoles(user.Id, roleIds); err != nil {
		return web.JsonError(err)
	}
	user = services.UserService.Get(user.Id)
	return web.JsonData(c.buildUserItem(user, true))
}

// 禁言
func (c *UserController) PostForbidden() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if user == nil {
		return web.JsonError(errs.NotLogin())
	}
	if !user.HasAnyRole(constants.RoleOwner, constants.RoleAdmin) {
		return web.JsonErrorMsg("无权限")
	}
	var (
		userId = params.FormValueInt64Default(c.Ctx, "userId", 0)
		days   = params.FormValueIntDefault(c.Ctx, "days", 0)
		reason = params.FormValue(c.Ctx, "reason")
	)
	if userId < 0 {
		return web.JsonErrorMsg("请传入：userId")
	}
	if days == 0 {
		services.UserService.RemoveForbidden(user.Id, userId, c.Ctx.Request())
	} else {
		if err := services.UserService.Forbidden(user.Id, userId, days, reason, c.Ctx.Request()); err != nil {
			return web.JsonError(err)
		}
	}
	return web.JsonSuccess()
}

// OptionsBatch_import 处理批量导入用户的OPTIONS预检请求
func (c *UserController) OptionsBatch_import() *web.JsonResult {
	// 对于OPTIONS请求，我们只需要返回成功的状态码
	c.Ctx.Header("Access-Control-Allow-Origin", "*")
	c.Ctx.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
	c.Ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	c.Ctx.Header("Access-Control-Allow-Credentials", "true")
	c.Ctx.StatusCode(iris.StatusOK)
	return nil
}

// PostBatch_import 批量导入用户
func (c *UserController) PostBatch_import() *web.JsonResult {
	// 1. 获取上传的文件
	file, header, err := c.Ctx.FormFile("file")
	if err != nil {
		return web.JsonErrorMsg("请选择要上传的文件")
	}
	defer file.Close()

	// 2. 验证文件类型
	filename := header.Filename
	if !strings.HasSuffix(filename, ".xlsx") && !strings.HasSuffix(filename, ".xls") {
		return web.JsonErrorMsg("请上传Excel文件（.xlsx或.xls格式）")
	}

	// 3. 保存临时文件
	tempDir := "./temp"
	os.MkdirAll(tempDir, 0755)

	timestamp := time.Now().Unix()
	tempFilePath := filepath.Join(tempDir, fmt.Sprintf("import_%d_%s", timestamp, filename))

	out, err := os.Create(tempFilePath)
	if err != nil {
		return web.JsonErrorMsg("保存文件失败")
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return web.JsonErrorMsg("保存文件失败")
	}

	defer os.Remove(tempFilePath) // 处理完后删除临时文件

	// 4. 执行导入
	result, err := services.UserBatchImportService.ImportUsersFromExcel(tempFilePath)
	if err != nil {
		return web.JsonError(err)
	}

	// 5. 返回结果
	return web.JsonData(result)
}

// DownloadImportTemplate 下载导入模板
func (c *UserController) DownloadImportTemplate() *web.JsonResult {
	// 设置响应头（提前设置，确保在任何情况下都有效）
	c.Ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Ctx.Header("Content-Disposition", "attachment; filename=\"用户导入模板.xlsx\"")
	
	// 1. 生成模板文件
	tempDir := "./temp"
	os.MkdirAll(tempDir, 0755)

	templatePath := filepath.Join(tempDir, "用户导入模板.xlsx")
	err := services.UserBatchImportService.GenerateTemplate(templatePath)
	if err != nil {
		c.Ctx.StatusCode(500)
		c.Ctx.Header("Content-Type", "application/json") // 错误情况返回JSON
		return web.JsonErrorMsg("生成模板失败: " + err.Error())
	}
	defer os.Remove(templatePath)

	// 2. 读取文件内容到字节流
	fileContent, err := os.ReadFile(templatePath)
	if err != nil {
		c.Ctx.StatusCode(500)
		c.Ctx.Header("Content-Type", "application/json") // 错误情况返回JSON
		return web.JsonErrorMsg("读取模板文件失败: " + err.Error())
	}

	// 设置Content-Length
	c.Ctx.Header("Content-Length", strconv.Itoa(len(fileContent)))
	
	// 添加调试日志
	fmt.Printf("Template file size: %d bytes\n", len(fileContent))
	fmt.Printf("Template file content preview: %s\n", string(fileContent[:int(math.Min(100, float64(len(fileContent))))]))

	// 3. 直接将文件内容写入响应
	bytesWritten, err := c.Ctx.Write(fileContent)
	if err != nil {
		c.Ctx.StatusCode(500)
		c.Ctx.Header("Content-Type", "application/json") // 错误情况返回JSON
		return web.JsonErrorMsg("发送文件失败: " + err.Error())
	}
	
	// 确保写入的字节数正确
	if bytesWritten != len(fileContent) {
		c.Ctx.StatusCode(500)
		c.Ctx.Header("Content-Type", "application/json") // 错误情况返回JSON
		return web.JsonErrorMsg("文件写入不完整")
	}

	// 文件已经通过 Write 发送，返回 nil 避免额外的 JSON 响应
	return nil
}

func (c *UserController) buildUserItem(user *models.User, buildRoleIds bool) map[string]interface{} {
	b := web.NewRspBuilder(user).
		Put("roles", user.GetRoles()).
		Put("username", user.Username.String).
		Put("email", user.Email.String).
		Put("score", user.Score).
		Put("forbidden", user.IsForbidden())
	if buildRoleIds {
		b.Put("roleIds", services.UserRoleService.GetUserRoleIds(user.Id))
	}
	return b.Build()
}
