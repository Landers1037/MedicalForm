package main

import (
	"testing"
)

// TestMedicineCrawler 测试药品爬虫功能
func TestMedicineCrawler(t *testing.T) {
	t.Log("开始测试药品爬虫功能...")

	// 创建数据库服务
	dbService := NewDatabaseService()
	if err := dbService.InitDB(); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	// 创建药品爬虫
	crawler := NewMedicineCrawler(dbService.GetDB())

	// 获取爬取前的药品数量
	beforeCount, err := crawler.GetMedicineCount()
	if err != nil {
		t.Logf("获取药品数量失败: %v", err)
	} else {
		t.Logf("爬取前数据库中有 %d 条药品记录", beforeCount)
	}

	// 开始爬取药品数据（只爬取第一页进行测试）
	t.Log("开始爬取药品数据...")
	err = crawler.CrawlMedicineDataWithPagination(1)
	if err != nil {
		t.Logf("爬取药品数据失败: %v", err)
	} else {
		t.Log("药品数据爬取完成")
	}

	// 获取爬取后的药品数量
	afterCount, err := crawler.GetMedicineCount()
	if err != nil {
		t.Logf("获取药品数量失败: %v", err)
	} else {
		t.Logf("爬取后数据库中有 %d 条药品记录", afterCount)
		t.Logf("本次爬取新增 %d 条药品记录", afterCount-beforeCount)
	}

	// 显示前几条药品信息
	var medicines []Medicine
	result := dbService.GetDB().Limit(5).Find(&medicines)
	if result.Error != nil {
		t.Logf("查询药品信息失败: %v", result.Error)
	} else {
		t.Log("前5条药品信息:")
		for i, medicine := range medicines {
			t.Logf("%d. 名称: %s, 厂家: %s, 价格: %s, 规格: %s\n",
				i+1, medicine.Name, medicine.Manufacturer, medicine.Price, medicine.Specification)
		}
	}

	t.Log("药品爬虫测试完成")
}

// 如果直接运行此文件，则执行测试
// 注意：这个函数不会在正常的应用程序中被调用
// 只是用于独立测试药品爬虫功能
func init() {
	// 可以在这里添加测试逻辑，但通常不建议在init中执行耗时操作
}
