CREATE TABLE IF NOT EXISTS VirtualMachines (
  Name TEXT PRIMARY KEY NOT NULL UNIQUE,
  ShutdownTime DATETIME NOT NULL,
  SkipToday BOOL NOT NULL DEFAULT 0
);

DELETE FROM VirtualMachines;

INSERT INTO VirtualMachines (Name, ShutdownTime, SkipToday) VALUES
  ('ubuntu01', '2024-06-30 22:00:00', 0),
  ('ubuntu02', '2024-06-30 22:00:00', 1),
  ('ubuntu03', '2024-06-30 20:00:00', 0),
  ('ubuntu04', '2024-06-30 19:00:00', 1);
