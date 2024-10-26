package users

import (
	"errors"
	"fmt"
	"time"

	"github.com/adamkirk/panoptes/internal/util/dt"
	"github.com/adamkirk/panoptes/internal/util/random"
	"github.com/google/uuid"
)

type AccessTokensRepo interface {
	Create(t *AccessToken) error
}
type AccessTokensService struct {
	encrypter Encrypter
	genString func(length int) (string)
	getNow func() time.Time
	users UsersRepo
	repo AccessTokensRepo
}

type TokenDefinition struct {
	ExpiryDays int
	UserID uuid.UUID
}

func (svc *AccessTokensService) Create(def TokenDefinition) (*AccessToken, error) {
	u, err := svc.users.Get(def.UserID)

	if err != nil {
		return nil, errors.New("user not found")
	}

	id := fmt.Sprintf("PAT_%s", svc.genString(16))
	secret:= svc.genString(32)

	hash, err := svc.encrypter.Encrypt(secret)

	if err != nil {
		return nil, err
	}

	if def.ExpiryDays < -1 || def.ExpiryDays == 0 {
		return nil, errors.New("expiry days must -1 or greater than 0")
	} 

	var expireAt *time.Time

	if def.ExpiryDays != -1 {
		now := svc.getNow().Add(time.Duration(def.ExpiryDays) * time.Hour * time.Duration(24))
		expireAt = &now
	}

	t := &AccessToken{
		ID: id,
		Secret: &secret,
		SecretHash: hash,
		ExpireAt: expireAt,
		User: u,
	}

	return t, svc.repo.Create(t)
}

type AccessTokensServiceOpt func(*AccessTokensService)

func NewAccessTokensService(encrypter Encrypter, users UsersRepo, repo AccessTokensRepo, opts... AccessTokensServiceOpt) *AccessTokensService {
	svc := &AccessTokensService{
		encrypter: encrypter,
		getNow: dt.NowUTC,
		genString: random.String,
		repo: repo,
		users: users,
	}

	// TODO add opts
	for _, opt := range opts {
		opt(svc)
	}

	return svc
}