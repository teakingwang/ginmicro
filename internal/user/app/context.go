package app

import (
	"github.com/teakingwang/ginmicro/internal/user/model"
	"github.com/teakingwang/ginmicro/internal/user/repository"
	"github.com/teakingwang/ginmicro/internal/user/service"
	"github.com/teakingwang/ginmicro/pkg/datastore"
	"github.com/teakingwang/ginmicro/pkg/db"
	"github.com/teakingwang/ginmicro/pkg/logger"
	"gorm.io/gorm"
	"time"
)

type AppContext struct {
	DB          *gorm.DB
	Redis       datastore.Store
	UserService service.UserService
}

func NewAppContext() (*AppContext, error) {
	// 初始化 DB
	gormDB, err := db.NewDB()
	if err != nil {
		logger.Errorf("Failed to initialize database: %v", err)
		return nil, err
	}

	// 初始化 Redis
	redisStore := datastore.NewRedisClient()
	if err != nil {
		logger.Errorf("Failed to initialize Redis: %v", err)
		return nil, err
	}

	// 初始化 UserService
	userRepo := repository.NewUserRepo(gormDB)
	if err := userRepo.Migrate(); err != nil {
		logger.Errorf("Failed to migrate user repository: %v", err)
		return nil, err
	}

	// 插入初始订单记录
	if err := seedInitialUser(gormDB); err != nil {
		logger.Warn("failed to seed initial user:", err)
	}

	userSrv := service.NewUserService(userRepo, redisStore)

	return &AppContext{
		DB:          gormDB,
		Redis:       redisStore,
		UserService: userSrv,
	}, nil
}

func seedInitialUser(db *gorm.DB) error {
	user := &model.User{
		UserID:    100000000000000001, // 使用你的 Snowflake 或其他 ID 生成器
		Username:  "guest",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}
