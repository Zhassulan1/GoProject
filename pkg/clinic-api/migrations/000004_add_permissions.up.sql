CREATE TABLE IF NOT EXISTS permissions
(
    id   BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions
(
    user_id       BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES ('doctors:read'),
       ('doctors:write'),
       ('patients:read'),
       ('patients:write'),
       ('appointments:read'),
       ('appointments:write'),
       ('clinics:read'),
       ('clinics:write');
