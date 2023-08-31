# avito-tech-backend-trainee-assigment-2023
Запуск проекта:  
```
docker build . -t avito_service:latest
```
```
docker-compose up -d --force-recreate
```
Запуск sql-скрипта для создания таблиц:  
![Image](https://github.com/m0rk0vka/images/raw/main/avito-service.drawio.png)
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
Добавление пользователя в сегмент:  ђ
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
Уточнения:  
1. При получении активных сегментов пользователя и при добавлении пользователя в сегмент - новый пользователь добавляется в таблицу
2. При добавлении пользователя в сегмент, если пытаемся добавить пользователя повторно в сегмент/пытаемся добавить пользователя в несуществующий сегмент/пытаемся убрать пользователя из сегмента, в котором его нет, то запрос прерывается
3. По предыдущим пунктам: добавление пользователя вынес бы отдельно, так же добавление пользователя в сегмент/удаление пользователя из сегмента - выполнял бы по одному за раз, меньше ошибок при обработке и проще логику описать - в рабочей среде предложил бы такое решение, но здесь нет такой возможности, поэтому реализовал по тз 
