# api_tickets
## comandos para go 
```
go mod init api_tickets
go get -u github.com/gorilla/mux
go get -u github.com/go-sql-driver/mysql
ejecutar el servidor de go en local
go run main.go
```

## comandos para crear la base de datos y la tabla en la base de datos mysql/mariadb en local

```
CREATE DATABASE db_ticket;

CREATE USER 'user_ticket'@'localhost' IDENTIFIED BY '12345';

GRANT ALL ON db_ticket.* TO 'user_ticket'@'localhost';
FLUSH PRIVILEGES;

USE db_ticket;

CREATE TABLE tickets (id int AUTO_INCREMENT, user varchar(50), date_create date, date_update date, state set('abierto','cerrado'), PRIMARY KEY(id));
```