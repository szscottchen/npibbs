package dto

// 积分配置
type ScoreConfig struct {
	PostTopicScore   int          `json:"postTopicScore"`   // 发帖获得积分
	PostCommentScore int          `json:"postCommentScore"` // 跟帖获得积分
	CheckInScore     int          `json:"checkInScore"`     // 签到积分
	ValueTypes       []ValueType  `json:"valueTypes"`       // 价值类型配置
}

// 价值类型
type ValueType struct {
	Label string `json:"label"` // 标签
	Score int    `json:"score"` // 积分
}

// SysConfigResponse
//
//	配置返回结构体
type SysConfigResponse struct {
	SiteTitle                  string        `json:"siteTitle"`
	SiteDescription            string        `json:"siteDescription"`
	SiteKeywords               []string      `json:"siteKeywords"`
	SiteLogo                   string        `json:"siteLogo"`
	SiteNavs                   []ActionLink  `json:"siteNavs"`
	SiteNotification           string        `json:"siteNotification"`
	RecommendTags              []string      `json:"recommendTags"`
	UrlRedirect                bool          `json:"urlRedirect"`
	ScoreConfig                ScoreConfig   `json:"scoreConfig"`
	DefaultNodeId              int64         `json:"defaultNodeId"`
	ArticlePending             bool          `json:"articlePending"`
	TopicCaptcha               bool          `json:"topicCaptcha"`
	UserObserveSeconds         int           `json:"userObserveSeconds"`
	TokenExpireDays            int           `json:"tokenExpireDays"`
	CreateTopicEmailVerified   bool          `json:"createTopicEmailVerified"`
	CreateArticleEmailVerified bool          `json:"createArticleEmailVerified"`
	CreateCommentEmailVerified bool          `json:"createCommentEmailVerified"`
	EnableHideContent          bool          `json:"enableHideContent"`
	Modules                    ModulesConfig `json:"modules"`
	EmailWhitelist             []string      `json:"emailWhitelist"` // 邮箱白名单
	UploadConfig               *UploadConfig `json:"uploadConfig"`   // 上传配置
	ValueTypes                 []ValueType   `json:"valueTypes"`     // 价值类型配置
}

// ModulesConfig
//
//	模块配置
type ModulesConfig struct {
	Tweet   bool `json:"tweet"`
	Topic   bool `json:"topic"`
	Article bool `json:"article"`
}

type AliyunSmsConfig struct {
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	SignName        string `json:"signName"`
	TemplateCode    string `json:"templateCode"`
}

type UploadMethod string

const (
	AliyunOss  UploadMethod = "AliyunOss"
	TencentCos UploadMethod = "TencentCos"
	Local      UploadMethod = "Local"
)

type UploadConfig struct {
	EnableUploadMethod UploadMethod           `json:"enableUploadMethod"`
	AliyunOss          AliyunOssUploadConfig  `json:"aliyunOss"`
	TencentCos         TencentCosUploadConfig `json:"tencentCos"`
	Local              LocalUploadConfig      `json:"local"`
}

type AliyunOssUploadConfig struct {
	Host            string `json:"host"`
	Bucket          string `json:"bucket"`
	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"accessKeyId"`
	AccessKeySecret string `json:"accessKeySecret"`
	StyleSplitter   string `json:"styleSplitter"`
	StyleAvatar     string `json:"styleAvatar"`
	StylePreview    string `json:"stylePreview"`
	StyleSmall      string `json:"styleSmall"`
	StyleDetail     string `json:"styleDetail"`
}

type TencentCosUploadConfig struct {
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
}

type LocalUploadConfig struct {
	UploadPath string `json:"uploadPath"` // 上传目录，默认值：uploads
	MaxSizeMB  int    `json:"maxSizeMB"`  // 文件大小限制（MB），默认值：10
}
