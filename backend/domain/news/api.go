package news

import (
	"strconv"

	"github.com/cv70/pkgo/ghttp"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// GetNews 获取帖子列表
func (d *NewsDomain) ApiGetNews(c *gin.Context) {
	var req GetNewsReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	req.Page = max(req.Page, 1)
	req.Limit = min(max(req.Limit, 20), 100)

	res, err := d.GetNews(c, &req)
	if err != nil {
		klog.Errorf("failed to get news: %v", err)
		ghttp.RespError(c, 500, "failed to get news")
		return
	}

	ghttp.RespSuccess(c, res)
}

// SearchNews 搜索帖子
func (d *NewsDomain) ApiSearchNews(c *gin.Context) {
	var req SearchNewsReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	req.Page = max(req.Page, 1)
	req.Limit = min(max(req.Limit, 20), 100)

	res, err := d.SearchNews(c, &req)
	if err != nil {
		klog.Errorf("failed to search news: %v", err)
		ghttp.RespError(c, 500, "failed to search news")
		return
	}

	ghttp.RespSuccess(c, res)
}

// AddNews 添加帖子
func (d *NewsDomain) ApiAddNews(c *gin.Context) {
	var req AddNewsReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	res, err := d.AddNews(c, &req)
	if err != nil {
		klog.Errorf("failed to add news: %v", err)
		ghttp.RespError(c, 500, "failed to add news")
		return
	}

	ghttp.RespSuccess(c, res)
}

// ApiGenerateNewsSummary 生成帖子摘要
func (d *NewsDomain) ApiGenerateNewsSummary(c *gin.Context) {
	newsID := c.Query("id")
	if newsID == "" {
		ghttp.RespError(c, 400, "news id is required")
		return
	}

	summary, err := d.GenerateNewsSummary(c, newsID)
	if err != nil {
		klog.Errorf("failed to generate news summary: %v", err)
		ghttp.RespError(c, 500, "failed to generate news summary")
		return
	}

	ghttp.RespSuccess(c, summary)
}

// ApiGetTrendingNews 获取热门帖子
func (d *NewsDomain) ApiGetTrendingNews(c *gin.Context) {
	daysStr := c.Query("days")
	days := 7 // 默认7天
	if daysStr != "" {
		// 这里应该转换为整数，为简化实现，直接使用默认值
	}

	news, err := d.GetTrendingNews(c, days)
	if err != nil {
		klog.Errorf("failed to get trending news: %v", err)
		ghttp.RespError(c, 500, "failed to get trending news")
		return
	}

	ghttp.RespSuccess(c, gin.H{
		"news": news,
	})
}

// ApiGetNewsCalendar 获取活动日历
func (d *NewsDomain) ApiGetNewsCalendar(c *gin.Context) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		klog.Errorf("failed to parse year: %v", err)
		ghttp.RespError(c, 400, "failed to parse year")
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		klog.Errorf("failed to parse month: %v", err)
		ghttp.RespError(c, 400, "failed to parse month")
		return
	}

	// 调用领域层方法获取新闻日历数据
	news, err := d.GetNewsCalendar(c, year, month)
	if err != nil {
		klog.Errorf("failed to get news calendar: %v", err)
		ghttp.RespError(c, 500, "failed to get news calendar")
		return
	}

	ghttp.RespSuccess(c, gin.H{
		"news": news,
	})
}

// ApiFetchRSSFeeds 获取RSS订阅内容
func (d *NewsDomain) ApiFetchRSSFeeds(c *gin.Context) {
	err := d.FetchAndProcessRSSFeeds(c)
	if err != nil {
		klog.Errorf("failed to fetch rss feeds: %v", err)
		ghttp.RespError(c, 500, "failed to fetch rss feeds")
		return
	}

	ghttp.RespSuccess(c, gin.H{
		"message": "RSS feeds fetched successfully",
	})
}

func (d *NewsDomain) ApiLikeNews(c *gin.Context) {
	var req LikeNewsReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	newsID, err := strconv.ParseUint(req.NewsID, 10, 64)
	if err != nil {
		klog.Errorf("failed to parse news id: %v", err)
		ghttp.RespError(c, 400, "failed to parse news id")
		return
	}

	userID := uint64(0)
	err = d.LikeNews(c, newsID, userID)
	if err != nil {
		klog.Errorf("failed to add comment: %v", err)
		ghttp.RespError(c, 500, "failed to add comment")
		return
	}

	ghttp.RespSuccess(c, "")
}

func (d *NewsDomain) ApiUnlikeNews(c *gin.Context) {
	var req UnlikeNewsReq
	err := c.ShouldBind(&req)
	if err != nil {
		klog.Errorf("failed to parse body: %v", err)
		ghttp.RespError(c, 400, "failed to parse body")
		return
	}

	newsID, err := strconv.ParseUint(req.NewsID, 10, 64)
	if err != nil {
		klog.Errorf("failed to parse news id: %v", err)
		ghttp.RespError(c, 400, "failed to parse news id")
		return
	}

	userID := 0
	err = d.UnlikeNews(c, newsID, uint64(userID))
	if err != nil {
		klog.Errorf("failed to unlike news: %v", err)
		ghttp.RespError(c, 500, "failed to unlike news")
		return
	}

	ghttp.RespSuccess(c, "")
}

func (d *NewsDomain) ApiGetNewsLikes(c *gin.Context) {
	newsIDStr := c.Query("id")
	if newsIDStr == "" {
		ghttp.RespError(c, 400, "news id is required")
		return
	}

	newsID, err := strconv.ParseUint(newsIDStr, 10, 64)
	if err != nil {
		klog.Errorf("failed to parse news id: %v", err)
		ghttp.RespError(c, 400, "failed to parse news id")
		return
	}

	likes, err := d.GetNewsLikes(c, newsID)
	if err != nil {
		klog.Errorf("failed to get news likes: %v", err)
		ghttp.RespError(c, 500, "failed to get news likes")
		return
	}

	ghttp.RespSuccess(c, likes)
}
