package models

// BatchImportUserRow Excel中每一行的数据
type BatchImportUserRow struct {
	RowNumber     int    `json:"rowNumber"`     // Excel行号（用于错误定位）
	Phone         string `json:"phone"`         // 电话
	Username      string `json:"username"`      // 用户名
	Email         string `json:"email"`         // 邮箱
	EmployeeID    string `json:"employeeId"`    // 员工编号
	BDivision     string `json:"bDivision"`     // 业务部门
	BDepartment   string `json:"bDepartment"`   // 业务单位
	JobPosition   string `json:"jobPosition"`   // 工作岗位
	CommunityRole string `json:"communityRole"` // 社区角色
	Nickname      string `json:"nickname"`      // 昵称
}

// BatchImportResult 批量导入结果
type BatchImportResult struct {
	TotalCount   int                `json:"totalCount"`   // 总条数
	SuccessCount int                `json:"successCount"` // 成功条数
	FailCount    int                `json:"failCount"`    // 失败条数
	Errors       []BatchImportError `json:"errors"`       // 错误详情
}

// BatchImportError 单条错误信息
type BatchImportError struct {
	RowNumber    int    `json:"rowNumber"`    // 行号
	EmployeeID   string `json:"employeeId"`   // 员工编号
	Nickname     string `json:"nickname"`     // 昵称
	ErrorMessage string `json:"errorMessage"` // 错误信息
	ErrorTime    int64  `json:"errorTime"`    // 错误发生时间
}
