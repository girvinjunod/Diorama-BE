CREATE DATABASE diorama;

\c diorama

CREATE USER diorama WITH PASSWORD 'diorama';
GRANT pg_read_all_data TO diorama;
GRANT pg_write_all_data TO diorama;

CREATE TABLE users (
	id SERIAL primary key,
    username VARCHAR (30) not null unique,
    email VARCHAR (50) not null unique,
    name VARCHAR (50),
    password TEXT not null,
    profile_picture BYTEA
);

INSERT INTO users (username, email, name, password) 
VALUES ('girvinjunod', 'girvinjunod@gmail.com', 'Girvin Junod', 'aaaaaa');

INSERT INTO users (username, email, name, password) 
VALUES ('kerangajaib', 'kerang@gmail.com', 'Kerang Ajaib', 'aaaaaa');

CREATE TABLE following (
    follower_id int not null,
    followed_id int not null,
    PRIMARY KEY (follower_id, followed_id),
    CONSTRAINT fk_follower
      FOREIGN KEY(follower_id) 
	  REFERENCES users(id)
	  ON DELETE CASCADE,
    CONSTRAINT fk_followed
      FOREIGN KEY(followed_id) 
	  REFERENCES users(id)
	  ON DELETE CASCADE
);

INSERT INTO following (follower_id, followed_id) 
VALUES (1,2);

INSERT INTO following (follower_id, followed_id) 
VALUES (2,1);

CREATE TABLE trips (
    id SERIAL primary key,
    user_id INT not null,
    start_date DATE,
    end_date DATE,
    trip_name VARCHAR (50),
    location_name VARCHAR (50),
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(id)
	  ON DELETE CASCADE
);

CREATE TABLE events (
    id SERIAL primary key,
    trip_id INT not null,
    user_id INT not null,
    caption VARCHAR (1000),
    event_date DATE,
    post_time TIMESTAMP,
    picture BYTEA,
    CONSTRAINT fk_trip
      FOREIGN KEY(trip_id) 
	  REFERENCES trips(id)
	  ON DELETE CASCADE,
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(id)
	  ON DELETE CASCADE
);

CREATE TABLE comments(
    id SERIAL primary key,
    event_id INT not null,
    user_id INT not null,
    text VARCHAR (500),
    comment_time TIMESTAMP,
    CONSTRAINT fk_event
      FOREIGN KEY(event_id) 
	  REFERENCES events(id)
	  ON DELETE CASCADE,
    CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES users(id)
	  ON DELETE CASCADE
);

pg_dump diorama > db/diorama.sql -U diorama -h localhost

ssh -i ~/.ssh/gcloud g13519097@34.101.123.15
./cloud_sql_proxy --instances=stei-itb-13519097:asia-southeast2:diorama-db=tcp:3306     --credential_file=key.json &>     proxy.log &