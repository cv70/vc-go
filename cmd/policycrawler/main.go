package main

import (
	"context"
	"log"
	"strings"
	"time"
	"vc-go/config"
	"vc-go/dao"
	"vc-go/datasource"
	"vc-go/infra"
	"vc-go/pkg/mistake"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig()
	mistake.Unwrap(err)

	// Initialize database
	db, err := infra.NewDB(ctx, cfg.Database)
	mistake.Unwrap(err)

	// 启动爬虫
	crawler := NewPolicyCrawler(db)
	crawler.Start()
}

type PolicyCrawler struct {
	db *dao.DB
}

func NewPolicyCrawler(db *dao.DB) *PolicyCrawler {
	return &PolicyCrawler{
		db: db,
	}
}

func (c *PolicyCrawler) Start() {
	log.Println("Starting policy crawler...")

	// 定义要爬取的政府网站列表
	governmentSites := []string{
		"http://www.gov.cn",
		"https://sme.miit.gov.cn",  // 工信部中小企业政策
		"https://www.stats.gov.cn", // 国家统计局
		// 可以添加更多政府网站
	}

	for _, site := range governmentSites {
		log.Printf("Crawling site: %s", site)
		policies, err := c.crawlSite(site)
		if err != nil {
			log.Printf("Error crawling site %s: %v", site, err)
			continue
		}

		// 保存政策数据
		for _, policy := range policies {
			err := c.db.InsertPolicy(&policy)
			if err != nil {
				log.Printf("Error saving policy: %v", err)
			}
		}
	}

	log.Println("Policy crawling completed")
}

func (c *PolicyCrawler) crawlSite(siteURL string) ([]datasource.Policy, error) {
	// 启动浏览器
	launcher := launcher.New()
	defer launcher.Cleanup()
	browser := rod.New().ControlURL(launcher.MustLaunch()).MustConnect()
	defer browser.Close()

	page := browser.MustPage(siteURL)

	// 等待页面加载
	page.MustWaitLoad()

	// 查找政策链接（这需要根据具体网站的结构进行调整）
	links := page.MustElements("a[href*='policy'], a[href*='zhengce'], a[href*='通知'], a[href*='公告'], a[href*='政策']")

	var policies []datasource.Policy

	for _, link := range links {
		href := link.MustAttribute("href")
		if href == nil {
			continue
		}

		// 构建完整URL
		policyURL := c.buildFullURL(siteURL, *href)

		// 获取政策详情
		policy, err := c.extractPolicyDetails(browser, policyURL)
		if err != nil {
			log.Printf("Error extracting policy from %s: %v", policyURL, err)
			continue
		}

		if policy != nil && policy.Title != "" && policy.Content != "" {
			policies = append(policies, *policy)
		}
	}

	return policies, nil
}

func (c *PolicyCrawler) extractPolicyDetails(browser *rod.Browser, url string) (*datasource.Policy, error) {
	page := browser.MustPage(url)
	page.MustWaitLoad()

	// 提取标题
	var title string
	titleEl := page.Timeout(5 * time.Second).Element("h1, .title, .article-title, [class*='title'], h2, h3")
	if titleEl != nil {
		title = titleEl.MustText()
	} else {
		// 如果没有找到标题元素，使用页面标题
		title = page.MustInfo().Title
	}

	// 提取内容
	var content string
	contentEl := page.Timeout(5 * time.Second).Element(".content, .article-content, .detail-content, [class*='content'], .main-content, .article-body, .post-content")
	if contentEl != nil {
		content = contentEl.MustText()
	} else {
		// 如果没有找到特定的内容元素，尝试获取body文本
		body := page.MustElement("body")
		content = body.MustText()
		// 清理内容，移除多余的空白字符
		content = strings.TrimSpace(content)
	}

	// 提取发布日期（这需要根据具体网站的结构进行调整）
	var publishDate string
	dateEl := page.Timeout(5 * time.Second).Element("[class*='date'], [class*='time'], .publish-time, time, [class*='publish']")
	if dateEl != nil {
		publishDate = dateEl.MustText()
	} else {
		// 如果没有找到日期元素，使用当前时间
		publishDate = time.Now().Format("2006-01-02")
	}

	// 简单提取地区和行业信息（实际项目中应使用更复杂的NLP技术）
	region := c.extractRegion(content)
	industry := c.extractIndustry(content)
	enterpriseSize := c.extractEnterpriseSize(content)

	policy := &datasource.Policy{
		Title:       strings.TrimSpace(title),
		Content:     strings.TrimSpace(content),
		Region:      region,
		Industry:    industry,
		PublishDate: publishDate,
		Source:      url,
	}

	return policy, nil
}

func (c *PolicyCrawler) buildFullURL(baseURL, href string) string {
	// 简单的URL构建逻辑，实际项目中应使用更完善的URL处理
	if strings.HasPrefix(href, "http") {
		return href
	}

	if strings.HasPrefix(href, "/") {
		// 提取域名部分
		parts := strings.Split(baseURL, "/")
		domain := parts[0] + "//" + parts[2]
		return domain + href
	}

	// 相对路径
	return baseURL + "/" + href
}

func (c *PolicyCrawler) extractRegion(content string) string {
	// 简单的地区提取，实际项目中应使用NLP技术
	// 这里只是示例，实际应根据具体需求实现
	regionKeywords := []string{"北京市", "上海市", "广东省", "深圳市", "浙江省", "杭州市", "江苏省", "南京市", "重庆市", "天津市", "四川省", "成都市", "湖北省", "武汉市", "陕西省", "西安市", "辽宁省", "沈阳市", "山东省", "青岛市", "福建省", "厦门市"}

	for _, keyword := range regionKeywords {
		if strings.Contains(content, keyword) {
			return keyword
		}
	}

	return ""
}

func (c *PolicyCrawler) extractIndustry(content string) string {
	// 简单的行业提取，实际项目中应使用NLP技术
	industryKeywords := []string{"科技", "制造业", "服务业", "互联网", "人工智能", "生物医药", "新能源", "新材料", "金融", "教育", "医疗", "农业", "能源", "交通", "建筑", "文化", "旅游", "环保", "信息产业", "高端制造", "集成电路", "软件", "通信", "电子商务", "大数据", "云计算", "物联网", "区块链", "智能硬件"}

	for _, keyword := range industryKeywords {
		if strings.Contains(content, keyword) {
			return keyword
		}
	}

	return ""
}

func (c *PolicyCrawler) extractEnterpriseSize(content string) string {
	// 简单的企业规模提取
	sizeKeywords := []string{"小微企业", "中小企业", "大型企业", "初创企业", "高新技术企业", "创新型企业", "民营企业", "国有企业", "合资企业", "外资企业"}

	for _, keyword := range sizeKeywords {
		if strings.Contains(content, keyword) {
			return keyword
		}
	}

	return ""
}
