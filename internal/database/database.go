package database

import (
	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"log/slog"
	"os"
	"time"
)

func New(dsn string, env string) (*gorm.DB, func() error, error) {
	var gormConfig *gorm.Config
	if env == "dev" {
		gormConfig = &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // You can customize the logger output
				logger.Config{
					SlowThreshold: time.Second, // SQL queries that take longer than this threshold will be logged as slow queries
					LogLevel:      logger.Info, // Set log level to Log mode to log all queries
					Colorful:      true,        // Enable colorful output
				},
			),
		}
	} else {
		gormConfig = &gorm.Config{}
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	return db, sqlDb.Close, nil
}

func Initialize(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.Channel{})

	if db.Migrator().HasTable(&models.User{}) && db.Migrator().HasTable(&models.Channel{}) {
		var usersCount int64

		db.Model(&models.User{}).Count(&usersCount)

		users := []*models.User{
			{
				Name:     "user1",
				Email:    "user1@gmail.com",
				Salt:     "7d59ac12ccec1b47d0d25b045d5aae3f",
				Password: "$2a$10$o9oxvoVCz6Qdq3jyqOQPdu0LlpBU8rfwEQuxYDVcvjhxihOlN2uoG", // password1
			},
			{
				Name:     "user2",
				Email:    "user2@gmail.com",
				Salt:     "7985cc34b27ef2e7ad20d490c234c5b0",
				Password: "$2a$10$yY2X3qWyx.Xyq4D6anSFQOh6bu5StPZUKK77t4bKroSnFZa9k5GTq", // password2
			},
			{
				Name:     "user3",
				Email:    "user3@gmail.com",
				Salt:     "43b68be506c2eb1cc07b530aade6644f",
				Password: "$2a$10$xRgdxef7krGas/RbLNXStOBBR3KIZuCYRVBpmqjM0iEHRSjw3TvqO", //password3
			},
		}

		if usersCount == 0 {

			for _, user := range users {
				db.Create(user)
			}

			for _, user := range users {
				for _, userForAssociations := range users {
					if userForAssociations.ID != user.ID {
						err := db.Model(&user).Association("Contacts").Append(&models.User{ID: userForAssociations.ID})
						if err != nil {
							slog.Error("Association", "err", err)
						}
					}
				}
			}
		}

		var channelsCount int64

		db.Model(&models.Channel{}).Count(&channelsCount)

		if channelsCount == 0 {
			channels := []*models.Channel{
				{
					Name: "channel1",
				},
				{
					Name: "channel2",
				},
				{
					Name: "channel3",
				},
			}

			for _, channel := range channels {
				db.Create(&channel)
			}

			for _, user := range users {
				err := db.Model(&user).Association("Channels").Append(channels)
				if err != nil {
					slog.Error("Association", "err", err)
				}
			}
		}
	}

	return err
}
