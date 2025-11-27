#!/bin/bash

# Скрипт для бекапа PostgreSQL
# Запускать каждые 3 часа через cron

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups"
DB_NAME="products"
HOST="postgres-master"
USER="admin"

# Создаем бекап
pg_dump -h $HOST -U $USER -d $DB_NAME -F c -b -v -f "$BACKUP_DIR/backup_$DATE.dump"

# Удаляем бекапы старше 7 дней
find $BACKUP_DIR -name "*.dump" -mtime +7 -delete

echo "Backup completed: backup_$DATE.dump"