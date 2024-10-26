package postgres

import (
	"github.com/adamkirk/panoptes/internal/domain/users"
	"github.com/adamkirk/panoptes/internal/repository/postgres/schema/panoptes/public/model"
	"github.com/adamkirk/panoptes/internal/repository/postgres/schema/panoptes/public/table"
	"github.com/adamkirk/panoptes/internal/util"
	"github.com/go-jet/jet/v2/postgres"
)

type UserAccessTokensRepository struct {
	conn *Connector
}

type dbUserAccessTokenUserRole struct {
	model.Roles

	Permissions []model.Permissions
}

type dbUserAccessTokenUser struct {
	model.Users

	Roles []dbUserAccessTokenUserRole
}

type dbUserAccessToken struct {
	model.UserAccessTokens

	User dbUserAccessTokenUser
}

func (r *UserAccessTokensRepository) Create(t *users.AccessToken) error {
	conn, err := r.conn.Connection()

	if err != nil {
		return err
	}

	stmt := table.UserAccessTokens.INSERT(table.UserAccessTokens.ID, table.UserAccessTokens.UserID, table.UserAccessTokens.Secret, table.UserAccessTokens.ExpiresAt).
	MODEL(model.UserAccessTokens{
			ID: t.ID,
			UserID: t.User.ID,
			Secret: t.SecretHash,
			ExpiresAt: t.ExpireAt,
		})

	_, err = stmt.Exec(conn)

	return err
}

func (r *UserAccessTokensRepository) ByID(id string) (*users.AccessToken, error) {
	conn, err := r.conn.Connection()

	if err != nil {
		return nil, err
	}

	q := table.UserAccessTokens.SELECT(table.UserAccessTokens.AllColumns, table.Users.AllColumns, table.Roles.AllColumns, table.Permissions.AllColumns).
		FROM(
			table.UserAccessTokens.
			LEFT_JOIN(table.Users, table.UserAccessTokens.UserID.EQ(table.Users.ID)).
			LEFT_JOIN(table.UserRoles, table.Users.ID.EQ(table.UserRoles.UserID)).
			LEFT_JOIN(table.Roles, table.UserRoles.RoleID.EQ(table.Roles.ID)).
			LEFT_JOIN(table.RolesPermissions, table.Roles.ID.EQ(table.RolesPermissions.RoleID)).
			LEFT_JOIN(table.Permissions, table.RolesPermissions.PermissionID.EQ(table.Permissions.ID)),
		).WHERE(table.UserAccessTokens.ID.EQ(postgres.String(id)))

	dest := []dbUserAccessToken{}

	if err := q.Query(conn, &dest); err != nil {
		return nil, err
	}


	if len(dest) == 0 {
		return nil, nil
	}

	t := dest[0]

	return &users.AccessToken{
		ID: t.ID,
		SecretHash: t.Secret,
		ExpireAt: t.ExpiresAt,
		User: &users.User{
			ID: t.User.ID,
			FirstName: t.User.FirstName,
			LastName: t.User.LastName,
			Email: t.User.Email,
			PasswordHash: t.User.Password,

			Roles: util.Map[dbUserAccessTokenUserRole, *users.Role](func(in dbUserAccessTokenUserRole) *users.Role {
				return &users.Role{
					ID: in.ID,
					Name: in.Name,

					Permissions: util.Map[model.Permissions](func (in model.Permissions) *users.Permission {
						return &users.Permission{
							ID: in.ID,
							Name: in.Name,
						}
					}, in.Permissions),
				}
			}, t.User.Roles),
		},
	}, nil
}

func NewUserAccessTokensRepository(conn *Connector) *UserAccessTokensRepository {
	return &UserAccessTokensRepository{
		conn: conn,
	}
}