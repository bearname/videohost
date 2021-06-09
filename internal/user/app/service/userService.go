package service

import (
	"encoding/json"
	"errors"
	"github.com/bearname/videohost/internal/user/app/dto"
	"github.com/bearname/videohost/internal/user/domain/model"
	"github.com/bearname/videohost/internal/user/domain/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
)

type UserService struct {
	userRepo   repository.UserRepo
	followRepo repository.FollowerRepo
}

func NewUserService(userRepo repository.UserRepo, followRepo repository.FollowerRepo) *UserService {
	s := new(UserService)
	s.followRepo = followRepo
	s.userRepo = userRepo

	return s
}

func (s *UserService) Find(usernameOrId string) (dto.FindUserDto, error) {
	var user model.User
	var err error
	uuid := s.isUUID(usernameOrId)

	if uuid {
		user, err = s.userRepo.FindById(usernameOrId)
	} else {
		user, err = s.userRepo.FindByUserName(usernameOrId)
	}
	if err != nil {
		return dto.FindUserDto{}, errors.New("user not exist")
	}

	return dto.FindUserDto{Username: usernameOrId,
		Email:        user.Email,
		IsSubscribed: user.IsSubscribed,
		Role:         user.Role.Values(),
	}, nil
}

func (s *UserService) isUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func (s *UserService) Follow(followingToUserId string, follower string, isFollowing bool) error {
	return s.followRepo.Follow(followingToUserId, follower, isFollowing)
}

func (s *UserService) GetUserStatistic(userId string) (*model.UserStatistic, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return nil, err
	}
	stats, err := s.followRepo.GetStats(userId)
	if err != nil {
		return nil, err
	}
	return model.NewUserStatistic(userId, user.Username, stats), nil
}

func (s *UserService) UpdatePassword(request *http.Request) error {
	var userDto dto.ChangePasswordUserDto
	err := json.NewDecoder(request.Body).Decode(&userDto)
	if err != nil {
		log.Error(err)
		return errors.New("cannot decode username/password struct")
	}
	if !IsUsernameContextOk(userDto.Username, request) {
		log.Error(err)
		return errors.New("is username context invalid")
	}

	userFromDb, err := s.userRepo.FindByUserName(userDto.Username)
	if (err == nil && userFromDb.Username != userDto.Username) || err != nil {
		log.Error(err)
		return errors.New("user not exist")
	}

	err = bcrypt.CompareHashAndPassword(userFromDb.Password, []byte(userDto.OldPassword))
	if err != nil {
		log.Error(err)
		return errors.New("wrong password")
	}

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(userDto.NewPassword), bcrypt.DefaultCost)

	if ok := s.userRepo.UpdatePassword(userDto.Username, passwordHash); !ok {
		log.Error(err)
		return errors.New("failed update password")
	}

	return nil
}
