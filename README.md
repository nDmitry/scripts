Scripts I'm using for my pet-projects where some full-scale monitoring or backup solutions feel like an overkill. All of them are written in Go.

## freespace

Checks available disk space and sends a notification to a Telegram chat if it's less than 10%. Encountered errors will be sent to a Telegram chat as well.

I don't like setting up email notifications for it is more cumbersome, but it's easy to add any API client you need in the `pkg/notifier` package.

Build the script with `make build-freespace` or a custom build command, the artifact will be located in the `./bin` directory. Upload it to a server and run periodically using cron, e.g.:

```cron
0 * * * * cd /where/you/uploaded/it && ./freespace -telegram-token=<bot token> -telegram-chat-id=<...> 2>&1 | /usr/bin/logger -t freespace
```

## pgdumpdoc

Runs pg_dump inside a specified Docker container and uploads the backup file to any S3 compatible object storage. The file will then be deleted from the host filesystem (or left untouched in case of an error). Encountered errors will be sent to a Telegram chat.

Build the script with `make build-pgdumpdoc` or a custom build command, the artifact will be located in the `./bin` directory. A cron job for it might look like:

```
0 3 * * * cd /where/you/uploaded/it && ./pgdumpdoc -container=my_postgres -u=postgres -db=my_db -o=/home/user/backups/my_db.pgdata -s3-endpoint=<S3 URL> -s3-key-id=<S3 key ID> -s3-key-secret=<S3 key secret> -s3-bucket=my-backups -s3-object=<my_db.pgdata> -telegram-token=<bot token> -telegram-chat-id=<...> 2>&1 | /usr/bin/logger -t pgdumpdoc
```

To restore from a backup download it from the object storage and run:

```bash
docker exec -i my_pg pg_restore --dbname=my_db --verbose --clean < /tmp/pg_dump.pgdata
```

## dirbackup

Creates a `tar.gz` archive with an arbitrary directory and uploads it to any S3 compatible object storage. The archive will then be deleted from the host filesystem (or left untouched in case of an error). Encountered errors will be sent to a Telegram chat.

Build the script with `make build-dirbackup` or a custom build command, the artifact will be located in the `./bin` directory. A cron job for it might look like:

```
0 4 * * * cd /where/you/uploaded/it && ./dirbackup -dir=/path/to/backup -o=/home/user/backups/backup.tar.gz -s3-endpoint=<S3 URL> -s3-key-id=<S3 key ID> -s3-key-secret=<S3 key secret> -s3-bucket=my-backups -s3-object=<my_db.pgdata> -telegram-token=<bot token> -telegram-chat-id=<...> 2>&1 | /usr/bin/logger -t dirbackup
```
