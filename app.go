package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// App struct
type App struct {
	ctx             context.Context
	db              *DatabaseService
	medicineCrawler *MedicineCrawler // 药品爬虫实例
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		db: NewDatabaseService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	// 初始化数据库
	if err := a.db.InitDB(); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
	}
	// 初始化药品爬虫
	a.medicineCrawler = NewMedicineCrawler(a.db.GetDB())
}

// GetAllPatients 获取所有患者数据
func (a *App) GetAllPatients() ([]Patient, error) {
	var patients []Patient
	result := a.db.GetDB().Order("id desc").Find(&patients)
	return patients, result.Error
}

// GetPatientByID 根据ID获取患者数据
func (a *App) GetPatientByID(id uint) (*Patient, error) {
	var patient Patient
	result := a.db.GetDB().First(&patient, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &patient, nil
}

// AddPatient 添加新患者数据
func (a *App) AddPatient(patient Patient) error {
	// 设置当前时间
	now := time.Now()
	patient.Time = now.Format("15:04:05")
	patient.Date = now.Format("2006-01-02")
	patient.CreatedAt = now
	patient.UpdatedAt = now

	result := a.db.GetDB().Create(&patient)
	return result.Error
}

// UpdatePatient 更新患者数据
func (a *App) UpdatePatient(patient Patient) error {
	patient.UpdatedAt = time.Now()
	result := a.db.GetDB().Save(&patient)
	return result.Error
}

// DeletePatient 删除患者数据
func (a *App) DeletePatient(id uint) error {
	result := a.db.GetDB().Delete(&Patient{}, id)
	return result.Error
}

// SearchPatients 搜索患者数据
func (a *App) SearchPatients(searchType, keyword string) ([]Patient, error) {
	var patients []Patient
	var result *gorm.DB

	switch searchType {
	case "name":
		result = a.db.GetDB().Where("name LIKE ?", "%"+keyword+"%").Find(&patients)
	case "date":
		result = a.db.GetDB().Where("date LIKE ?", "%"+keyword+"%").Find(&patients)
	case "sex":
		result = a.db.GetDB().Where("sex = ?", keyword).Find(&patients)
	case "phone":
		result = a.db.GetDB().Where("phone LIKE ?", "%"+keyword+"%").Find(&patients)
	default:
		return nil, fmt.Errorf("invalid search type")
	}

	return patients, result.Error
}

// GetDatabaseInfo 获取数据库信息
func (a *App) GetDatabaseInfo() (map[string]interface{}, error) {
	// 获取数据库文件大小
	fileInfo, err := os.Stat(filepath.Join(GetAppRuntimePath(), DBName))
	if err != nil {
		return nil, err
	}
	size := fileInfo.Size()

	// 获取记录数量
	var count int64
	a.db.GetDB().Model(&Patient{}).Count(&count)

	// 获取备份信息
	backupDir := "backup"
	var latestBackup string
	if _, err := os.Stat(filepath.Join(GetAppRuntimePath(), backupDir)); err == nil {
		files, _ := os.ReadDir(filepath.Join(GetAppRuntimePath(), backupDir))
		if len(files) > 0 {
			latestBackup = strings.Replace(files[len(files)-1].Name(), "patient-", "", 1)
			latestBackup = strings.Replace(latestBackup, ".db", "", 1)
		}
	}

	return map[string]interface{}{
		"size":   size,
		"count":  count,
		"new_bk": latestBackup,
	}, nil
}

// CreateDatabase 创建数据库表
func (a *App) CreateDatabase() error {
	return a.db.AutoMigrate()
}

// DropDatabase 删除所有数据
func (a *App) DropDatabase() error {
	return a.db.GetDB().Exec("DELETE FROM patient").Error
}

// BackupDatabase 备份数据库
func (a *App) BackupDatabase(clearAfterBackup bool) error {
	// 创建备份目录
	backupDir := "backup"
	backupDir = filepath.Join(GetAppRuntimePath(), backupDir)
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		os.Mkdir(backupDir, 0755)
	}

	// 生成备份文件名
	backupName := fmt.Sprintf("patient-%s.db", time.Now().Format("2006-01-02_15-04-05"))
	backupPath := filepath.Join(backupDir, backupName)

	// 复制数据库文件
	dbPath := filepath.Join(GetAppRuntimePath(), DBName)
	sourceFile, err := os.Open(dbPath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(backupPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// 如果需要清空当前数据库
	if clearAfterBackup {
		return a.DropDatabase()
	}

	return nil
}

// ExportPatientsToExcel 导出患者数据到Excel文件
func (a *App) ExportPatientsToExcel() (string, error) {
	// 获取所有患者数据
	patients, err := a.GetAllPatients()
	if err != nil {
		return "", fmt.Errorf("获取患者数据失败: %v", err)
	}

	// 创建Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// 设置工作表名称
	sheetName := "患者信息"
	f.SetSheetName("Sheet1", sheetName)

	// 设置表头
	headers := []string{
		"ID", "姓名", "性别", "年龄", "身份证号", "电话", "联系方式", "监护人", "工作", "地址", "病史", "患病时间", "医师",
		"诊断", "治疗方案", "医嘱", "费用", "费用备注", "创建时间", "更新时间",
	}

	// 写入表头
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for i, patient := range patients {
		row := i + 2 // 从第二行开始写入数据

		// 格式化时间
		createdAt := patient.CreatedAt.Format("2006-01-02 15:04:05")
		updatedAt := patient.UpdatedAt.Format("2006-01-02 15:04:05")

		data := []interface{}{
			patient.ID,
			patient.Name,
			patient.Sex,
			patient.Age,
			patient.IDCard,
			patient.Phone,
			patient.Contact,
			patient.Parent,
			patient.Work,
			patient.Address,
			patient.AllergyHistory,
			patient.IllTime,
			patient.Doc,
			patient.Detail,
			patient.Solution,
			patient.MedicalAdvice,
			patient.Fee,
			patient.Money,
			createdAt,
			updatedAt,
		}

		for j, value := range data {
			cell, _ := excelize.CoordinatesToCellName(j+1, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 8)  // ID
	f.SetColWidth(sheetName, "B", "B", 12) // 姓名
	f.SetColWidth(sheetName, "C", "C", 8)  // 性别
	f.SetColWidth(sheetName, "D", "D", 8)  // 年龄
	f.SetColWidth(sheetName, "E", "E", 20) // 身份证号
	f.SetColWidth(sheetName, "F", "F", 15) // 电话
	f.SetColWidth(sheetName, "G", "G", 20) // 联系方式
	f.SetColWidth(sheetName, "H", "H", 12) // 监护人
	f.SetColWidth(sheetName, "I", "I", 15) // 工作
	f.SetColWidth(sheetName, "J", "J", 20) // 地址
	f.SetColWidth(sheetName, "K", "K", 25) // 病史
	f.SetColWidth(sheetName, "L", "L", 15) // 患病时间
	f.SetColWidth(sheetName, "M", "M", 12) // 医师
	f.SetColWidth(sheetName, "N", "N", 25) // 诊断
	f.SetColWidth(sheetName, "O", "O", 25) // 治疗方案
	f.SetColWidth(sheetName, "P", "P", 25) // 医嘱
	f.SetColWidth(sheetName, "Q", "Q", 12) // 费用
	f.SetColWidth(sheetName, "R", "R", 15) // 费用备注
	f.SetColWidth(sheetName, "S", "S", 20) // 创建时间
	f.SetColWidth(sheetName, "T", "T", 20) // 更新时间

	// 获取桌面路径
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户目录失败: %v", err)
	}
	desktopPath := filepath.Join(homeDir, "Desktop")

	// 生成文件名
	now := time.Now()
	filename := fmt.Sprintf("patients-%s.xlsx", now.Format("20060102-150405"))
	filePath := filepath.Join(desktopPath, filename)

	// 保存文件
	if err := f.SaveAs(filePath); err != nil {
		return "", fmt.Errorf("保存Excel文件失败: %v", err)
	}

	return filePath, nil
}

// GetColumnConfigs 获取患者表格列配置
func (a *App) GetColumnConfigs() ([]ColumnConfig, error) {
	var configs []ColumnConfig
	result := a.db.GetDB().Where("table_name = ?", "patient").Order("sort_order asc").Find(&configs)
	return configs, result.Error
}

// SaveColumnConfigs 保存患者表格列配置
func (a *App) SaveColumnConfigs(configs []ColumnConfig) error {
	// 开启事务
	tx := a.db.GetDB().Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 遍历配置，根据columnKey更新或新增
	for i, config := range configs {
		config.Name = "patient"
		config.SortOrder = i

		// 查找是否存在该columnKey的配置
		var existingConfig ColumnConfig
		err := tx.Where("table_name = ? AND column_key = ?", "patient", config.ColumnKey).First(&existingConfig).Error

		if err == nil {
			// 存在则更新
			if err := tx.Model(&existingConfig).Updates(map[string]interface{}{
				"visible":    config.Visible,
				"sort_order": config.SortOrder,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 不存在则新增
			if err := tx.Model(&ColumnConfig{}).Create(&config).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit().Error
}

// InitDefaultColumnConfigs 初始化患者表默认列配置
func (a *App) InitDefaultColumnConfigs() error {
	// 检查是否已有配置
	var count int64
	a.db.GetDB().Model(&ColumnConfig{}).Where("table_name = ?", "patient").Count(&count)
	if count > 0 {
		return nil // 已有配置，不需要初始化
	}

	// 患者表的默认列配置
	defaultConfigs := []ColumnConfig{
		{Name: "patient", ColumnKey: "id", Visible: true, SortOrder: 0},
		{Name: "patient", ColumnKey: "name", Visible: true, SortOrder: 1},
		{Name: "patient", ColumnKey: "sex", Visible: true, SortOrder: 2},
		{Name: "patient", ColumnKey: "age", Visible: true, SortOrder: 3},
		{Name: "patient", ColumnKey: "phone", Visible: true, SortOrder: 4},
		{Name: "patient", ColumnKey: "address", Visible: true, SortOrder: 5},
		{Name: "patient", ColumnKey: "parent", Visible: true, SortOrder: 6},
		{Name: "patient", ColumnKey: "work", Visible: true, SortOrder: 7},
		{Name: "patient", ColumnKey: "detail", Visible: true, SortOrder: 8},
		{Name: "patient", ColumnKey: "contact", Visible: false, SortOrder: 9},
		{Name: "patient", ColumnKey: "date", Visible: false, SortOrder: 10},
		{Name: "patient", ColumnKey: "ill_time", Visible: false, SortOrder: 11},
		{Name: "patient", ColumnKey: "doc", Visible: false, SortOrder: 12},
		{Name: "patient", ColumnKey: "solution", Visible: false, SortOrder: 13},
		{Name: "patient", ColumnKey: "medical_advice", Visible: false, SortOrder: 14},
		{Name: "patient", ColumnKey: "fee", Visible: false, SortOrder: 15},
		{Name: "patient", ColumnKey: "money", Visible: false, SortOrder: 16},
	}

	for _, config := range defaultConfigs {
		if err := a.db.GetDB().Create(&config).Error; err != nil {
			return err
		}
	}

	return nil
}

// CrawlMedicineData 开始爬取药品数据
func (a *App) CrawlMedicineData() error {
	if a.medicineCrawler == nil {
		return fmt.Errorf("药品爬虫未初始化")
	}
	return a.medicineCrawler.CrawlMedicineData()
}

// CrawlMedicineDataWithPages 分页爬取药品数据
func (a *App) CrawlMedicineDataWithPages(maxPages int) error {
	if a.medicineCrawler == nil {
		return fmt.Errorf("药品爬虫未初始化")
	}
	if maxPages <= 0 {
		maxPages = 5 // 默认爬取5页
	}
	return a.medicineCrawler.CrawlMedicineDataWithPagination(maxPages)
}

// GetMedicineCount 获取数据库中药品数量
func (a *App) GetMedicineCount() (int64, error) {
	if a.medicineCrawler == nil {
		return 0, fmt.Errorf("药品爬虫未初始化")
	}
	return a.medicineCrawler.GetMedicineCount()
}

// ClearAllMedicines 清空所有药品数据
func (a *App) ClearAllMedicines() error {
	result := a.db.GetDB().Where("1 = 1").Delete(&Medicine{})
	return result.Error
}

// GetConfig 获取配置
func (a *App) GetConfig(key string) (string, error) {
	var config Config
	result := a.db.GetDB().Where("key = ?", key).First(&config)
	if result.Error != nil {
		return "", result.Error
	}
	return config.Value, nil
}

// SetConfig 设置配置
func (a *App) SetConfig(key, value string) error {
	var config Config
	result := a.db.GetDB().Where("key = ?", key).First(&config)

	if result.Error != nil {
		// 配置不存在，创建新配置
		config = Config{
			Key:       key,
			Value:     value,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return a.db.GetDB().Create(&config).Error
	} else {
		// 配置存在，更新配置
		config.Value = value
		config.UpdatedAt = time.Now()
		return a.db.GetDB().Model(&config).Where("key = ?", key).Updates(map[string]interface{}{
			"value":      value,
			"updated_at": time.Now(),
		}).Error
	}
}

// GetAllConfigs 获取所有配置
func (a *App) GetAllConfigs() (map[string]string, error) {
	var configs []Config
	result := a.db.GetDB().Find(&configs)
	if result.Error != nil {
		return nil, result.Error
	}

	configMap := make(map[string]string)
	for _, config := range configs {
		configMap[config.Key] = config.Value
	}

	return configMap, nil
}

// SelectDatabaseFile 选择数据库文件
func (a *App) SelectDatabaseFile() (string, error) {
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择数据库文件",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "SQLite数据库文件 (*.db)",
				Pattern:     "*.db",
			},
			{
				DisplayName: "所有文件 (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil {
		return "", err
	}

	// 如果用户选择了文件，更新数据库连接
	if filePath != "" {
		// 这里可以添加切换数据库的逻辑
		// 目前只返回路径，实际切换数据库需要重新初始化数据库服务
	}

	return filePath, nil
}

// GetDatabasePath 获取数据库路径
func (a *App) GetDatabasePath() (string, error) {
	// 获取应用程序数据目录
	appDataDir := os.Getenv("APPDATA")
	if appDataDir == "" {
		appDataDir = "."
	}

	// 获取可执行文件名称
	execPath, err := os.Executable()
	var appName string
	if err != nil {
		appName = "medical-form"
	} else {
		appName = filepath.Base(execPath)
		if ext := filepath.Ext(appName); ext != "" {
			appName = appName[:len(appName)-len(ext)]
		}
	}

	// 构建数据库文件路径
	dbPath := filepath.Join(appDataDir, appName, DBName)

	// 检查文件是否存在并获取大小
	fileInfo, err := os.Stat(dbPath)
	if err != nil {
		return dbPath, err
	}

	return fmt.Sprintf("%s (%.2f KB)", dbPath, float64(fileInfo.Size())/1024), nil
}

// GetAllDoctors 获取所有医师数据
func (a *App) GetAllDoctors() ([]Doctor, error) {
	var doctors []Doctor
	result := a.db.GetDB().Order("id desc").Find(&doctors)
	return doctors, result.Error
}

// GetDoctorByID 根据ID获取医师数据
func (a *App) GetDoctorByID(id uint) (*Doctor, error) {
	var doctor Doctor
	result := a.db.GetDB().First(&doctor, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &doctor, nil
}

// AddDoctor 添加医师
func (a *App) AddDoctor(doctor Doctor) error {
	now := time.Now()
	doctor.CreatedAt = now
	doctor.UpdatedAt = now
	result := a.db.GetDB().Create(&doctor)
	return result.Error
}

// UpdateDoctor 更新医师信息
func (a *App) UpdateDoctor(doctor Doctor) error {
	doctor.UpdatedAt = time.Now()
	result := a.db.GetDB().Save(&doctor)
	return result.Error
}

// DeleteDoctor 删除医师
func (a *App) DeleteDoctor(id uint) error {
	result := a.db.GetDB().Delete(&Doctor{}, id)
	return result.Error
}

// SearchDoctors 搜索医师
func (a *App) SearchDoctors(searchType, keyword string) ([]Doctor, error) {
	var doctors []Doctor
	db := a.db.GetDB()

	switch searchType {
	case "name":
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	case "level":
		db = db.Where("level LIKE ?", "%"+keyword+"%")
	case "phone":
		db = db.Where("phone LIKE ?", "%"+keyword+"%")
	default:
		// 全字段搜索
		db = db.Where("name LIKE ? OR level LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	result := db.Order("id desc").Find(&doctors)
	return doctors, result.Error
}

// GetAppInfo 获取应用程序信息
func (a *App) GetAppInfo() (map[string]interface{}, error) {
	// 获取可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		execPath = "未知"
	}

	// 获取应用程序名称
	var appName string
	if execPath != "未知" {
		appName = filepath.Base(execPath)
		if ext := filepath.Ext(appName); ext != "" {
			appName = appName[:len(appName)-len(ext)]
		}
	} else {
		appName = "medical-form"
	}

	// 获取数据目录
	appDataDir := os.Getenv("APPDATA")
	if appDataDir == "" {
		appDataDir = "."
	}
	dataDir := filepath.Join(appDataDir, appName)

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "未知"
	}

	// 获取环境信息（如果有runtime context）
	var envInfo map[string]interface{}
	if a.ctx != nil {
		env := runtime.Environment(a.ctx)
		envInfo = map[string]interface{}{
			"platform":  env.Platform,
			"arch":      env.Arch,
			"buildType": env.BuildType,
		}
	}

	return map[string]interface{}{
		"executablePath": execPath,
		"appName":        appName,
		"dataDirectory":  dataDir,
		"workingDir":     currentDir,
		"environment":    envInfo,
	}, nil
}

// GetMetaInfo 获取程序元信息
func (a *App) GetMetaInfo() map[string]string {
	return map[string]string{
		"title":   Title,   // 程序标题
		"version": Version, // 程序版本
	}
}

// GetAllTemplates 获取所有模板数据
func (a *App) GetAllTemplates() ([]Template, error) {
	var templates []Template
	result := a.db.GetDB().Order("type, id desc").Find(&templates)
	return templates, result.Error
}

// CreateTemplate 创建新模板
func (a *App) CreateTemplate(template Template) (*Template, error) {
	// 检查名称是否重复
	var existingTemplate Template
	result := a.db.GetDB().Where("name = ?", template.Name).First(&existingTemplate)
	if result.Error == nil {
		return nil, fmt.Errorf("模板名称已存在")
	}
	if result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	// 创建新模板
	result = a.db.GetDB().Create(&template)
	if result.Error != nil {
		return nil, result.Error
	}
	return &template, nil
}

// UpdateTemplate 更新模板
func (a *App) UpdateTemplate(template Template) error {
	// 检查名称是否与其他模板重复
	var existingTemplate Template
	result := a.db.GetDB().Where("name = ? AND id != ?", template.Name, template.ID).First(&existingTemplate)
	if result.Error == nil {
		return fmt.Errorf("模板名称已存在")
	}
	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	// 更新模板
	result = a.db.GetDB().Save(&template)
	return result.Error
}

// DeleteTemplate 删除模板
func (a *App) DeleteTemplate(id uint) error {
	result := a.db.GetDB().Delete(&Template{}, id)
	return result.Error
}

// GetTemplatesByType 根据类型获取模板
func (a *App) GetTemplatesByType(templateType string) ([]Template, error) {
	var templates []Template
	result := a.db.GetDB().Where("type = ?", templateType).Find(&templates)
	return templates, result.Error
}

// GetAllMedicines 获取所有药品数据
func (a *App) GetAllMedicines() ([]Medicine, error) {
	var medicines []Medicine
	result := a.db.GetDB().Order("id desc").Find(&medicines)
	return medicines, result.Error
}

// GetMedicineByID 根据ID获取药品数据
func (a *App) GetMedicineByID(id uint) (*Medicine, error) {
	var medicine Medicine
	result := a.db.GetDB().First(&medicine, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &medicine, nil
}

// CreateMedicine 创建新药品
func (a *App) CreateMedicine(medicine Medicine) (*Medicine, error) {
	// 设置创建时间
	medicine.CreatedAt = time.Now()
	medicine.UpdatedAt = time.Now()

	result := a.db.GetDB().Create(&medicine)
	if result.Error != nil {
		return nil, result.Error
	}
	return &medicine, nil
}

// UpdateMedicine 更新药品信息
func (a *App) UpdateMedicine(medicine Medicine) error {
	// 设置更新时间
	medicine.UpdatedAt = time.Now()

	result := a.db.GetDB().Save(&medicine)
	return result.Error
}

// DeleteMedicine 删除药品
func (a *App) DeleteMedicine(id uint) error {
	result := a.db.GetDB().Delete(&Medicine{}, id)
	return result.Error
}

// SearchMedicines 搜索药品
func (a *App) SearchMedicines(searchType, keyword string) ([]Medicine, error) {
	var medicines []Medicine
	var result *gorm.DB

	switch searchType {
	case "name":
		result = a.db.GetDB().Where("name LIKE ?", "%"+keyword+"%").Find(&medicines)
	case "manufacturer":
		result = a.db.GetDB().Where("manufacturer LIKE ?", "%"+keyword+"%").Find(&medicines)
	default:
		// 默认搜索药品名称和生产厂家
		result = a.db.GetDB().Where("name LIKE ? OR manufacturer LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&medicines)
	}

	return medicines, result.Error
}

// MigrateOldDatabase 迁移旧版本数据库
func (a *App) MigrateOldDatabase(oldDbPath string) error {
	// 检查旧数据库文件是否存在
	if _, err := os.Stat(oldDbPath); os.IsNotExist(err) {
		return fmt.Errorf("旧数据库文件不存在: %s", oldDbPath)
	}

	// 创建临时数据库服务连接旧数据库
	oldDbService := NewDatabaseService()
	oldDbService.dbPath = oldDbPath

	// 连接旧数据库
	if err := oldDbService.connectDB(); err != nil {
		return fmt.Errorf("连接旧数据库失败: %v", err)
	}
	defer oldDbService.CloseDB()

	// 检查旧数据库是否包含旧表结构
	var tableExists string
	err := oldDbService.db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name='pig'").Scan(&tableExists).Error
	if err != nil {
		return fmt.Errorf("检查旧表结构失败: %v", err.Error())
	}

	// 如果没有旧表结构，尝试直接从pig表读取
	var oldPatients []OldPigTable
	if tableExists != "" {
		// 从旧表结构读取数据
		result := oldDbService.db.Table("pig").Find(&oldPatients)
		if result.Error != nil {
			return fmt.Errorf("读取旧数据失败: %v", result.Error)
		}
	} else {
		return fmt.Errorf("旧数据库没有pig表")
	}

	// 转换并插入新数据库
	for _, oldPatient := range oldPatients {
		// 转换旧数据结构到新数据结构
		newPatient := Patient{
			Date:           oldPatient.Date,
			Time:           oldPatient.Time,
			Name:           oldPatient.Name,
			Sex:            oldPatient.Sex,
			Age:            oldPatient.Age,
			IllTime:        oldPatient.IllTime,
			Phone:          oldPatient.Phone,
			Address:        oldPatient.Address,
			Parent:         oldPatient.Parent,
			Work:           oldPatient.Work,
			AllergyHistory: "", // 旧版本没有过敏史字段
			Detail:         oldPatient.Detail,
			Solution:       oldPatient.Solution,
			MedicalAdvice:  oldPatient.MedicalAdvice,
			Addon:          oldPatient.Addon,
			Money:          oldPatient.Money,
			Doc:            oldPatient.Doc,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// 插入到新数据库
		result := a.db.GetDB().Create(&newPatient)
		if result.Error != nil {
			return fmt.Errorf("插入迁移数据失败: %v", result.Error)
		}
	}

	return nil
}
