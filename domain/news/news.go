package news

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"
	"vc-go/datasource/scylladao"
	"vc-go/pkg/gslice"

	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"k8s.io/klog/v2"
)

// AddNews 添加帖子
func (d *NewsDomain) AddNews(ctx *gin.Context, req *AddNewsReq) (*AddNewsRes, error) {
	news := &scylladao.NewsRow{
		Title:       req.Title,
		Content:     req.Content,
		Category:    req.Category,
		Tags:        req.Tags,
		Region:      req.Region,
		Industry:    req.Industry,
	}

	err := d.Scylla.SaveNews(news)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to insert news")
	}

	// 同步到milvus
	return &AddNewsRes{
		ID: news.ID.String(),
	}, nil
}

// GetNews 获取帖子列表
func (d *NewsDomain) GetNews(ctx *gin.Context, req *GetNewsReq) (*GetNewsResp, error) {
	news, total, err := d.DB.GetNews(req.Category, req.Region, req.Industry, req.Tags, req.Page, req.Limit)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get news")
	}

	newsResp := gslice.Map(news, func(newsItem *datasource.News) *NewsItemResp {
		newsItemResp := &NewsItemResp{}
		newsItemResp.From(newsItem)
		return newsItemResp
	})

	return &GetNewsResp{
		NewsList: newsResp,
		Total:    total,
	}, nil
}

// SearchNews 搜索帖子
func (d *NewsDomain) SearchNews(ctx *gin.Context, req *SearchNewsReq) (*SearchNewsRes, error) {
	news, total, err := d.DB.SearchNews(req.Query, req.Page, req.Limit)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get news")
	}

	return &SearchNewsRes{
		News:  news,
		Total: total,
	}, nil
}

// syncToES 同步到Elasticsearch
func (d *NewsDomain) syncToES(news News) error {
	// 这里应该实现同步到ES的逻辑
	// 为简化实现，这里只做模拟
	klog.Infof("Syncing news to ES: %s", news.Title)
	return nil
}

// FetchAndProcessRSSFeeds 获取并处理RSS订阅
func (d *NewsDomain) FetchAndProcessRSSFeeds(ctx context.Context) error {
	// 获取所有活跃的RSS源
	feeds, err := d.getRSSFeeds(true)
	if err != nil {
		return errors.WithMessage(err, "failed to get rss feeds")
	}

	// 遍历每个RSS源并获取最新内容
	for _, feed := range feeds {
		err := d.fetchAndProcessFeed(feed)
		if err != nil {
			klog.Errorf("failed to fetch feed %s: %v", feed.URL, err)
			continue
		}
	}

	return nil
}

// getRSSFeeds 获取RSS订阅源列表
func (d *NewsDomain) getRSSFeeds(activeOnly bool) ([]*RSSFeed, error) {
	// 这里在实际项目中应该从数据库获取RSS订阅源
	// 为简化实现，返回一些示例数据
	feeds := []*RSSFeed{
		{
			ID:          "1",
			URL:         "https://example.com/rss",
			Name:        "政策解读RSS",
			Active:      true,
			LastFetched: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			ID:          "2",
			URL:         "https://example.com/funding/rss",
			Name:        "融资动态RSS",
			Active:      true,
			LastFetched: time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	if activeOnly {
		activeFeeds := make([]*RSSFeed, 0)
		for _, feed := range feeds {
			if feed.Active {
				activeFeeds = append(activeFeeds, feed)
			}
		}
		return activeFeeds, nil
	}

	return feeds, nil
}

// fetchAndProcessFeed 获取并处理单个RSS源
func (d *NewsDomain) fetchAndProcessFeed(feed *RSSFeed) error {
	fp := gofeed.NewParser()
	feedData, err := fp.ParseURL(feed.URL)
	if err != nil {
		return errors.WithMessage(err, "failed to parse rss feed")
	}

	// 处理每个RSS项目
	for _, item := range feedData.Items {
		// 检查是否已存在该帖子
		exists, err := d.checkNewsExists(item.Link)
		if err != nil || exists {
			continue // 如果已存在或检查出错，跳过
		}

		// 创建新的帖子条目
		newsItem := &AddNewsReq{
			Title:       item.Title,
			Content:     item.Description,
			Category:    d.categorizeNews(item.Title, item.Description),
			Source:      feed.Name,
			Tags:        d.extractTags(item.Title, item.Description),
			Region:      "", // 从内容中提取地区信息
			Industry:    "", // 从内容中提取行业信息
			PublishDate: item.Published,
		}

		// 添加到数据库
		_, err = d.AddNews(nil, newsItem)
		if err != nil {
			klog.Errorf("failed to add rss item: %v", err)
			continue
		}
	}

	// 更新最后获取时间
	d.updateLastFetched(feed.ID)

	return nil
}

// checkNewsExists 检查帖子是否已存在（通过链接）
func (d *NewsDomain) checkNewsExists(link string) (bool, error) {
	// 这里在实际项目中应该查询数据库检查是否已存在
	// 为简化实现，总是返回false
	return false, nil
}

// updateLastFetched 更新最后获取时间
func (d *NewsDomain) updateLastFetched(feedID string) {
	// 这里在实际项目中应该更新数据库中的最后获取时间
	// 为简化实现，只打印日志
	klog.Infof("Updated last fetched time for feed %s", feedID)
}

// categorizeNews 根据标题和内容分类新闻
func (d *NewsDomain) categorizeNews(title, content string) string {
	title = strings.ToLower(title)
	content = strings.ToLower(content)

	// 定义分类关键词
	categoryKeywords := map[string][]string{
		"政策解读": {"政策", "政府", "补贴", "优惠", "法规", "条例", "通知", "发布", "改革"},
		"成功案例": {"成功", "案例", "创业", "经验", "分享", "故事", "企业", "成就", "突破"},
		"融资动态": {"融资", "投资", "资金", "投资方", "融资额", "B轮", "A轮", "天使轮", "VC", "PE"},
		"活动日历": {"活动", "会议", "峰会", "论坛", "报名", "日程", "议程", "参加"},
	}

	for category, keywords := range categoryKeywords {
		for _, keyword := range keywords {
			if strings.Contains(title, keyword) || strings.Contains(content, keyword) {
				return category
			}
		}
	}

	// 默认分类
	return "创业帖子"
}

// extractTags 从标题和内容中提取标签
func (d *NewsDomain) extractTags(title, content string) []string {
	tags := make([]string, 0)

	// 定义标签关键词
	tagKeywords := []string{
		"人工智能", "大数据", "云计算", "区块链", "物联网", "新能源", "生物医药", "新材料",
		"互联网+", "智能制造", "数字经济", "乡村振兴", "绿色发展", "科技创新", "创业",
		"投资", "融资", "政策", "补贴", "优惠", "创业大赛", "孵化器", "加速器",
	}

	allText := strings.ToLower(title + " " + content)

	for _, keyword := range tagKeywords {
		if strings.Contains(allText, strings.ToLower(keyword)) {
			tags = append(tags, keyword)
		}
	}

	// 使用正则表达式提取其他可能的标签（如地区、行业等）
	// 这里简化实现，只提取一些常见的地名和行业
	regionPattern := regexp.MustCompile(`(北京|上海|深圳|广州|杭州|成都|武汉|西安|南京|苏州|重庆|天津|青岛|大连|宁波|厦门|深圳|广州)`)
	regionMatches := regionPattern.FindAllString(allText, -1)
	for _, region := range regionMatches {
		if !contains(tags, region) {
			tags = append(tags, region)
		}
	}

	return tags
}

// contains 检查字符串切片是否包含指定字符串
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// GenerateNewsSummary 生成帖子摘要
func (d *NewsDomain) GenerateNewsSummary(ctx context.Context, newsID string) (*NewsSummary, error) {
	// 获取帖子详情
	newsIDInt, err := strconv.Atoi(newsID)
	if err != nil {
		return nil, errors.New("invalid news id")
	}

	news, err := d.DB.DetailNews(uint(newsIDInt))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get news detail")
	}

	// 使用AI生成摘要（简化实现，返回内容的前100个字符）
	summary := news.Content
	if len(summary) > 100 {
		summary = summary[:100] + "..."
	}

	return &NewsSummary{
		ID:      newsID,
		Title:   news.Title,
		Summary: summary,
	}, nil
}

// GetTrendingNews 获取热门帖子
func (d *NewsDomain) GetTrendingNews(ctx context.Context, days int) ([]*datasource.News, error) {
	if days <= 0 {
		days = 7 // 默认获取7天内的热门帖子
	}

	// 计算开始日期
	_ = time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	// 这里在实际项目中应该根据访问量、分享量等指标来判断热门帖子
	// 为简化实现，返回最近的帖子
	newsList, _, err := d.DB.GetNews("", "", "", nil, 1, 20)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get trending news")
	}

	return newsList, nil
}

// GetNewsCalendar 获取活动日历
func (d *NewsDomain) GetNewsCalendar(ctx context.Context, year, month int) ([]*datasource.News, error) {
	// 这里在实际项目中应该查询指定月份的活动日历
	// 为简化实现，返回活动日历类别的帖子
	newsList, _, err := d.DB.GetNews("活动日历", "", "", nil, 1, 100)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get calendar news")
	}

	return newsList, nil
}

// LikeNews 点赞帖子
func (d *NewsDomain) LikeNews(ctx context.Context, newsID, userID int64) error {
	// 插入用户点赞记录
	err := d.Scylla.SaveUserLikeNews(userID, newsID)
	if err != nil {
		return errors.WithMessage(err, "failed to insert news like")
	}

	// 更新点赞计数
	err = d.Scylla.IncrNewsLikeCount(newsID)
	if err != nil {
		return errors.WithMessage(err, "failed to create scylla session")
	}
	return nil
}

// UnLikeNews 取消点赞帖子
func (d *NewsDomain) UnlikeNews(ctx context.Context, newsID, userID int64) error {
	// 删除用户点赞记录
	err := d.Scylla.DeleteUserLikeNews(userID, newsID)
	if err != nil {
		return errors.WithMessage(err, "failed to delete news like")
	}

	// 更新点赞计数
	err = d.Scylla.DecrNewsLikeCount(newsID)
	if err != nil {
		return errors.WithMessage(err, "failed to update news like count")
	}
	return nil
}

// GetNewsLikeCount 获取帖子点赞数
func (d *NewsDomain) GetNewsLikeCount(ctx context.Context, newsID int64) (int64, error) {
	// 查询点赞计数
	count, err := d.Scylla.GetNewsLikeCount(newsID)
	if err != nil {
		return 0, errors.WithMessage(err, "failed to get news like count")
	}
	return count, err
}

// FavoriteNews 收藏帖子
func (d *NewsDomain) FavoriteNews(ctx context.Context, userID, newsID int64) error {
	// 插入用户收藏记录
	err := d.Scylla.SaveUserFavoriteNews(userID, newsID)
	if err != nil {
		return errors.WithMessage(err, "failed to insert news favorite")
	}

	// 更新收藏计数
	err = d.Scylla.IncrUserFavoriteNewsCount(newsID)
	if err != nil {
		return errors.WithMessage(err, "failed to update news favorite count")
	}
	return nil
}

// UnfavoriteNews 取消收藏帖子
func (d *NewsDomain) UnfavoriteNews(ctx context.Context, newsID, userID int64) error {
	// 删除用户收藏记录
	err := d.Scylla.DeleteUserFavoriteNews(userID, newsID)
	if err != nil {
		return errors.WithMessage(err, "failed to delete news favorite")
	}

	// 更新收藏计数
	err = d.Scylla.DecrUserFavoriteNewsCount(newsID)
	if err != nil {
		return errors.WithMessage(err, "failed to update news favorite count")
	}
	return nil
}

// GetFavoriteNews 查询用户收藏的帖子
func (d *NewsDomain) GetFavoriteNews(ctx context.Context, userID int64) ([]*NewsItemResp, error) {
	// 查询用户收藏的帖子ID列表
	newsIDs, err := d.Scylla.GetUserFavoriteNewsIDs(userID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get user favorite news ids")
	}

	// 如果没有收藏记录，返回空列表
	if len(newsIDs) == 0 {
		return []*NewsItemResp{}, nil
	}

	// 根据帖子ID列表查询帖子详情
	newsList, err := d.Scylla.GetNewsList(newsIDs...)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get news list")
	}

	newsItems := gslice.Map(newsList, func(news *datasource.News) *NewsItemResp {
		newsItem := &NewsItemResp{}
		newsItem.From(news)
		return newsItem
	})

	return newsItems, err
}

// CommentNews 评论帖子
func (d *NewsDomain) CommentNews(ctx context.Context, newsID, userID, commentID uint64, comment string) error {

	// 插入评论记录
	if err := session.Query(`INSERT INTO news_comment (destination, news_id, user_id, comment, created_at) VALUES (?, ?, ?, ?, ?)`, 0, newsID, userID, comment, time.Now().Format("2006-01-02 15:04:05")).Exec(); err != nil {
		return errors.WithMessage(err, "failed to insert news comment")
	}
	return nil
}

// DeleteComment 删除评论
func (d *NewsDomain) DeleteComment(ctx context.Context, commentID uint64) error {
	// 创建 ScyllaDB 会话
	session, err := d.Scylla.CreateSession()
	if err != nil {
		return errors.WithMessage(err, "failed to create scylla session")
	}
	defer session.Close()

	// 删除评论记录
	if err := session.Query(`DELETE FROM news_comment WHERE id = ?`, commentID).Exec(); err != nil {
		return errors.WithMessage(err, "failed to delete news comment")
	}
	return nil
}

// LikeComment 点赞评论
func (d *NewsDomain) LikeComment(ctx context.Context, commentID, userID uint64) error {
	// 创建 ScyllaDB 会话
	session, err := d.Scylla.CreateSession()
	if err != nil {
		return errors.New("invalid news id")
	}
	defer session.Close()

	// 插入评论点赞记录（主键为 (comment_id, user_id)）
	if err := session.Query(`INSERT INTO news_comment_like (comment_id, user_id, created_at) VALUES (?, ?, ?)`, commentID, userID, time.Now().Format("2006-01-02 15:04:05")).Exec(); err != nil {
		return errors.WithMessage(err, "failed to insert news comment like")
	}

	// 更新评论点赞计数
	if err := session.Query(`UPDATE news_comment_like_count SET count = count + 1 WHERE comment_id = ?`, commentID).Exec(); err != nil {
		return errors.WithMessage(err, "failed to update news comment like count")
	}
	return nil
}

// UnlikeComment 取消点赞评论
func (d *NewsDomain) UnlikeComment(ctx context.Context, commentID, userID uint64) error {
	// 创建 ScyllaDB 会话
	session, err := d.Scylla.CreateSession()
	if err != nil {
		return errors.New("invalid news id")
	}
	defer session.Close()
	// 删除点赞记录
	if err := session.Query(`DELETE FROM news_comment_like WHERE comment_id = ? AND user_id = ?`, commentID, userID).Exec(); err != nil {
		return errors.WithMessage(err, "failed to delete news comment like")
	}

	// 更新评论点赞计数（减一）
	if err := session.Query(`UPDATE news_comment_like_count SET count = count - 1 WHERE comment_id = ?`, commentID).Exec(); err != nil {
		return errors.WithMessage(err, "failed to update news comment like count")
	}
	return nil
}

// 点赞 收藏 评论 分享 订阅
