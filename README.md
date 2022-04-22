#

```

go get github.com/jackc/pgx/v5

go run main.go

```



# Database configuration

```

    CREATE TABLE IF NOT EXISTS USERS(
        ID          SERIAL   PRIMARY KEY,
        USERNAME    VARCHAR(20) NOT NULL UNIQUE
    );
```


# Docker-compose config
```
    docker-compose up


```
