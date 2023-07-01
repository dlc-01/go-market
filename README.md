
# go-market

 Проект маркета с использованием ситемы лояльности. 

# Техническое задание

Подробное ТЗ находится в файле [tz.md](doc/tz.md).

# Особенности реализации

1. В качестве СУБД используется PostgreSQL.
3. Авторизация реализована с помощью JWT-токенов в cookies.
4. Хранение паролей в БД пользователей организовано с помощью bcrypt.
5. Используется фраемврок GIN, для работы с БД pgx, для логирование zap, для тестирования testify/mock.


