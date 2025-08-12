package main

// 迁移兼容旧版数据库
// 旧版数据库表
type OldPigTable struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"` // 主键ID
	Date          string `gorm:"size:50" json:"date"`                // 日期
	Time          string `gorm:"size:50" json:"time"`                // 准确时间
	Name          string `gorm:"size:20" json:"name"`                // 姓名
	Sex           string `gorm:"size:10" json:"sex"`                 // 性别
	Age           string `gorm:"size:10" json:"age"`                 // 年龄
	IllTime       string `gorm:"size:100" json:"ill_time"`           // 患病时间
	Phone         string `gorm:"size:20" json:"phone"`               // 手机号
	Work          string `gorm:"size:200" json:"work"`               // 工作
	Parent        string `gorm:"size:200" json:"parent"`             // 家长
	Address       string `gorm:"type:text" json:"address"`           // 地址
	Detail        string `gorm:"type:text" json:"detail"`            // 详情
	Solution      string `gorm:"type:text" json:"solution"`          // 解决方案
	MedicalAdvice string `gorm:"type:text" json:"medical_advice"`    // 医嘱
	Addon         string `gorm:"type:text" json:"addon"`             // 附加信息
	Money         string `gorm:"size:50" json:"money"`               // 费用
	Doc           string `gorm:"size:20" json:"doc"`                 // 医生
}
