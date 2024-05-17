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

