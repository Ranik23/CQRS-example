#!/usr/bin/env python3

from clickhouse_driver import Client
import psycopg2
import os

# без dotenv, переменные приходят из docker-compose
pg_conn = psycopg2.connect(
    dbname=os.getenv("PG_DBNAME"),
    user=os.getenv("PG_USER"),
    password=os.getenv("PG_PASSWORD"),
    host=os.getenv("PG_HOST"),
    port=os.getenv("PG_PORT")
)

ch_client = Client(
    host=os.getenv("CH_HOST"),
    port=int(os.getenv("CH_PORT")),
    user=os.getenv("CH_USER"),
    password=os.getenv("CH_PASSWORD")
)

with pg_conn.cursor() as cur:
    cur.execute("SELECT id, name, created_at FROM users")
    rows = cur.fetchall()

ch_client.execute("""
    CREATE TABLE IF NOT EXISTS users (
        id UInt32,
        name String,
        created_at DateTime
    ) ENGINE = MergeTree()
    ORDER BY id
""")

ch_client.execute("INSERT INTO users VALUES", rows)
