# Medical Clinic REST API

## Project description:
Clinic API - это веб-приложение на языке Golang, предназначенное для управления данными клиники. Оно предоставляет набор эндпоинтов для выполнения различных операций с данными о врачах, записях на прием и пациентах.
Ключевые особенности: 
 1. Управление врачами:  
    * Создание новых врачей.  
    * Получение информации о врачах по идентификатору.  
    * Обновление информации о врачах.  
    * Удаление врачей.  
2. Управление записями на прием:  
    * Создание новых записей на прием.  
    * Получение информации о записях на прием по идентификатору.  
    * Обновление информации о записях на прием.  
    * Удаление записей на прием.  
3. Управление пациентами:  
    * Создание новых пациентов.  
    * Получение информации о пациентах по идентификатору.  
    * Обновление информации о пациентах.  
    * Удаление пациентов.  

---

####  Используемые технологии и библиотеки:  
  -  Язык программирования Golang; 
  -  Gorilla Mux;  
  -  Библиотека для работы с базой данных PostgreSQL: database/sql;  
  -  Драйвер для работы с базой данных PostgreSQL: github.com/lib/pq;  


### Team members:

```
Zhassulan Kainazarov, 22B030547
Damir Kakarov, 22B030548
Yerlan Kaliyev, 22B030373
```


## API endpoints:

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

## DB Structure:

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