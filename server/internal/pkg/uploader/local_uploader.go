package uploader

import (
	"bbs-go/internal/models/dto"
	"bbs-go/internal/pkg/bbsurls"
	"bbs-go/internal/pkg/config"
	"mime"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/mlogclub/simple/common/dates"
	"github.com/mlogclub/simple/common/digests"
	"github.com/mlogclub/simple/common/strs"
)

type LocalUploader struct{}

func (u *LocalUploader) PutImage(cfg *dto.UploadConfig, data []byte, contentType string) (string, error) {
	if strs.IsBlank(contentType) {
		contentType = "image/jpeg"
	}
	key := generateLocalImageKey(cfg, data, contentType)
	return u.PutObject(cfg, key, data, contentType)
}

func (u *LocalUploader) PutObject(cfg *dto.UploadConfig, key string, data []byte, contentType string) (string, error) {
	// 如果key为空，生成一个合适的key
	if key == "" {
		key = generateLocalFileKey(cfg, data, contentType)
	}
	
	// 确保目录存在
	dir := filepath.Dir(key)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	// 写入文件
	if err := os.WriteFile(key, data, 0644); err != nil {
		return "", err
	}

	// 返回可访问的URL
	return bbsurls.UrlJoin(config.Instance.BaseURL, key), nil
}

func (u *LocalUploader) CopyImage(cfg *dto.UploadConfig, originUrl string) (string, error) {
	data, contentType, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return u.PutImage(cfg, data, contentType)
}

// generateLocalImageKey 生成本地存储的图片Key
func generateLocalImageKey(cfg *dto.UploadConfig, data []byte, contentType string) string {
	md5 := digests.MD5Bytes(data)
	ext := ""
	if strs.IsNotBlank(contentType) {
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil && len(exts) > 0 {
			ext = exts[0]
		}
	}

	// 本地存储使用配置的目录，默认为uploads
	uploadPath := "uploads"
	if cfg != nil && cfg.Local.UploadPath != "" {
		uploadPath = cfg.Local.UploadPath
	}
	
	// 使用path.Join而不是filepath.Join，确保在所有平台上都使用正斜杠
	// 这样生成的URL在浏览器中可以正确解析
	return path.Join(uploadPath, "images", dates.Format(time.Now(), "2006/01/02"), md5+ext)
}

// generateLocalFileKey 生成本地存储的任意文件Key
func generateLocalFileKey(cfg *dto.UploadConfig, data []byte, contentType string) string {
	md5 := digests.MD5Bytes(data)
	ext := ""
	if strs.IsNotBlank(contentType) {
		exts, err := mime.ExtensionsByType(contentType)
		if err == nil && len(exts) > 0 {
			ext = exts[0]
		}
	}
	
	// 根据内容类型确定存储目录
	dir := "files"
	if strings.HasPrefix(contentType, "image/") {
		dir = "images"
	}
	
	// 本地存储使用配置的目录，默认为uploads
	uploadPath := "uploads"
	if cfg != nil && cfg.Local.UploadPath != "" {
		uploadPath = cfg.Local.UploadPath
	}
	
	// 使用path.Join而不是filepath.Join，确保在所有平台上都使用正斜杠
	// 这样生成的URL在浏览器中可以正确解析
	return path.Join(uploadPath, dir, dates.Format(time.Now(), "2006/01/02"), md5+ext)
}
