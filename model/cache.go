package model

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"go-message-pusher/wechat"
	"log"
	"sync"
	"time"
)

var RedisEnabled bool

var rdb *redis.Client

var userSMap sync.Map
var appTokenSMap sync.Map
var corpTokenSMap sync.Map

func InitCache(enableRedis bool) error {
	RedisEnabled = enableRedis
	if RedisEnabled {
		return initRedis()
	} else {
		return initMap()
	}
}

/* User cache */
/* 用户信息 4 小时后过期 */

func UpdateUserCache(u User) error {
	if RedisEnabled {
		return updateUserRedis(u)
	} else {
		return updateUserMap(u)
	}
}

func RetrieveUserCacheByName(name string) (User, error) {
	if RedisEnabled {
		return retrieveUserRedis(name)
	} else {
		return retrieveUserMap(name)
	}
}

func RetrieveAllUsersCache() (users []User, err error) {
	if RedisEnabled {
		return retrieveAllUsersRedis()
	} else {
		return retrieveAllUsersMap()
	}
}

func DeleteUserCacheByName(name string) error {
	if RedisEnabled {
		// 只用删除用户信息，accessToken 会自动过期
		return deleteUserRedis(name)
	} else {
		return deleteUserMap(name)
	}
}

/* AccessToken cache */
/* Redis 中的 access token 2 小时后过期，每次访问都会更新过期时间 */

func UpdateAllAccessTokens() error {
	users, err := RetrieveAllUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		appID, appSecret := u.App.AppID, u.App.AppSecret
		if appID != "" && appSecret != "" {
			r, err := wechat.GetAppAccessToken(appID, appSecret)
			if err != nil {
				log.Println(err.Error())
			} else if r.ErrCode != 0 {
				log.Println("Wechat Server:", r.ErrMsg)
			} else {
				_ = UpdateAppAccessTokenCache(u.Name, r.AccessToken)
			}
		}

		corpID, agentSecret := u.Corp.CorpID, u.Corp.AgentSecret
		if corpID != "" && agentSecret != "" {
			r, err := wechat.GetCropAccessToken(corpID, agentSecret)
			if err != nil {
				log.Println(err.Error())
			} else if r.ErrCode != 0 {
				log.Println("Wechat Server:", r.ErrMsg)
			} else {
				_ = UpdateCorpAccessTokenCache(u.Name, r.AccessToken)
			}
		}
	}

	return nil
}

func RetrieveAppAccessTokenCache(name string) (accessToken string, err error) {
	if RedisEnabled {
		ctx := context.Background()
		accessToken, err = rdb.GetEx(ctx, "appAccessToken:"+name, 2*time.Hour).Result()
		if err == redis.Nil {
			// 缓存未命中，重新查数据库并向微信服务器获取
			u, err := RetrieveUserCacheByName(name)
			if err != nil {
				return "", errors.New("user does not exist")
			}

			appID, appSecret := u.App.AppID, u.App.AppSecret
			if appID == "" || appSecret == "" {
				return "", errors.New("empty appID or appSecret")
			}

			r, err := wechat.GetAppAccessToken(appID, appSecret)
			if err != nil {
				return "", err
			}
			if r.ErrCode != 0 {
				return "", errors.New(r.ErrMsg)
			}

			// 更新缓存
			_ = UpdateAppAccessTokenCache(u.Name, r.AccessToken)
			return r.AccessToken, nil
		}
		return accessToken, err
	} else {
		accessToken, ok := appTokenSMap.Load(name)
		if !ok {
			// 缓存未命中，重新查数据库并向微信服务器获取
			u, err := RetrieveUserCacheByName(name)
			if err != nil {
				return "", errors.New("user does not exist")
			}

			appID, appSecret := u.App.AppID, u.App.AppSecret
			if appID == "" || appSecret == "" {
				return "", errors.New("empty appID or appSecret")
			}

			r, err := wechat.GetAppAccessToken(appID, appSecret)
			if err != nil {
				return "", err
			}
			if r.ErrCode != 0 {
				return "", errors.New(r.ErrMsg)
			}

			// 更新缓存
			_ = UpdateAppAccessTokenCache(u.Name, r.AccessToken)
			return r.AccessToken, nil
		}
		return accessToken.(string), nil
	}
}
func RetrieveCorpAccessTokenCache(name string) (accessToken string, err error) {
	if RedisEnabled {
		ctx := context.Background()
		accessToken, err = rdb.GetEx(ctx, "corpAccessToken:"+name, 2*time.Hour).Result()
		if err == redis.Nil {
			// 缓存未命中，重新查数据库并向微信服务器获取
			u, err := RetrieveUserCacheByName(name)
			if err != nil {
				return "", errors.New("user does not exist")
			}

			corpID, agentSecret := u.Corp.CorpID, u.Corp.AgentSecret
			if corpID == "" || agentSecret == "" {
				return "", errors.New("empty corpID or agentSecret")
			}

			r, err := wechat.GetCropAccessToken(corpID, agentSecret)
			if err != nil {
				return "", err
			}
			if r.ErrCode != 0 {
				return "", errors.New(r.ErrMsg)
			}

			// 更新缓存
			_ = UpdateCorpAccessTokenCache(u.Name, r.AccessToken)
			return r.AccessToken, nil
		}
		return accessToken, err

	} else {
		accessToken, ok := corpTokenSMap.Load(name)
		if !ok {
			// 缓存未命中，重新查数据库并向微信服务器获取
			u, err := RetrieveUserCacheByName(name)
			if err != nil {
				return "", errors.New("user does not exist")
			}

			corpID, agentSecret := u.Corp.CorpID, u.Corp.AgentSecret
			if corpID == "" || agentSecret == "" {
				return "", errors.New("empty corpID or agentSecret")
			}

			r, err := wechat.GetCropAccessToken(corpID, agentSecret)
			if err != nil {
				return "", err
			}
			if r.ErrCode != 0 {
				return "", errors.New(r.ErrMsg)
			}

			// 更新缓存
			_ = UpdateCorpAccessTokenCache(u.Name, r.AccessToken)
			return r.AccessToken, nil
		}
		return accessToken.(string), nil
	}
}

func UpdateAppAccessTokenCache(name, accessToken string) (err error) {
	if RedisEnabled {
		ctx := context.Background()
		err = rdb.Set(ctx, "appAccessToken:"+name, accessToken, 2*time.Hour).Err()
		return err
	} else {
		appTokenSMap.Store(name, accessToken)
		return nil
	}
}

func UpdateCorpAccessTokenCache(name, accessToken string) (err error) {
	if RedisEnabled {
		ctx := context.Background()
		err = rdb.Set(ctx, "corpAccessToken:"+name, accessToken, 2*time.Hour).Err()
		return
	} else {
		corpTokenSMap.Store(name, accessToken)
		return nil
	}
}

/* private functions */

func initRedis() error {
	ctx := context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	users, err := RetrieveAllUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		v, err := json.Marshal(u)
		if err != nil {
			return err
		}

		err = rdb.Set(ctx, "user:"+u.Name, v, 4*time.Hour).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func initMap() error {
	users, err := RetrieveAllUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		userSMap.Store(u.Name, u)
	}

	return nil
}

func retrieveUserRedis(name string) (user User, err error) {
	ctx := context.Background()
	val, err := rdb.GetEx(ctx, "user:"+name, 4*time.Hour).Result()

	switch err {
	case redis.Nil: // 缓存未命中，查数据库并更新缓存
		user, err = RetrieveUserByName(name)
		if err != nil {
			return
		}
		err = UpdateUserCache(user)
	case nil:
		err = json.Unmarshal([]byte(val), &user)
	}
	return
}

func retrieveUserMap(name string) (User, error) {
	if user, ok := userSMap.Load(name); ok {
		return user.(User), nil
	} else {
		user, err := RetrieveUserByName(name)
		if err != nil {
			return User{}, err
		}
		err = UpdateUserCache(user)
		return user, err
	}
}

func retrieveAllUsersRedis() (users []User, err error) {
	ctx := context.Background()
	iter := rdb.Scan(ctx, 0, "user:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		user := User{}
		err = json.Unmarshal([]byte(rdb.Get(ctx, key).Val()), &user)
		if err != nil {
			return
		}
		users = append(users, user)
	}
	err = iter.Err()
	return
}

func retrieveAllUsersMap() (users []User, err error) {
	userSMap.Range(func(k, v interface{}) bool {
		users = append(users, v.(User))
		return true
	})
	return
}

func updateUserRedis(u User) error {
	ctx := context.Background()

	v, err := json.Marshal(u)
	if err != nil {
		return err
	}

	rdb.Set(ctx, "user:"+u.Name, v, 4*time.Hour)
	return nil
}

func updateUserMap(u User) error {
	userSMap.Store(u.Name, u)
	return nil
}

func deleteUserRedis(name string) error {
	ctx := context.Background()
	err := rdb.Del(ctx, "user:"+name).Err()
	return err
}

func deleteUserMap(name string) error {
	userSMap.Delete(name)
	appTokenSMap.Delete(name)
	corpTokenSMap.Delete(name)
	return nil
}
