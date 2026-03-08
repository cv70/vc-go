package news

import (
	"vc-go/datasource"
)

type GetNewsReq struct {
	Category string   `json:"category"` // 分类筛选
	Region   string   `json:"region"`   // 地区筛选
	Industry string   `json:"industry"` // 行业筛选
	Tags     []string `json:"tags"`     // 标签筛选
	Page     int      `json:"page" binding:"min=1"`
	Limit    int      `json:"limit" binding:"min=1,max=20"`
}

type NewsItemResp struct {
	ID          uint64   `json:"id,string"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Category    string   `json:"category"` // 分类：政策解读、成功案例、融资动态、活动日历
	Source      string   `json:"source"`   // 来源
	PublishDate string   `json:"publish_date"`
	Tags        []string `json:"tags"`     // 标签
	Region      string   `json:"region"`   // 地区
	Industry    string   `json:"industry"` // 行业
}

func (r *NewsItemResp) From(news *datasource.News) {
	r.ID = news.ID
	r.Title = news.Title
	r.Content = news.Content
	r.Category = news.Category
	r.Source = news.Source
	r.PublishDate = news.PublishDate
	r.Tags = news.Tags
	r.Region = news.Region
	r.Industry = news.Industry
}

type GetNewsResp struct {
	NewsList []*NewsItemResp `json:"news"`
	Total    int64           `json:"total"`
}

type AddNewsReq struct {
	Title       string   `json:"title" binding:"required"`
	Content     string   `json:"content" binding:"required"`
	Category    string   `json:"category" binding:"required"`
	Source      string   `json:"source" binding:"required"`
	Tags        []string `json:"tags"`
	Region      string   `json:"region"`
	Industry    []string `json:"industry"`
	PublishDate string   `json:"publish_date"`
}

type AddNewsRes struct {
	ID string `json:"id"`
}

type SearchNewsReq struct {
	Query string `json:"query" binding:"required"`
	Page  int    `json:"page" binding:"min=1"`
	Limit int    `json:"limit" binding:"min=1,max=20"`
}

type SearchNewsRes struct {
	News  []*datasource.News `json:"news"`
	Total int                `json:"total"`
}

// RSSFeed RSS订阅源
type RSSFeed struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	LastFetched string `json:"last_fetched"`
}

// RSSItem RSS项目
type RSSItem struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	PubDate     string `json:"pub_date"`
	Content     string `json:"content"`
	Source      string `json:"source"`
}

// NewsSummary 帖子摘要
type NewsSummary struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

// AddRSSFeedReq 添加RSS订阅源请求
type AddRSSFeedReq struct {
	URL  string `json:"url" binding:"required"`
	Name string `json:"name" binding:"required"`
}

// AddRSSFeedRes 添加RSS订阅源响应
type AddRSSFeedRes struct {
	ID string `json:"id"`
}

// GetRSSFeedsReq 获取RSS订阅源请求
type GetRSSFeedsReq struct {
	Page  int `json:"page" binding:"min=1"`
	Limit int `json:"limit" binding:"min=1,max=20"`
}

// GetRSSFeedsRes 获取RSS订阅源响应
type GetRSSFeedsRes struct {
	Feeds []*RSSFeed `json:"feeds"`
	Total int        `json:"total"`
}

// LikeNewsReq 点赞帖子请求
type LikeNewsReq struct {
	NewsID string `json:"news_id" binding:"required"`
}

// UnlikeNewsReq 取消点赞帖子请求
type UnlikeNewsReq struct {
	NewsID string `json:"news_id" binding:"required"`
}
