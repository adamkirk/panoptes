package postgres

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/adamkirk/panoptes/internal/domain/users"
	"github.com/adamkirk/panoptes/internal/repository/postgres/schema/panoptes/public/model"
	"github.com/adamkirk/panoptes/internal/repository/postgres/schema/panoptes/public/table"
	"github.com/adamkirk/panoptes/internal/util"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
)

type UsersRepository struct {
	conn *Connector
}

type dbUser struct {
	model.Users

	Roles []model.Roles
}

func (r *UsersRepository) Create(u *users.User) error {
	conn, err := r.conn.Connection()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	stmt := table.Users.INSERT(table.Users.ID, table.Users.Email, table.Users.FirstName, table.Users.LastName, table.Users.Password).
	MODEL(model.Users{
			ID: u.ID,
			Email: u.Email,
			FirstName: u.FirstName,
			LastName: u.LastName,
			Password: u.PasswordHash,
		})

	if _, err := stmt.Exec(tx); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return errors.New(fmt.Sprintf("%s: %s", txErr.Error(), err.Error()))
		}

		return err
	}

	for _, role := range u.Roles {
		rolesStmt := table.UserRoles.INSERT(table.UserRoles.ID, table.UserRoles.RoleID, table.UserRoles.UserID).
			VALUES(uuid.New(), role.ID, u.ID)

		if _, err := rolesStmt.Exec(tx); err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				return errors.New(fmt.Sprintf("%s: %s", txErr.Error(), err.Error()))
			}
			return err
		}
	}

	return tx.Commit()
}

func (r *UsersRepository) Get(id uuid.UUID) (*users.User, error) {
	conn, err := r.conn.Connection()

	if err != nil {
		return nil, err
	}

	stmt := table.GithubWebhooks.SELECT(table.Users.AllColumns, table.Roles.AllColumns, table.Permissions.AllColumns).
		FROM(table.Users.
			LEFT_JOIN(table.UserRoles, table.UserRoles.UserID.EQ(table.Users.ID)).
			LEFT_JOIN(table.Roles, table.Roles.ID.EQ(table.UserRoles.RoleID)).
			LEFT_JOIN(table.RolesPermissions, table.Roles.ID.EQ(table.RolesPermissions.RoleID)).
			LEFT_JOIN(table.Permissions, table.RolesPermissions.PermissionID.EQ(table.Permissions.ID)),
		).
		WHERE(table.Users.ID.EQ(postgres.UUID(id)))

	slog.Debug("query.user.get", "query", stmt.DebugSql())
	dest := []dbUser{}
	if err := stmt.Query(conn, &dest); err != nil {
		return nil, err
	}

	if len(dest) == 0 {
		return nil, nil
	}

	u := dest[0]

	return &users.User{
		ID: u.ID,
		FirstName: u.FirstName,
		LastName: u.LastName,
		Email: u.Email,
		PasswordHash: u.Password,
		Roles: util.Map[model.Roles, *users.Role](func (v model.Roles) *users.Role {
			return &users.Role{
				ID: v.ID,
				Name: v.Name,
			}
		}, u.Roles),
	}, nil
}

func (r *UsersRepository) ByEmail(email string) (*users.User, error) {
	conn, err := r.conn.Connection()

	if err != nil {
		return nil, err
	}

	q := table.GithubWebhooks.SELECT(table.Users.AllColumns).
		FROM(table.Users).
		WHERE(table.Users.Email.EQ(postgres.String(email))).
		LIMIT(1)

	dest := []model.Users{}
	if err := q.Query(conn, &dest); err != nil {
		return nil, err
	}

	if len(dest) == 0 {
		return nil, nil
	}

	u := dest[0]

	return &users.User{
		ID: u.ID,
		FirstName: u.FirstName,
		LastName: u.LastName,
		Email: u.Email,
		PasswordHash: u.Password,
	}, nil
}

func NewUsersRepository(conn *Connector) *UsersRepository {
	return &UsersRepository{
		conn: conn,
	}
}