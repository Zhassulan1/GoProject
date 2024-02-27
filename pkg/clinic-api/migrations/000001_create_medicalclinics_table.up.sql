CREATE TABLE IF NOT EXISTS patients 
(
  id            bigserial PRIMARY KEY,
  created_at    timestamp(0) with time zone NOT NULL DEFAULT now(),
  updated_at    timestamp(0) with time zone NOT NULL DEFAULT now(),
  name          text                        NOT NULL,
  birthdate     date                        NOT NULL,
  gender        text                        NOT NULL
);

CREATE TABLE IF NOT EXISTS doctors
(
    id            bigserial PRIMARY KEY,
    created_at    timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at    timestamp(0) with time zone NOT NULL DEFAULT now(),
    name          text                        NOT NULL,
    specialty     text                        NOT NULL
);

CREATE TABLE IF NOT EXISTS appointments
(
    id            bigserial PRIMARY KEY,
    created_at    timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at    timestamp(0) with time zone NOT NULL DEFAULT now(),
    patient_id    bigint                      NOT NULL,
    doctor_id     bigint                      NOT NULL,
    date          date                        NOT NULL,
    start_time    time                        NOT NULL,
    end_time      time                        NOT NULL,
    status        text                        NOT NULL,
    FOREIGN KEY (patient_id) REFERENCES patients (id),
    FOREIGN KEY (doctor_id) REFERENCES doctors (id)
);