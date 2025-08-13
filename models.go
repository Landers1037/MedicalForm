package main

import (
	"time"
)

// Patient 患者信息模型
type Patient struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	// 兼容性字段
	Date string `gorm:"size:50" json:"date"` // 日期
	// 兼容性字段
	Time string `gorm:"size:50" json:"time"` // 准确时间
	Name string `gorm:"size:20" json:"name"` // 姓名
	Sex  string `gorm:"size:10" json:"sex"`  // 性别
	Age  string `gorm:"size:10" json:"age"`  // 年龄
	// 兼容性字段
	IllTime        string    `gorm:"size:100" json:"ill_time"`         // 患病时间
	Phone          string    `gorm:"size:20" json:"phone"`             // 电话
	Contact        string    `gorm:"type:text" json:"contact"`         // 联系方式
	Parent         string    `gorm:"size:20" json:"parent"`            // 家长
	Work           string    `gorm:"size:255" json:"work"`             // 工作
	Address        string    `gorm:"size:255" json:"address"`          // 地址
	AllergyHistory string    `gorm:"type:text" json:"allergy_history"` // 过敏史
	Detail         string    `gorm:"type:text" json:"detail"`          // 详情
	Solution       string    `gorm:"type:text" json:"solution"`        // 解决方案
	MedicalAdvice  string    `gorm:"type:text" json:"medical_advice"`  // 医嘱
	Addon          string    `gorm:"type:text" json:"addon"`           // 附加信息
	Money          string    `gorm:"size:50" json:"money"`             // 费用
	Doc            string    `gorm:"size:20" json:"doc"`               // 医生
	CreatedAt      time.Time `json:"created_at"`                       // 创建时间
	UpdatedAt      time.Time `json:"updated_at"`                       // 更新时间
}

// TableName 指定表名为小写下划线格式
func (Patient) TableName() string {
	return "patient"
}

// Config 全局配置模型
type Config struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Key       string    `gorm:"uniqueIndex;size:50" json:"key"`     // 配置键
	Value     string    `gorm:"size:200" json:"value"`              // 配置值
	CreatedAt time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                         // 更新时间
}

// TableName 指定表名为小写下划线格式
func (Config) TableName() string {
	return "config"
}

// Doctor 医师信息模型
type Doctor struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Name      string    `gorm:"size:50" json:"name"`                // 姓名
	Sex       string    `gorm:"size:10" json:"sex"`                 // 性别
	Level     string    `gorm:"size:50" json:"level"`               // 级别
	Phone     string    `gorm:"size:20" json:"phone"`               // 电话
	Email     string    `gorm:"size:100" json:"email"`              // 邮箱
	Remark    string    `gorm:"type:text" json:"remark"`            // 备注
	CreatedAt time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                         // 更新时间
}

// TableName 指定表名为小写下划线格式
func (Doctor) TableName() string {
	return "doctor"
}

const (
	TemplateTypeDiagnosis = "diagnosis"
	TemplateTypeAdvice    = "advice"
	TemplateTypeTreatment = "treatment"
)

// 文本模板 三类：诊断，医嘱，治疗方案
type Template struct {
	ID      uint   `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Name    string `gorm:"size:50,unique" json:"name"`         // 名称
	Type    string `gorm:"size:50" json:"type"`                // 类型
	Content string `gorm:"type:text" json:"content"`           // 内容
}

// TableName 指定表名为小写下划线格式
func (Template) TableName() string {
	return "template"
}

// Medicine 药品信息模型
type Medicine struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Name         string    `gorm:"size:100" json:"name"`               // 药品名
	Manufacturer string    `gorm:"size:200" json:"manufacturer"`       // 生产厂家
	Price        string    `gorm:"size:50" json:"price"`               // 价格
	Specification string   `gorm:"size:200" json:"specification"`      // 药品规格
	CreatedAt    time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt    time.Time `json:"updated_at"`                         // 更新时间
}

// TableName 指定表名为小写下划线格式
func (Medicine) TableName() string {
	return "medicine"
}
