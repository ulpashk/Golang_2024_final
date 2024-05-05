CREATE TABLE groups (
    group_id SERIAL PRIMARY KEY,
    name VARCHAR(25) NOT NULL,
    num_of_members INTEGER NOT NULL
);

CREATE TABLE album (
    album_id SERIAL PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    genre VARCHAR(25) NOT NULL,
    tracks INTEGER NOT NULL,
    group_id INTEGER REFERENCES groups(group_id) NOT NULL
);

CREATE TABLE songs (
    song_id SERIAL PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    length INTEGER,
    album_id INTEGER REFERENCES album(album_id) NOT NULL
);