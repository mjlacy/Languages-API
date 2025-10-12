FROM mariadb:latest

ENV MYSQL_DATABASE languages

COPY ./setup.sql /docker-entrypoint-initdb.d/
