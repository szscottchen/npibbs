package bbsurls

import (
	"log/slog"
	"net/url"
	"strconv"
	"strings"

	"bbs-go/internal/pkg/config"
)

// 是否是内部链接
func IsInternalUrl(href string) bool {
	if IsAnchor(href) {
		return true
	}
	u, err := url.Parse(config.Instance.BaseURL)
	if err != nil {
		slog.Error(err.Error(), slog.Any("err", err))
		return false
	}
	return strings.Contains(href, u.Host)
}

// 是否是锚链接
func IsAnchor(href string) bool {
	return strings.Index(href, "#") == 0
}

func AbsUrl(path string) string {
	if config.Instance.BaseURL == "/" {
		return path
	}
	return config.Instance.BaseURL + path
}

// 用户主页
func UserUrl(userId int64) string {
	return AbsUrl("/user/" + strconv.FormatInt(userId, 10))
}

// 文章详情
func ArticleUrl(articleId int64) string {
	return AbsUrl("/article/" + strconv.FormatInt(articleId, 10))
}

// 标签文章列表
func TagArticlesUrl(tagId int64) string {
	return AbsUrl("/articles/" + strconv.FormatInt(tagId, 10))
}

// 话题详情
func TopicUrl(topicId int64) string {
	return AbsUrl("/topic/" + strconv.FormatInt(topicId, 10))
}

func UrlJoin(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	
	sep := "/"
	var ss []string
	
	for i, part := range parts {
		part = strings.TrimSpace(part)
		
		// Skip empty parts
		if part == "" {
			continue
		}
		
		// Special handling for the first part (usually BaseURL)
		if i == 0 {
			// If it's root path, add an empty string to represent leading slash correctly
			if part == "/" {
				ss = append(ss, "")
			} else {
				// Remove trailing slash if present (but keep the part if it's more than just a slash)
				if strings.HasSuffix(part, sep) && len(part) > 1 {
					part = part[:len(part)-1]
				}
				ss = append(ss, part)
			}
		} else {
			// For subsequent parts, remove leading and trailing slashes
			from, to := 0, len(part)
			if strings.HasPrefix(part, sep) {
				from = 1
			}
			if strings.HasSuffix(part, sep) {
				to = len(part) - 1
			}
			if from < to {
				part = part[from:to]
				if part != "" {
					ss = append(ss, part)
				}
			}
		}
	}
	
	// Join all parts with slash
	result := strings.Join(ss, sep)
	
	// Ensure root path is displayed as "/"
	if result == "" {
		return "/"
	}
	
	return result
}

// TestUrlJoin 用于测试UrlJoin函数的行为
func TestUrlJoin() string {
	// 测试案例1: baseURL为"/"的情况
	result1 := UrlJoin("/", "uploads/images/2025/11/17/test.jpg")
	
	// 测试案例2: baseURL为"http://localhost:8082/"的情况
	result2 := UrlJoin("http://localhost:8082/", "uploads/images/2025/11/17/test.jpg")
	
	// 测试案例3: baseURL为"http://localhost:8082"的情况（没有结尾斜杠）
	result3 := UrlJoin("http://localhost:8082", "uploads/images/2025/11/17/test.jpg")
	
	return "Test 1 (/): " + result1 + "\nTest 2 (http://localhost:8082/): " + result2 + "\nTest 3 (http://localhost:8082): " + result3
}
