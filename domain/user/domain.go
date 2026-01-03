package user

import (
	"vc-go/dao"

	"github.com/redis/go-redis/v9"
)

type UserDomain struct {
	DB    *dao.DB
	Redis *redis.Client
}
