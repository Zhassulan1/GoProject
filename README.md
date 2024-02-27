# Go_Project


## Medical Clinic REST API
```
POST /patients: Создать нового пациента.
GET /patients/:id: Получить информацию о пациенте по его ID.
PUT /patients/:id: Обновить информацию о пациенте по его ID.
DELETE /patients/:id: Удалить пациента по его ID.

POST /doctors: Создать нового врача.
GET /doctors/:id: Получить информацию о враче по его ID.
PUT /doctors/:id: Обновить информацию о враче по его ID.
DELETE /doctors/:id: Удалить врача по его ID.

POST /appointments: Создать новую запись на прием.
GET /appointments/:id: Получить информацию о записи на прием по ее ID.
PUT /appointments/:id: Обновить информацию о записи на прием по ее ID.
DELETE /appointments/:id: Удалить запись на прием по ее ID.
```

## DB Structure
```
Table patients {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  name text
  birthdate date
  gender text
}

Table doctors {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  name text
  specialty text
}

// many-to-many
Table appointments {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  patient_id bigint [ref: > patients.id]
  doctor_id bigint [ref: > doctors.id]
  date_time DATETIME
}
```