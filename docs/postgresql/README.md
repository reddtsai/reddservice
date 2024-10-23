# Postgresql 17

透過 docker compose 建立 replication

## Install

```
make docker-compose-up
```

## Stream Replication

Set Primary Database

- 身份驗證

  建立帳號

  ```
  CREATE ROLE replicator WITH REPLICATION LOGIN ENCRYPTED PASSWORD '1qaz2wsx';
  ```

  驗證控制，例如 docker network 分配了 172.30.0.0 網段，將此地址加入 pg_hba.conf

  ```
  # TYPE  DATABASE        USER            ADDRESS                 METHOD
  host    replication     all             172.30.0.0/16           md5
  ```

  > Client authentication for replication is controlled by a `pg_hba.conf` record specifying replication in the database field.

- 資料庫配置

  修改 postgresql.conf 中以下內容

  ```
  listen_addresses = '*'
  wal_level = replica
  max_wal_senders = 10
  wal_keep_size = 16MB
  synchronous_commit = on
  ```

  > 可透過 postgres -c max_wal_senders = 10 設定 `postgresql.conf`

  建立 slot

  ```
  SELECT * FROM pg_create_physical_replication_slot('replica_1_slot');
  ```

- 重啟資料庫

- 檢查從節點

  當 Standby Database 設定完成

  ```
  SELECT * FROM pg_stat_replication;
  ```

Set Standby Database

- 備份

  備份 Primary Database

  ```
  su postgres
  pg_basebackup -h auth-primary-db -D /var/lib/postgresql/data/replica -U replicator -X stream
  ```

- 資料庫配置

  修改 postgresql.conf 中以下內容

  ```
  listen_addresses = '*'
  wal_level = replica
  max_wal_senders = 10
  wal_keep_size = 16MB
  synchronous_commit = on
  primary_conninfo = 'host=auth-primary-db port=5432 user=replicator password=1qaz2wsx application_name=s1'
  primary_slot_name = 'replica_1_slot'
  ```

  > 可透過 postgres -c max_wal_senders = 10 設定 `postgresql.conf`

- 開始同步

  停止資料庫後，以下動作因為是使用 docker compose，可直接在主機操作

  ```
  # 以待命模式啟動 /var/lib/postgresql/data/standby.signal
  touch /data/replica/standby.signal
  # 移除或備分 old data /var/lib/postgresql/data/data
  mv /data/replica /data/replica-old
  mv /data/replica-old/replica /data/replica
  ```

- 重啟資料庫

  確認備份

  ```
  # 資料庫變成 read only
  database system is ready to accept read-only connections
  # 查看同步資料
  SELECT * FROM pg_stat_wal_receiver;
  ```
