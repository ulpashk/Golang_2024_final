project
# Golang_2024_final
# K-pop Entertainment

## Members

```
Kuanyshbek Ulpan 22B030552
Kemel Merey 22B030615
```

## Healthcheck
```
GET /healthcheck
```

## Groups REST API
```
# list of all groups items
GET /groups
POST /groups
GET /songs/:id
PUT /songs/:id
DELETE /songs/:id
GET groups/:id/albums
```

## Albums REST API
```
# list of all albums items
GET /albums
POST /albums
GET /albums/:id
PUT /albums/:id
DELETE /albums/:id
GET albums/:id/songs
```

## Songs REST API
```
# list of all songs items
GET /songs
POST /song
GET /songs/:id
PUT /songs/:id
DELETE /songs/:id
```

## Users
```
POST /users
PUT /users/activated
```
## Token
```
POST /token/login
```
## DB Structure
```

TABLE groups (
    group_id INTEGER PRIMARY KEY,
    name VARCHAR(25) NOT NULL,
    num_of_members INTEGER NOT NULL,
    launch_date DATE DEFAULT CURRENT_DATE
);

TABLE album (
    album_id INTEGER PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    genre VARCHAR(25) NOT NULL,
    num_of_tracks INTEGER NOT NULL,
    group_id INTEGER REFERENCES groups(group_id) NOT NULL
);

TABLE songs (
    song_id INTEGER PRIMARY KEY,
    title VARCHAR(25) NOT NULL,
    length INTEGER,
    album_id INTEGER REFERENCES album(album_id) NOT NULL
);

TABLE users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    email citext UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    activated bool NOT NULL,
    version integer NOT NULL DEFAULT 1
);

TABLE tokens (
    hash bytea PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);
```

