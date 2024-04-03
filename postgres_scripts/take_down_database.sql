-- UNDO All the changes made by the setup script

-- Drop the database
DROP DATABASE IF EXISTS stalk;

-- Drop the user
DROP USER IF EXISTS default_superuser;

-- Drop the tables
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS users;