CREATE DATABASE ourproject;

-- Creating tables
CREATE TABLE groups (
    group_id INTEGER PRIMARY KEY,
    name VARCHAR(25) NOT NULL,
    num_of_members INTEGER NOT NULL,
    launch_date DATE DEFAULT CURRENT_DATE
);

CREATE TABLE singer (
    singer_id INTEGER PRIMARY KEY,
    first_name VARCHAR(25) NOT NULL,
    last_name VARCHAR(25) NOT NULL,
    birthday DATE NOT NULL,
    group_id INTEGER REFERENCES groups(group_id)
);

CREATE TABLE album (
    album_id INTEGER PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    genre VARCHAR(25) NOT NULL,
    num_of_tracks INTEGER NOT NULL,
    group_id INTEGER REFERENCES groups(group_id) NOT NULL
);

CREATE TABLE song (
    song_id INTEGER PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    length INTEGER,
    album_id INTEGER REFERENCES album(album_id) NOT NULL
);


-- Inserting initial values for created tables

INSERT INTO groups(group_id, name, num_of_members)
VALUES (1, 'BTS', 7),
       (2, 'BlackPink', 4),
       (3, 'EXO', 9);

INSERT INTO singer(singer_id, first_name, last_name, birthday, group_id)
VALUES (1, 'Nam-joon', 'Kim','12-09-1994',1),
       (2,'Ji-min', 'Park','13-10-1995',1),
       (3,'Yoon-gi', 'Min','09-03-1993',1),
       (4,'Ji-soo', 'Kim','03-01-1995',2),
       (5,'Jennie', 'Kim','16-01-1996',2),
       (6,'Lisa', 'Manoban','27-03-1997',2),
       (7,'Roseanne', 'Park','11-02-1997',2),
       (8,'Baek-hyun', 'Byun','06-05-1992',3),
       (9,'Chan-yeol', 'Park','27-11-1992',3),
       (10,'Min-seok', 'Kim','26-03-1990',3);

INSERT INTO album(album_id, title, genre, num_of_tracks, group_id)
VALUES (100,'Map of the Soul', 'EDM',16,1),
       (101,'Born Pink','Pop',8,2),
       (102,'Dont mess up my Tempo', 'Dance',11,3);

INSERT INTO song(song_id, title, length, album_id)
VALUES (111,'ON',250,100),
       (112,'Pink Venom',Null,101),
       (113,'Love Shot', NULL,102);