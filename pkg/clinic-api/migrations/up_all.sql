CREATE TABLE IF NOT EXISTS clinics
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT                        NOT NULL,
    city       TEXT                        NOT NULL,
    address    TEXT                        NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT now(),
    updated_at TIMESTAMP(0) with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS patients
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT                        NOT NULL,
    birthdate  DATE                        NOT NULL,
    gender     TEXT                        NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT now(),
    updated_at TIMESTAMP(0) with time zone NOT NULL DEFAULT now(),
    user_id    INT REFERENCES users (id)
);


CREATE TABLE IF NOT EXISTS doctors
(
    id         BIGSERIAL PRIMARY KEY,
    clinic_id  BIGINT                      NOT NULL REFERENCES clinics (id),
    name       TEXT                        NOT NULL,
    specialty  TEXT                        NOT NULL,
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT now(),
    updated_at TIMESTAMP(0) with time zone NOT NULL DEFAULT now(),
    user_id    INT REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS appointments
(
    id         BIGSERIAL PRIMARY KEY,
    patient_id BIGINT                      NOT NULL,
    doctor_id  BIGINT                      NOT NULL,
    date       DATE                        NOT NULL,
    start_time TIME                        NOT NULL,
    end_time   TIME                        NOT NULL,
    status     TEXT                        NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients (id),
    FOREIGN KEY (doctor_id) REFERENCES doctors (id),
    created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT now(),
    updated_at TIMESTAMP(0) with time zone NOT NULL DEFAULT now()
);
--! 1 ends






-- add citext extension
CREATE EXTENSION IF NOT EXISTS citext;

-- Admin@kbtu.kz = admin@kbtu.kz
CREATE TABLE IF NOT EXISTS users
(
    id            BIGSERIAL PRIMARY KEY,
    created_at    TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name          TEXT                        NOT NULL,
    email         CITEXT UNIQUE               NOT NULL,
    password_hash BYTEA                       NOT NULL,
    activated     BOOL                        NOT NULL,
    version       INTEGER                     NOT NULL DEFAULT 1
);
--! 2 ends






CREATE TABLE IF NOT EXISTS tokens
(
    hash        BYTEA PRIMARY KEY,
    user_id     BIGINT                      NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry      TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope       TEXT                        NOT NULL,
    plain_token TEXT
);
--! 3 ends






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
--! 4 ends
