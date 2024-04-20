#!/bin/bash
set -e

sed -i 's/max_wal_size = 1GB/max_wal_size = 2GB/' /var/lib/postgresql/data/postgresql.conf
echo "shared_preload_libraries = 'pg_cron'" >> /var/lib/postgresql/data/postgresql.conf
echo "cron.use_background_workers = on" >> /var/lib/postgresql/data/postgresql.conf
