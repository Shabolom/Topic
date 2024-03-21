package repository

import (
	"Arkadiy_Servis_authorization/config"
	"Arkadiy_Servis_authorization/iternal/domain"
	"Arkadiy_Servis_authorization/iternal/tools"
	"fmt"
	"github.com/gofrs/uuid"
	"net/http"
)

type UserRepo struct {
}

func NewUserRepo() *UserRepo {
	return &UserRepo{}
}

func (ur *UserRepo) Register(user domain.User) tools.UserResult {
	if _, err := ur.FindOnUser("login", user.Login); err == nil {
		return tools.UserResult{
			Message: "такой пользователь уже существует",
			Status:  http.StatusBadRequest,
		}
	}
	err := config.DB.
		Create(&user).
		Error

	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Err:     nil,
		Message: "успешно зарегестрированны",
		Status:  http.StatusCreated,
		Result:  user,
	}
}

func (ur *UserRepo) FindOnUser(what string, find any) (domain.User, error) {
	var user domain.User

	err := config.DB.
		Where(what+" =?", find).
		First(&user).
		Error

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (ur *UserRepo) JoinTopic(userID uuid.UUID, topicID string) (domain.User, error) {
	id, _ := uuid.FromString(topicID)

	err := config.DB.
		Model(&domain.Topic{
			Base: domain.Base{
				ID: id,
			},
		}).
		Association("Users").
		Append(domain.User{
			Base: domain.Base{
				ID: userID,
			},
		}).
		Error

	if err != nil {
		return domain.User{}, err
	}

	return domain.User{}, nil
}

func (ur *UserRepo) PostMassage(massage domain.Massage) tools.UserResult {

	err := config.DB.
		Create(&massage).
		Error

	if err != nil {
		return tools.UserResult{
			Err:    err,
			Status: http.StatusBadRequest,
		}
	}

	return tools.UserResult{
		Err:    nil,
		Status: http.StatusOK,
	}
}

func (ur *UserRepo) FindDiz(massageID uuid.UUID, userID uuid.UUID) (domain.Massage, error) {
	var massage domain.Massage

	err := config.DB.
		Model(domain.Massage{}).
		Where("id = ?", massageID).
		Preload("DizLikes", "user_id IN (?)", userID).
		Preload("Likes", "user_id IN (?)", userID).
		//Preload("DizLikes").
		//Preload("Likes").
		Find(&massage).
		Error

	if err != nil {
		return domain.Massage{}, err
	}

	return massage, nil
}

func (ur *UserRepo) CreateDizLike(dizLike domain.DizLike) error {

	err := config.DB.
		Create(&dizLike).
		Error

	return err
}

func (ur *UserRepo) DeleteDizLike(massageID uuid.UUID) error {
	var dizLike domain.DizLike

	err := config.DB.
		Where("massage_id =?", massageID).
		Delete(&dizLike).
		Error

	return err
}

func (ur *UserRepo) CreateLike(like domain.Like) error {

	err := config.DB.
		Create(&like).
		Error

	return err
}

func (ur *UserRepo) DeleteLike(massageID uuid.UUID) error {
	var like domain.Like

	err := config.DB.
		Where("massage_id =?", massageID).
		Delete(&like).
		Error

	return err
}

func (ur *UserRepo) FindUserInTopic(topicID uuid.UUID, userID uuid.UUID) error {
	var topic domain.Topic

	err := config.DB.
		Model(&domain.Topic{}).
		Preload("Users", "topic_id IN (?)", topicID).
		First(&topic).
		Error

	return err
}

func (ur *UserRepo) TopicMassages(topicID uuid.UUID) (domain.Topic, error) {
	var topic domain.Topic
	fmt.Println(topicID)
	err := config.DB.
		Model(&domain.Topic{}).
		Order("created-at").
		Preload("Massages", "topic_id IN (?)", topicID).
		Preload("Massages.Likes").
		Preload("Massages.DizLikes").
		Find(&topic).
		Error

	if err != nil {
		return domain.Topic{}, err
	}

	return topic, err
}
