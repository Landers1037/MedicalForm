package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/gorm"
)

// MedicineCrawler 药品爬虫结构体
type MedicineCrawler struct {
	db     *gorm.DB     // 数据库连接
	client *http.Client // HTTP客户端
}

// NewMedicineCrawler 创建新的药品爬虫实例
func NewMedicineCrawler(db *gorm.DB) *MedicineCrawler {
	return &MedicineCrawler{
		db: db,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CrawlMedicineData 爬取药品数据
func (mc *MedicineCrawler) CrawlMedicineData() error {
	log.Println("开始爬取国产药品数据...")

	// 基础URL
	baseURL := "http://www.54md.com/guochanyaopin/"

	// 获取药品列表页面
	doc, err := mc.fetchPage(baseURL)
	if err != nil {
		return fmt.Errorf("获取药品列表页面失败: %v", err)
	}

	// 解析药品信息
	medicines, err := mc.parseMedicineList(doc)
	if err != nil {
		return fmt.Errorf("解析药品信息失败: %v", err)
	}

	// 保存到数据库
	for _, medicine := range medicines {
		if err := mc.saveMedicine(medicine); err != nil {
			log.Printf("保存药品信息失败: %s, 错误: %v", medicine.Name, err)
			continue
		}
		log.Printf("成功保存药品: %s", medicine.Name)

		// 添加延时避免请求过快
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("药品数据爬取完成，共处理 %d 条记录", len(medicines))
	return nil
}

// fetchPage 获取网页内容
func (mc *MedicineCrawler) fetchPage(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头模拟浏览器
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := mc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP状态码错误: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// parseMedicineList 解析药品列表
func (mc *MedicineCrawler) parseMedicineList(doc *goquery.Document) ([]Medicine, error) {
	var medicines []Medicine
	// 查找药品信息的容器
	// 根据网站结构调整选择器
	doc.Find(".datatable tbody tr td a").Each(func(i int, s *goquery.Selection) {
		medicine := Medicine{}
		// 尝试提取药品名称
		name := s.First().Text()
		if name == "" {
			// 如果没有找到，尝试从链接文本中提取
			name = strings.TrimSpace(s.Find("a").First().Text())
		}

		// 清理药品名称
		name = mc.cleanText(name)
		if name == "" || len(name) < 2 {
			return // 跳过无效的药品名称
		}

		medicine.Name = name
		// TODO

		// 设置时间戳
		now := time.Now()
		medicine.CreatedAt = now
		medicine.UpdatedAt = now

		// 只有药品名称不为空才添加到列表
		if medicine.Name != "" {
			medicines = append(medicines, medicine)
		}
	})

	return medicines, nil
}

// extractText 从选择器中提取文本
func (mc *MedicineCrawler) extractText(s *goquery.Selection, selectors string) string {
	for _, selector := range strings.Split(selectors, ", ") {
		text := strings.TrimSpace(s.Find(selector).First().Text())
		if text != "" {
			return text
		}
	}
	return ""
}

// extractManufacturerFromText 从文本中提取厂家信息
func (mc *MedicineCrawler) extractManufacturerFromText(text string) string {
	// 常见的厂家关键词模式
	patterns := []string{
		`([^\s]+(?:制药|药业|医药|生物|科技|有限公司|股份|集团)[^\s]*)`,
		`生产企业[：:](.*?)(?:\s|$)`,
		`厂家[：:](.*?)(?:\s|$)`,
		`生产厂家[：:](.*?)(?:\s|$)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// parseMedicineFromText 从文本中解析药品信息
func (mc *MedicineCrawler) parseMedicineFromText(text string) Medicine {
	medicine := Medicine{}

	// 尝试分割文本获取药品信息
	parts := strings.Split(text, " ")
	if len(parts) > 0 {
		medicine.Name = mc.cleanText(parts[0])
	}

	// 查找厂家信息
	manufacturer := mc.extractManufacturerFromText(text)
	medicine.Manufacturer = manufacturer

	// 查找价格信息
	pricePattern := regexp.MustCompile(`(\d+(?:\.\d+)?)[元￥]`)
	priceMatches := pricePattern.FindStringSubmatch(text)
	if len(priceMatches) > 1 {
		medicine.Price = priceMatches[1] + "元"
	}

	// 查找规格信息
	specPattern := regexp.MustCompile(`(\d+(?:\.\d+)?(?:mg|g|ml|片|粒|支|盒)[^\s]*)`)
	specMatches := specPattern.FindStringSubmatch(text)
	if len(specMatches) > 1 {
		medicine.Specification = specMatches[1]
	}

	return medicine
}

// cleanText 清理文本内容
func (mc *MedicineCrawler) cleanText(text string) string {
	// 移除多余的空白字符
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")

	// 移除特殊字符
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.ReplaceAll(text, "\r", " ")

	return strings.TrimSpace(text)
}

// saveMedicine 保存药品信息到数据库
func (mc *MedicineCrawler) saveMedicine(medicine Medicine) error {
	// 检查是否已存在相同的药品
	var existingMedicine Medicine
	result := mc.db.Where("name = ? AND manufacturer = ?", medicine.Name, medicine.Manufacturer).First(&existingMedicine)

	if result.Error == nil {
		// 如果已存在，更新信息
		medicine.ID = existingMedicine.ID
		medicine.CreatedAt = existingMedicine.CreatedAt
		medicine.UpdatedAt = time.Now()
		return mc.db.Save(&medicine).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		// 如果不存在，创建新记录
		return mc.db.Create(&medicine).Error
	}

	return result.Error
}

// CrawlMedicineDataWithPagination 分页爬取药品数据
func (mc *MedicineCrawler) CrawlMedicineDataWithPagination(maxPages int) error {
	log.Println("开始分页爬取国产药品数据...")

	baseURL := "http://www.54md.com/guochanyaopin/"
	totalMedicines := 0

	for page := 1; page <= maxPages; page++ {
		var url string
		if page == 1 {
			url = baseURL
		} else {
			url = fmt.Sprintf("%s?page=%d", baseURL, page)
		}

		log.Printf("正在爬取第 %d 页: %s", page, url)

		doc, err := mc.fetchPage(url)
		if err != nil {
			log.Printf("获取第 %d 页失败: %v", page, err)
			continue
		}

		medicines, err := mc.parseMedicineList(doc)
		if err != nil {
			log.Printf("解析第 %d 页失败: %v", page, err)
			continue
		}

		if len(medicines) == 0 {
			log.Printf("第 %d 页没有找到药品信息，可能已到最后一页", page)
			break
		}

		// 保存药品信息
		for _, medicine := range medicines {
			// 检查是否存在重名药品
			var existingMedicine Medicine
			result := mc.db.Where("name = ?", medicine.Name).First(&existingMedicine)
			if result.Error == nil {
				// 如果存在重名药品，跳过插入
				log.Printf("药品 %s 已存在，跳过插入", medicine.Name)
				continue
			}
			
			if err := mc.saveMedicine(medicine); err != nil {
				log.Printf("保存药品信息失败: %s, 错误: %v", medicine.Name, err)
				continue
			}
			totalMedicines++
		}

		log.Printf("第 %d 页处理完成，本页获取 %d 条药品信息", page, len(medicines))

		// 添加延时避免请求过快
		time.Sleep(1 * time.Second)
	}

	log.Printf("分页爬取完成，共获取 %d 条药品信息", totalMedicines)
	return nil
}

// GetMedicineCount 获取数据库中药品数量
func (mc *MedicineCrawler) GetMedicineCount() (int64, error) {
	var count int64
	err := mc.db.Model(&Medicine{}).Count(&count).Error
	return count, err
}
