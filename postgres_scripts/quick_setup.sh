#! /bin/sh

echo "Checking if stalk database exists, if so, drop it"
if psql -U postgres -lqt | cut -d \| -f 1 | grep -qw stalk; then
  echo "Dropping stalk database"
  psql -U postgres -d stalk -a -f postgres_scripts/take_down_database.sql
fi

echo "Creating stalk database"
psql -U postgres -a -f postgres_scripts/set_up_database.sql

echo "Running all scripts with number prefixes in postgres_scripts"
for script in $(ls postgres_scripts | grep '^[0-9]'); do
  echo "Running $script"
  psql -U postgres -d stalk -a -f postgres_scripts/$script
done