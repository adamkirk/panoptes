package postgres

import (
	"github.com/adamkirk/panoptes/internal/domain/users"
	"github.com/adamkirk/panoptes/internal/repository/postgres/schema/panoptes/public/model"
	"github.com/adamkirk/panoptes/internal/repository/postgres/schema/panoptes/public/table"
	"github.com/adamkirk/panoptes/internal/util"
	"github.com/go-jet/jet/v2/postgres"
)

type RolesRepository struct {
	conn *Connector
}

type dbRole struct {
	model.Roles
 
	Permissions []model.Permissions
}

func (r *RolesRepository) ByNames(names []string) ([]*users.Role, error) {
	if len(names) == 0 {
		return []*users.Role{}, nil
	}
	conn, err := r.conn.Connection()

	if err != nil {
		return nil, err
	}

	namesIn := util.Map[string, postgres.Expression](func (v string) postgres.Expression {
		return postgres.String(v)
	}, names)

	stmt := table.Roles.SELECT(table.Roles.AllColumns).
		FROM(table.Roles.
			LEFT_JOIN(table.Permissions, table.Roles.ID.EQ(table.Permissions.ID))).
		WHERE(table.Roles.Name.IN(namesIn...))

	var dest []dbRole

	if err := stmt.Query(conn, &dest); err != nil {
		return nil, err
	}

	if len(dest) == 0 {
		return nil, nil
	}

	roles := util.Map[dbRole, *users.Role](func (in dbRole) *users.Role {
		return &users.Role{
			ID: in.ID,
			Name: in.Name,

			Permissions: util.Map[model.Permissions, *users.Permission](func (in model.Permissions) *users.Permission{
				return &users.Permission{
					ID: in.ID,
					Name: in.Name,
				}
			}, in.Permissions),
		}
	}, dest)
		
	return roles, nil
}

func NewRolesRepository(conn *Connector) *RolesRepository {
	return &RolesRepository{
		conn: conn,
	}
}