package render

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/dto"
	"bbs-go/internal/services"
	"log/slog"

	"github.com/mlogclub/simple/common/jsons"
	"github.com/mlogclub/simple/common/strs"
)

func BuildImageList(imageListStr string) (imageList []models.ImageInfo) {
	if strs.IsNotBlank(imageListStr) {
		var images []models.ImageDTO
		if err := jsons.Parse(imageListStr, &images); err == nil {
			if len(images) > 0 {
				for _, image := range images {
					// 在本地上传模式下，直接使用image.Url作为Url和Preview
					// 在OSS模式下，使用HandleOssImageStyle处理图片URL
					cfg := services.SysConfigService.GetUploadConfig()
					if cfg.EnableUploadMethod == dto.AliyunOss || cfg.EnableUploadMethod == dto.TencentCos {
						imageList = append(imageList, models.ImageInfo{
							Url:      HandleOssImageStyleDetail(image.Url),
							Preview:  HandleOssImageStylePreview(image.Url),
							Filename: image.Filename,
						})
					} else {
						imageList = append(imageList, models.ImageInfo{
							Url:      image.Url,
							Preview:  image.Url,
							Filename: image.Filename,
						})
					}
				}
			}
		} else {
			slog.Error(err.Error(), slog.Any("err", err))
		}
	}
	return
}

func BuildImage(imageStr string) *models.ImageInfo {
	if strs.IsBlank(imageStr) {
		return nil
	}
	var img *models.ImageDTO
	if err := jsons.Parse(imageStr, &img); err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return nil
	} else {
		// 在本地上传模式下，直接使用img.Url作为Url和Preview
		// 在OSS模式下，使用HandleOssImageStyle处理图片URL
		cfg := services.SysConfigService.GetUploadConfig()
		if cfg.EnableUploadMethod == dto.AliyunOss || cfg.EnableUploadMethod == dto.TencentCos {
			return &models.ImageInfo{
				Url:      HandleOssImageStyleDetail(img.Url),
				Preview:  HandleOssImageStylePreview(img.Url),
				Filename: img.Filename,
			}
		} else {
			return &models.ImageInfo{
				Url:      img.Url,
				Preview:  img.Url,
				Filename: img.Filename,
			}
		}
	}
}
