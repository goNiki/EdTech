package permission

import (
	"context"
	"edtech/internal/infrastructure/txmanager"
	"fmt"
)

func (r *repository) GetRolePermissions(ctx context.Context, rolename string) ([]string, error) {
	const op = "repository.permission.getrolepermissions"
	q := txmanager.GetQueryExecutor(ctx, r.Pool)

	query := `
		SELECT permissions.resource || ':' || permissions.action as permisions
		FROM role_permissions
		JOIN roles ON roles.id = role_permissions.roles_id 
		JOIN permissions ON permissions.id = role_permissions.permission_id
		WHERE roles.name = $1
	`
	rows, err := q.Query(ctx, query, rolename)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	var permisions []string

	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		permisions = append(permisions, perm)
	}

	return permisions, nil

}
