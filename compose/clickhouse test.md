```bash
docker exec -it clickhouse clickhouse local -m
```

```bash
create database test;
use test;
```

```bash
CREATE TABLE test.donky(
    id UInt32
    blob String,
    meta Tuple(uri String, id String, stream String, topic String, domain String)
)
ENGINE = Kafka('redpanda-0:9092', 'donkey-locations', 'click-test', 'JSONEachRow')
```

```bash
CREATE table wiki (
id UInt32,
blob String,
topic String,
)
ENGINE = MergeTree
Order by id;
```

```bash
 Create materialized view test_mv to test.donky as
select * from wiki
SETTINGS
stream_like_engine_allow_direct_select = 1;
```

```bash

```

```bash

```

```bash

```
