package services

import (
	"bbs-go/internal/models"
	"bbs-go/internal/models/constants"
	"bbs-go/internal/pkg/validate"
	"bbs-go/internal/repositories"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mlogclub/simple/common/passwd"
	"github.com/mlogclub/simple/sqls"
	"github.com/xuri/excelize/v2"
)

var UserBatchImportService = newUserBatchImportService()

func newUserBatchImportService() *userBatchImportService {
	return &userBatchImportService{}
}

type userBatchImportService struct{}

// ImportUsersFromExcel 从Excel文件导入用户
func (s *userBatchImportService) ImportUsersFromExcel(filePath string) (*models.BatchImportResult, error) {
	// 1. 打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}


	defer f.Close()

	// 2. 读取第一个Sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("Excel文件中没有工作表")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("读取Excel数据失败: %v", err)
	}

	// 添加调试信息
	slog.Info("Excel文件读取调试信息",
		slog.Any("总行数", len(rows)),
		slog.Any("工作表名称", sheets[0]),
	)

	// 输出每行的内容用于调试
	for i, row := range rows {
		slog.Info("Excel行数据",
			slog.Any("行号", i+1),
			slog.Any("列数", len(row)),
			slog.Any("内容", row),
		)
	}

	if len(rows) <= 1 {
		return nil, fmt.Errorf("Excel文件中没有数据行，实际读取到%d行", len(rows))
	}

	// 3. 初始化结果
	result := &models.BatchImportResult{
		TotalCount:   len(rows) - 1, // 减去表头
		SuccessCount: 0,
		FailCount:    0,
		Errors:       make([]models.BatchImportError, 0),
	}

	// 4. 初始化错误日志文件
	errorLogPath := s.getErrorLogPath()

	// 5. 跳过表头，逐行处理
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNumber := i + 1 // Excel行号从1开始

		// 解析行数据
		importRow := s.parseRow(row, rowNumber)

		// 创建用户
		err := s.createUserFromImportRow(importRow)
		if err != nil {
			// 记录错误
			result.FailCount++
			importErr := models.BatchImportError{
				RowNumber:    rowNumber,
				EmployeeID:   importRow.EmployeeID,
				Nickname:     importRow.Nickname,
				ErrorMessage: err.Error(),
				ErrorTime:    time.Now().Unix(),
			}
			result.Errors = append(result.Errors, importErr)

			// 写入错误日志文件
			s.writeErrorLog(errorLogPath, importErr)
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// parseRow 解析Excel行数据
func (s *userBatchImportService) parseRow(row []string, rowNumber int) *models.BatchImportUserRow {
	// 确保有足够的列
	for len(row) < 9 {
		row = append(row, "")
	}

	return &models.BatchImportUserRow{
		RowNumber:     rowNumber,
		Phone:         strings.TrimSpace(row[0]),
		Username:      strings.TrimSpace(row[1]),
		Email:         strings.TrimSpace(row[2]),
		EmployeeID:    strings.TrimSpace(row[3]),
		BDivision:     strings.TrimSpace(row[4]),
		BDepartment:   strings.TrimSpace(row[5]),
		JobPosition:   strings.TrimSpace(row[6]),
		CommunityRole: strings.TrimSpace(row[7]),
		Nickname:      strings.TrimSpace(row[8]),
	}
}

// createUserFromImportRow 从导入行创建用户
func (s *userBatchImportService) createUserFromImportRow(row *models.BatchImportUserRow) error {
	// 1. 数据验证
	if err := s.validateImportRow(row); err != nil {
		return err
	}

	// 2. 构建User对象
	user := &models.User{
		Type:          1, // 1表示员工
		Status:        constants.StatusOk,
		Password:      passwd.EncodePassword("123456"), // 默认密码
		Nickname:      row.Nickname,
		EmployeeID:    row.EmployeeID,
		BDivision:     row.BDivision,
		BDepartment:   row.BDepartment,
		Jobposition:   row.JobPosition,
		CommunityRole: row.CommunityRole,
		CreateTime:    time.Now().Unix(),
		UpdateTime:    time.Now().Unix(),
	}

	// 设置邮箱
	if len(row.Email) > 0 {
		user.Email = sql.NullString{String: row.Email, Valid: true}
	}

	// 设置用户名
	if len(row.Username) > 0 {
		user.Username = sql.NullString{String: row.Username, Valid: true}
	}

	// 设置手机号
	if len(row.Phone) > 0 {
		user.Phone = sql.NullString{String: row.Phone, Valid: true}
	}

	// 3. 创建用户
	err := repositories.UserRepository.Create(sqls.DB(), user)
	if err != nil {
		return fmt.Errorf("创建用户失败: %v", err)
	}

	// 4. 为新用户分配默认"user"角色
	// 查找"user"角色
	userRole := repositories.RoleRepository.FindOne(sqls.DB(), sqls.NewCnd().Eq("code", "user"))
	if userRole != nil {
		// 为用户分配角色
		err = UserRoleService.UpdateUserRoles(user.Id, []int64{userRole.Id})
		if err != nil {
			// 记录错误但不中断流程
			slog.Error("为用户分配角色失败", slog.String("error", err.Error()), slog.Int64("userId", user.Id))
		}
	} else {
		slog.Warn("未找到'user'角色，跳过角色分配")
	}

	return nil
}

// validateImportRow 验证导入行数据
func (s *userBatchImportService) validateImportRow(row *models.BatchImportUserRow) error {
	// 1. 昵称验证
	if len(row.Nickname) == 0 {
		return fmt.Errorf("昵称不能为空")
	}

	// 2. 邮箱验证
	if len(row.Email) == 0 {
		return fmt.Errorf("邮箱不能为空")
	}

	// 邮箱格式验证
	if err := validate.IsEmail(row.Email); err != nil {
		return fmt.Errorf("邮箱格式不正确: %v", err)
	}

	// 检查邮箱是否已存在
	if UserService.GetByEmail(row.Email) != nil {
		return fmt.Errorf("邮箱：%s 已被占用", row.Email)
	}

	// 3. 用户名验证（如果提供了用户名）
	if len(row.Username) > 0 {
		// 检查用户名格式
		if err := validate.IsUsername(row.Username); err != nil {
			return fmt.Errorf("用户名格式不正确: %v", err)
		}
		// 检查用户名是否已存在
		if UserService.GetByUsername(row.Username) != nil {
			return fmt.Errorf("用户名：%s 已被占用", row.Username)
		}
	}

	// 4. 手机号验证（如果提供了手机号）
	if len(row.Phone) > 0 {
		// 检查手机号是否已存在
		if UserService.GetByPhone(row.Phone) != nil {
			return fmt.Errorf("手机号：%s 已被占用", row.Phone)
		}
	}

	return nil
}

// getErrorLogPath 获取错误日志文件路径
func (s *userBatchImportService) getErrorLogPath() string {
	dateStr := time.Now().Format("20060102")
	logDir := "./logs"

	// 确保logs目录存在
	os.MkdirAll(logDir, 0755)

	return filepath.Join(logDir, fmt.Sprintf("batchinputerr_%s.txt", dateStr))
}

// writeErrorLog 写入错误日志
func (s *userBatchImportService) writeErrorLog(logPath string, err models.BatchImportError) {
	// 格式：employee_id,nickname,错误信息,错误时间
	timeStr := time.Unix(err.ErrorTime, 0).Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("%s,%s,%s,%s\n",
		err.EmployeeID,
		err.Nickname,
		err.ErrorMessage,
		timeStr,
	)

	// 追加写入文件
	f, e := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		// 如果无法写入日志，至少记录到控制台
		fmt.Printf("无法写入错误日志: %v\n", e)
		return
	}
	defer f.Close()

	f.WriteString(logLine)
}

// GenerateTemplate 生成Excel模板文件
func (s *userBatchImportService) GenerateTemplate(filePath string) error {
	f := excelize.NewFile()
	defer f.Close()

	// 创建Sheet，使用默认的"Sheet1"名称以确保导入时能正确读取
	sheetName := "Sheet1"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}

	// 设置表头
	headers := []string{
		"手机号(phone)",
		"用户名(username)",
		"邮箱(email)*必填",
		"员工编号(employee_id)",
		"业务部门(b_division)",
		"业务单位(b_department)",
		"工作岗位(job_position)",
		"社区角色(community_role)",
		"昵称(nickname)*必填",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// 添加示例数据
	exampleData := []interface{}{
		"13800138000",
		"zhangsan",
		"zhangsan@example.com",
		"E001",
		"技术部",
		"研发中心",
		"高级工程师",
		"技术专家",
		"张三",
	}

	for i, data := range exampleData {
		cell := fmt.Sprintf("%c2", 'A'+i)
		f.SetCellValue(sheetName, cell, data)
	}

	// 设置活动Sheet
	f.SetActiveSheet(index)

	// 保存文件
	if err := f.SaveAs(filePath); err != nil {
		return err
	}

	return nil
}
