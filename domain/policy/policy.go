package policy

import (
	"context"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"vc-go/datasource"

	"github.com/pkg/errors"
)

// MatchPoliciesByEnterprise 匹配适合企业的政策
func (d *PolicyDomain) MatchPoliciesByEnterprise(ctx context.Context, req *PolicyMatchReq) (*PolicyMatchResp, error) {
	// 获取所有政策数据
	policies, _, err := d.DB.GetPolicies("", "", 1, 1000, true)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get policies")
	}

	// 计算匹配度
	matchedPolicies := make([]MatchedPolicy, 0)
	for _, policy := range policies {
		score, reason := calculateMatchScore(policy, req)
		if score > 0.3 { // 设置阈值，只有匹配度超过30%的才返回
			matchedPolicies = append(matchedPolicies, MatchedPolicy{
				ID:          strconv.Itoa(int(policy.ID)),
				Title:       policy.Title,
				Content:     policy.Content,
				Region:      policy.Region,
				Industry:    policy.Industry,
				Score:       score,
				MatchReason: reason,
			})
		}
	}

	// 按匹配度排序
	sort.Slice(matchedPolicies, func(i, j int) bool {
		return matchedPolicies[i].Score > matchedPolicies[j].Score
	})

	return &PolicyMatchResp{
		MatchedPolicies: matchedPolicies,
	}, nil
}

// calculateMatchScore 计算政策与企业的匹配度
func calculateMatchScore(policy datasource.Policy, req *PolicyMatchReq) (float64, string) {
	score := 0.0
	reasons := make([]string, 0)

	// 地区匹配
	if strings.Contains(policy.Region, req.Region) || req.Region == "" {
		score += 0.3
		reasons = append(reasons, "地区匹配")
	}

	// 行业匹配
	if strings.Contains(policy.Industry, req.Industry) || req.Industry == "" {
		score += 0.3
		reasons = append(reasons, "行业匹配")
	}

	// 从政策内容中提取企业类型要求并匹配
	enterpriseTypeMatched := checkEnterpriseType(policy.Content, req.EnterpriseType)
	if enterpriseTypeMatched {
		score += 0.2
		reasons = append(reasons, "企业类型匹配")
	}

	// 从政策内容中提取企业规模要求并匹配
	sizeMatched := checkEnterpriseSize(policy.Content, req.EnterpriseSize)
	if sizeMatched {
		score += 0.2
		reasons = append(reasons, "企业规模匹配")
	}

	// 限制最大分数为1.0
	if score > 1.0 {
		score = 1.0
	}

	return score, strings.Join(reasons, ", ")
}

// checkEnterpriseType 检查企业类型是否匹配
func checkEnterpriseType(policyContent, enterpriseType string) bool {
	// 定义不同类型企业的关键词
	keywords := map[string][]string{
		"初创企业":   {"初创", "创业", "初创期", "创业初期", "小微企业"},
		"中小企业":   {"中小企业", "中型企业", "小型企业", "微型", "小企业"},
		"高新技术企业": {"高新技术", "高新技术企业", "科技型企业", "创新型企业"},
	}

	if enterpriseType == "" {
		return true
	}

	if keywords, exists := keywords[enterpriseType]; exists {
		for _, keyword := range keywords {
			if strings.Contains(policyContent, keyword) {
				return true
			}
		}
	}

	return false
}

// checkEnterpriseSize 检查企业规模是否匹配
func checkEnterpriseSize(policyContent, enterpriseSize string) bool {
	// 解析企业规模字符串，例如 "1-50人", "51-200人" 等
	if enterpriseSize == "" {
		return true
	}

	// 提取数字范围
	re := regexp.MustCompile(`(\d+)-(\d+)`)
	matches := re.FindStringSubmatch(enterpriseSize)
	if len(matches) == 3 {
		min, err1 := strconv.Atoi(matches[1])
		max, err2 := strconv.Atoi(matches[2])
		if err1 == nil && err2 == nil {
			// 检查政策内容是否包含相关规模要求
			return checkSizeInPolicy(policyContent, min, max)
		}
	}

	// 如果是特殊描述如"大型企业"、"小型企业"等
	return checkSizeDescription(policyContent, enterpriseSize)
}

// checkSizeInPolicy 检查政策中是否包含特定规模要求
func checkSizeInPolicy(policyContent string, min, max int) bool {
	// 检查政策内容中是否提到员工数量、年收入等规模相关要求
	// 这里使用简单的关键词匹配，实际项目中可以使用更复杂的NLP算法
	if min <= 50 && max >= 1 {
		// 小微企业相关关键词
		microKeywords := []string{"小微企业", "小型", "微型", "初创", "1-50人"}
		for _, keyword := range microKeywords {
			if strings.Contains(policyContent, keyword) {
				return true
			}
		}
	}

	if min <= 200 && max >= 51 {
		// 中型企业相关关键词
		smallMediumKeywords := []string{"中型", "51-200人", "中小", "中等规模"}
		for _, keyword := range smallMediumKeywords {
			if strings.Contains(policyContent, keyword) {
				return true
			}
		}
	}

	return false
}

// checkSizeDescription 检查规模描述
func checkSizeDescription(policyContent, sizeDesc string) bool {
	// 检查是否包含特定规模描述
	scaleKeywords := map[string][]string{
		"大型企业": {"大型企业", "大型", "大企业"},
		"中型企业": {"中型企业", "中型", "中等企业", "中等规模"},
		"小型企业": {"小型企业", "小型", "小企业"},
		"微型企业": {"微型企业", "微型", "微小", "小微企业"},
	}

	if keywords, exists := scaleKeywords[sizeDesc]; exists {
		for _, keyword := range keywords {
			if strings.Contains(policyContent, keyword) {
				return true
			}
		}
	}

	return false
}
