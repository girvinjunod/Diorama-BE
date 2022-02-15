CREATE DATABASE diorama;

\c diorama

CREATE TABLE users (
	id SERIAL primary key,
    username VARCHAR (30) not null unique,
    email VARCHAR (50) not null unique,
    password TEXT not null,
    profile_picture TEXT
);

INSERT INTO users (username, email, password, profile_picture) 
VALUES ('girvinjunod', 'girvinjunod@gmail.com', 'aaaaaa', 'profile-picture/elephant-seal.jpg');

INSERT INTO users (username, email, password, profile_picture) 
VALUES ('kerangajaib', 'kerang@gmail.com', 'kerang', 'profile-picture/kerang.jpg');

CREATE USER diorama WITH PASSWORD 'diorama';
GRANT pg_read_all_data TO diorama;
GRANT pg_write_all_data TO diorama;
