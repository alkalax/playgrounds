import sqlite3

conn = sqlite3.connect("test.db")

cursor = conn.cursor()

cursor.execute("""
    CREATE TABLE IF NOT EXISTS TestTable (
        Id INTEGER PRIMARY KEY AUTOINCREMENT,
        Name TEXT NOT NULL,
        Age INTEGER
    )
""")

cursor.execute("INSERT INTO TestTable (Name, Age) VALUES (?, ?)", ("Bob", 34))

conn.commit()

cursor.execute("SELECT * FROM TestTable")
rows = cursor.fetchall()

for row in rows:
    print(row)

conn.close()
