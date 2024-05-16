project
# Golang_2024_final
# K-pop Entertainment

## Members

```
Kuanyshbek Ulpan 22B030552
Kemel Merey 22B030615
```

## Songs REST API
```
POST /songs 
GET /songs/id
PUT /songs/id
DELETE /songs/
```

## DB Structure
```

TABLE groups (
    group_id INTEGER PRIMARY KEY,
    name VARCHAR(25) NOT NULL,
    num_of_members INTEGER NOT NULL,
    launch_date DATE DEFAULT CURRENT_DATE
);

TABLE singer (
    singer_id INTEGER PRIMARY KEY,
    first_name VARCHAR(25) NOT NULL,
    last_name VARCHAR(25) NOT NULL,
    birthday DATE NOT NULL,
    group_id INTEGER REFERENCES groups(group_id)
);

TABLE album (
    album_id INTEGER PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    genre VARCHAR(25) NOT NULL,
    num_of_tracks INTEGER NOT NULL,
    group_id INTEGER REFERENCES groups(group_id) NOT NULL
);

TABLE song (
    song_id INTEGER PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    length INTEGER,
    album_id INTEGER REFERENCES album(album_id) NOT NULL
);
```
