-- +goose Up
-- +goose StatementBegin

-- Админ получает все разрешения
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin';

-- Учитель получает разрешения на курсы, уроки и ресурсы
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'teacher' 
  AND p.resource IN ('course', 'lesson', 'resource');

-- Студент получает только разрешения на просмотр
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'student' 
  AND p.action = 'view' AND p.resource IN ('course', 'lesson', 'resource');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions WHERE role_id IN (SELECT id FROM roles WHERE name = 'admin');
DELETE FROM role_permissions WHERE role_id IN (SELECT id FROM roles WHERE name = 'teacher');
DELETE FROM role_permissions WHERE role_id IN (SELECT id FROM roles WHERE name = 'student');
-- +goose StatementEnd
