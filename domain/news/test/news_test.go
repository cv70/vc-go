package news_test

import (
	"context"
	"testing"

	"vc-go/domain/news"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/assert"
)

func TestCollectNews(t *testing.T) {
	// 初始化测试数据
	domain := &news.NewsDomain{
		Scylla: &gocql.ClusterConfig{
			Hosts: []string{"localhost"},
		},
	}

	// 测试收藏功能
	err := domain.CollectNews(context.Background(), "1", "user1")
	assert.NoError(t, err, "收藏帖子应该成功")
}

func TestUncollectNews(t *testing.T) {
	// 初始化测试数据
	domain := &news.NewsDomain{
		Scylla: &gocql.ClusterConfig{
			Hosts: []string{"localhost"},
		},
	}

	// 测试取消收藏功能
	err := domain.UncollectNews(context.Background(), "1", "user1")
	assert.NoError(t, err, "取消收藏应该成功")
}

func TestGetNewsBookmarks(t *testing.T) {
	// 初始化测试数据
	domain := &news.NewsDomain{
		Scylla: &gocql.ClusterConfig{
			Hosts: []string{"localhost"},
		},
	}

	// 测试查询收藏列表
	bookmarks, err := domain.GetNewsBookmarks(context.Background(), "user1")
	assert.NoError(t, err, "查询收藏列表应该成功")
	assert.NotNil(t, bookmarks, "收藏列表不应为nil")
}
