package api

import (
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/locales"
	"io"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/mlogclub/simple/web"

	"bbs-go/internal/services"
)

type UploadController struct {
	Ctx iris.Context
}

func (c *UploadController) Post() *web.JsonResult {
	user := services.UserTokenService.GetCurrent(c.Ctx)
	if err := services.UserService.CheckPostStatus(user); err != nil {
		return web.JsonError(err)
	}

	// 支持上传任意类型的文件，不再限制为"image"
	file, header, err := c.Ctx.FormFile("file")
	if err != nil {
		// 兼容旧的"image"字段名
		file, header, err = c.Ctx.FormFile("image")
		if err != nil {
			return web.JsonError(err)
		}
	}
	defer file.Close()

	if header.Size > constants.UploadMaxBytes {
		return web.JsonErrorMsg(locales.Getf("upload.file_too_large", constants.UploadMaxM))
	}

	contentType := header.Header.Get("Content-Type")
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return web.JsonError(err)
	}

	slog.Info("上传文件：", slog.Any("filename", header.Filename), slog.Any("size", header.Size), slog.Any("contentType", contentType))

	// 检查文件类型是否被允许
	if !isAllowedFileType(header.Filename, contentType) {
		return web.JsonErrorMsg("不支持的文件类型：该文件类型可能存在安全风险，禁止上传")
	}

	// 使用新的PutObject方法处理任意文件类型
	url, err := services.UploadService.PutObject("", fileBytes, contentType)
	if err != nil {
		return web.JsonError(err)
	}
	
	// 返回文件信息，包括文件名和类型
	fileExt := strings.ToLower(filepath.Ext(header.Filename))
	isImage := strings.HasPrefix(contentType, "image/")
	
	return web.NewEmptyRspBuilder().
		Put("url", url).
		Put("filename", header.Filename).
		Put("contentType", contentType).
		Put("isImage", isImage).
		Put("fileExt", fileExt).
		JsonResult()
}

// isAllowedFileType 检查文件类型是否被允许上传
// 现在支持任意文件类型，但会检查危险文件类型
func isAllowedFileType(filename, contentType string) bool {
	// 禁止的危险文件类型
	dangerousExts := map[string]bool{
		".exe":  true, ".bat": true, ".cmd": true, ".sh": true, ".ps1": true,
		".dll":  true, ".sys": true, ".com": true, ".scr": true, ".pif": true,
		".app":  true, ".apk": true, ".ipa": true, ".deb": true, ".rpm": true,
		".msi":  true, ".msp": true, ".mst": true,
		".vbs":  true, ".js":  true, ".jar": true, ".class": true, ".php": true,
		".py":   true, ".pl":  true, ".rb":  true, ".cgi": true, ".asp": true,
		".aspx": true, ".jsp": true, ".html": true, ".htm": true,
	}
	
	// 禁止的危险内容类型
	dangerousContentTypes := map[string]bool{
		"application/x-msdownload":             true, // Windows可执行文件
		"application/x-msdos-program":          true,
		"application/x-executable":             true,
		"application/x-sh":                     true, // Shell脚本
		"application/x-shellscript":            true,
		"application/vnd.android.package-archive": true, // APK文件
		"application/x-iphone-app":             true, // iOS应用
		"application/x-java-archive":           true, // JAR文件
		"application/javascript":               true, // JavaScript文件
		"application/x-httpd-php":              true, // PHP文件
		"text/html":                           true, // HTML文件
		"application/x-cgi":                    true, // CGI脚本
	}
	
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	
	// 检查是否为危险文件类型
	if dangerousExts[ext] || dangerousContentTypes[contentType] {
		return false
	}
	
	// 允许所有其他文件类型
	return true
}
