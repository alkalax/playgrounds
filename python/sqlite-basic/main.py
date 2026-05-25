import sqlite3
import random

conn = sqlite3.connect("test.db")

cursor = conn.cursor()

cursor.execute("""
    CREATE TABLE IF NOT EXISTS TestTable (
        Id INTEGER PRIMARY KEY AUTOINCREMENT,
        Name TEXT NOT NULL,
        Age INTEGER
    )
""")

names = ["Alice", "Bob", "Charlie", "Dan", "Ellie", "Francesca", "Garry", "Holly", "Ian", "Jane"]

for name in names:
    age = random.randint(20, 70)
    cursor.execute("INSERT INTO TestTable (Name, Age) VALUES (?, ?)", (name, age))

conn.commit()

cursor.execute("SELECT * FROM TestTable")
rows = cursor.fetchall()

for row in rows:
    print(row)

conn.close()
