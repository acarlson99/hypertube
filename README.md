# hypertube
Basically..

## dependencies
- github.com/gorilla/mux
- github.com/lib/pq
- github.com/anacrolix/torrent

## Setup

### Database

install postgresql

Set up and run database

```
export PGDATA=/tmp/postgres/
initdb $PGDATA
postgres
```

Create hypertube user

```
createdb -O `whoami`
psql
create user postgres with superuser;
```
