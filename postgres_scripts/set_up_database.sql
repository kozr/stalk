-- Create a new database
CREATE DATABASE stalk;

-- Create a new user
CREATE USER default_superuser WITH ENCRYPTED PASSWORD 'default_password';

-- Grant privileges to your user on your database
GRANT ALL PRIVILEGES ON DATABASE stalk TO default_superuser;
