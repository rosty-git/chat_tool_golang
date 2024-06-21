package database

import (
	"log"
	"os"
	"time"

	"github.com/elef-git/chat_tool_golang/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config interface {
	GetEnv() string
	GetDsn() string
	GetGormDebug() bool
}

func New(c Config) (*gorm.DB, func() error, error) {
	var gormConfig *gorm.Config
	//if c.GetGormDebug() {
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
	//} else {
	//	gormConfig = &gorm.Config{}
	//}

	db, err := gorm.Open(mysql.Open(c.GetDsn()), gormConfig)
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
	err := db.AutoMigrate(&models.User{}, &models.Channel{}, &models.ChannelMembers{}, &models.Post{}, &models.Status{}, &models.File{}, &models.TmpFile{})

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
		}

		var channelsCount int64

		db.Model(&models.Channel{}).Count(&channelsCount)

		if channelsCount == 0 {
			channels := []*models.Channel{
				{
					Name: "Channel 1",
					Type: models.ChannelTypeOpen,
				},
				{
					Name: "Channel 2",
					Type: models.ChannelTypeOpen,
				},
				{
					Name: "Channel 3",
					Type: models.ChannelTypeOpen,
				},
			}

			for _, channel := range channels {
				db.Create(&channel)

				for _, user := range users {
					channelMembers := models.ChannelMembers{
						UserID:       user.ID,
						ChannelID:    channel.ID,
						LastViewedAt: time.Now(),
					}
					db.Create(&channelMembers)
				}
			}

			for _, user1 := range users {
				for _, user2 := range users {
					if user1.ID != user2.ID {
						var directChannelsCount int64

						db.Where(models.Channel{Name: user1.ID + "__" + user2.ID}).Or(models.Channel{Name: user2.ID + "__" + user1.ID}).Find(&models.Channel{}).Count(&directChannelsCount)

						if directChannelsCount == 0 {
							newChannel := models.Channel{
								Type: models.ChannelTypeDirect,
								Name: user1.ID + "__" + user2.ID,
							}

							db.Create(&newChannel)

							channelMembers := models.ChannelMembers{
								UserID:       user1.ID,
								ChannelID:    newChannel.ID,
								LastViewedAt: time.Now(),
							}
							db.Create(&channelMembers)

							channelMembers1 := models.ChannelMembers{
								UserID:       user2.ID,
								ChannelID:    newChannel.ID,
								LastViewedAt: time.Now(),
							}
							db.Create(&channelMembers1)
						}
					}
				}
			}
		}
	}

	return err
}
