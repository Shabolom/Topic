package migrate

import (
	"Arkadiy_Servis_authorization/config"
	"Arkadiy_Servis_authorization/iternal/domain"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/gormigrate.v1"
)

// Migrate запустите миграцию для всех объектов и добавьте для них ограничения
// создаем таблицы и закидываем в бд тут
func Migrate() {
	db := config.DB

	userID, _ := uuid.NewV4()
	likeID, _ := uuid.NewV4()
	dizLikeID, _ := uuid.NewV4()
	topicID, _ := uuid.NewV4()
	massageID, _ := uuid.NewV4()

	// создаем объект миграции данная строка всегда статична (всегда такая)
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			// id всех миграций кторые были проведены
			ID: userID.String(),
			// переписываем так при создании таблицы изменяется только структура которую мы передаем
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.User{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("users").Error
				if err != nil {
					return err
				}
				return nil
			},
		}, {
			// id всех миграций кторые были проведены
			ID: massageID.String(),
			// переписываем так при создании таблицы изменяется только структура которую мы передаем
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.Massage{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("massages").Error
				if err != nil {
					return err
				}
				return nil
			},
		}, {
			// id всех миграций кторые были проведены
			ID: likeID.String(),
			// переписываем так при создании таблицы изменяется только структура которую мы передаем
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.Like{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("likes").Error
				if err != nil {
					return err
				}
				return nil
			},
		}, {
			// id всех миграций кторые были проведены
			ID: dizLikeID.String(),
			// переписываем так при создании таблицы изменяется только структура которую мы передаем
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.DizLike{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("diz_like").Error
				if err != nil {
					return err
				}
				return nil
			},
		}, {
			// id всех миграций кторые были проведены
			ID: topicID.String(),
			// переписываем так при создании таблицы изменяется только структура которую мы передаем
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&domain.Topic{}).Error
				if err != nil {
					return err
				}
				return nil
			},
			// это метод отмены миграции ни разу не использовал
			Rollback: func(tx *gorm.DB) error {
				err := tx.DropTable("topics").Error
				if err != nil {
					return err
				}
				return nil
			},
		},
	})

	err := m.Migrate()
	if err != nil {
		log.WithField("component", "migration").Panic(err)
	}

	if err == nil {
		log.WithField("component", "migration").Info("Migration did run successfully")
	} else {
		log.WithField("component", "migration").Infof("Could not migrate: %v", err)
	}
}
