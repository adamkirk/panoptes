package users

import (
	"errors"
	"time"

	"github.com/adamkirk/panoptes/internal/domain/validation"
	"github.com/adamkirk/panoptes/internal/util/dt"
	"github.com/adamkirk/panoptes/internal/util/random"
	"github.com/google/uuid"
)

type CreateDTO struct {
	Email     string `validate:"required"`
	FirstName string `validate:"required"`
	LastName  string `validate:"required"`
	Password string `validate:"required"`
	Roles []string
}

type UsersRepo interface {
	ByEmail(email string) (*User, error)
	Create(u *User) error
	Get(id uuid.UUID) (*User, error)
}

type RolesRepo interface {
	ByNames(names []string) ([]*Role, error)
}

type UsersService struct {
	repo UsersRepo
	roles RolesRepo
	encrypter Encrypter
	genString func(length int) (string)
	getNow func() time.Time
	validator *validation.Validator
}

func (svc *UsersService) Create(dto CreateDTO) (*User, error) {
	if err := svc.validator.Validate(dto); err != nil {
		return nil, err
	}

	roles, err := svc.roles.ByNames(dto.Roles)

	if len(roles) != len(dto.Roles) {
		// TOOD improve error message
		return nil, errors.New("some roles we're not found")
	}

	if u, err := svc.repo.ByEmail(dto.Email); err != nil {
		return nil, err
	} else if u != nil {
		return nil, errors.New("email already in use")
	}

	passwordHash, err := svc.encrypter.Encrypt(dto.Password)

	if err != nil {
		return nil, err
	}

	u := &User{
		ID: uuid.New(),
		Email: dto.Email,
		FirstName: dto.FirstName,
		LastName: dto.LastName,
		PasswordHash: passwordHash,
		Roles: roles,
	}

	return u, svc.repo.Create(u)
}

type GetDTO struct {
	ID uuid.UUID `validate:"required"`
}

func (svc *UsersService) Get(dto GetDTO) (*User, error) {
	if err := svc.validator.Validate(dto); err != nil {
		return nil, err
	}

	return svc.repo.Get(dto.ID)
}

type UsersServiceOpt func(*UsersService)

func NewUsersService(encrypter Encrypter, repo UsersRepo, roles RolesRepo, validator *validation.Validator, opts... UsersServiceOpt) *UsersService {
	svc := &UsersService{
		encrypter: encrypter,
		getNow: dt.NowUTC,
		genString: random.String,
		repo: repo,
		validator: validator,
		roles: roles,
	}

	// TODO add opts
	for _, opt := range opts {
		opt(svc)
	}

	return svc
}