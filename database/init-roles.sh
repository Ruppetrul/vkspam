#!/bin/bash

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
    DO \$\$ BEGIN
      IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'vkspam') THEN
                CREATE ROLE vkspam WITH LOGIN PASSWORD 'vkspam';
            END IF;
      GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO vkspam;
    END \$\$;
EOSQL
