package v1

import (
	"context"
	"net/http"

	"github.com/adamkirk/panoptes/internal/api/operations"
	"github.com/adamkirk/panoptes/internal/domain/users"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)


type UsersService interface {
	Get(dto users.GetDTO) (*users.User, error)
}

type UsersController struct {
	svc UsersService
}

func (c *UsersController) RegisterRoutes(api huma.API) {
	huma.Register[GetUserRequest, GetUserResponse](api, huma.Operation{
		OperationID:  "v1.users.get",
		Method:       http.MethodGet,
		Path:         "/users/{id}",
		Summary:      "Get a User By ID",
		DefaultStatus: http.StatusOK,
		Metadata: map[string]any{
			operations.OptDisableNotFound: true,
		},
		Security: []map[string][]string{
			{"scopes": {"users.get"}},
		},
	}, ErrorHandler(true, c.Get))
}

func NewUsersController(svc UsersService) *UsersController {
	return &UsersController{
		svc: svc,
	}
}

type GetUserRequest struct {
	ID string `path:"id" required:"true"`
}

type GetUserResponse struct {
	Body *users.User
}

func (c *UsersController) Get(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {

	id, err := uuid.Parse(req.ID)

	if err != nil {
		return nil, huma.Error404NotFound(err.Error())
	}

	u, err := c.svc.Get(users.GetDTO{
		ID: id,
	})

	if err != nil {
		return nil, err
	}
	return &GetUserResponse{
		Body: u,
	}, nil
}
