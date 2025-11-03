package permission

import "context"

func (r *repository) HasPermission(ctx context.Context, roleName string, resource string, action string) (bool, error) {

	query := `
		SELECT EXISTS( 
			SELECT 1
			FROM role_permissions rp
			JOIN roles r ON r.id = rp.role_id
			JOIN permissions p ON p.id = rp.permission_id
			WHERE r.name = $1 
				AND p.resource = $2
				AND p.action = $3
		)`

	var exists bool

	err := r.Pool.QueryRow(ctx, query, roleName, resource, action).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil

}
