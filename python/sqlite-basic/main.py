import sys
import sqlite3
import random

conn = sqlite3.connect("test.db")

table_name = "TestTable"

cursor = conn.cursor()

try:
    cursor.execute(f"""
        CREATE TABLE IF NOT EXISTS {table_name} (
            Id INTEGER PRIMARY KEY AUTOINCREMENT,
            Name TEXT NOT NULL,
            Age INTEGER
        )
    """)

    names = ["Alice", "Bob", "Charlie", "Dan", "Ellie", "Francesca", "Garry", "Holly", "Ian", "Jane"]

    for name in names:
        age = random.randint(20, 70)
        cursor.execute(f"INSERT INTO {table_name} (Name, Age) VALUES (?, ?)", (name, age))

    conn.commit()

    cursor.execute(f"SELECT * FROM {table_name}")
    rows = cursor.fetchall()
except Exception as e:
    print(e)
    sys.exit(1)
finally:
    conn.close()

for row in rows:
    print(row)

