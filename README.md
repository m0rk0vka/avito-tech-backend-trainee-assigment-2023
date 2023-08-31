# avito-tech-backend-trainee-assigment-2023
Запуск проекта:  
```
docker build . -t avito_service:latest
```
```
docker-compose up -d --force-recreate
```
Запуск sql-скрипта для создания таблиц:
```
psql -U avito_service -d avito_service_db -W < avito_service_db.sql
```
Тестирование через postman:
Создание сегмента:
URL:  
```
http://localhost:8080/api/createsegment
```
Body:  
```
{
  "name": "SEGMENT_1"
}
```
Удаление сегмента:
URL:  
```
http://localhost:8080/api/deletesegment
```
Body:  
```
{
  "name": "SEGMENT_1"
}
```
Получение активных сегментов пользователя:
URL:  
```
http://localhost:8080/api/getusersegments
```
Body:  
```
{
  "id": 1010
}
```
Добавление пользователя в сегмент:
URL:  
```
http://localhost:8080/api/updateusersegments
```
Body:  
```
{
    "segments_to_add": ["SEGMENT_1", "SEGMENT_2"],
    "segments_to_delete": ["SEGMENT_3", "SOME_SEGMENT_4"],
    "user_id": 1010
}
```
