Scripts I'm using for my pet-projects where some full-scale monitoring or backup solutions feel like an overkill. All of them are written in Go.

## freespace

Checks available disk space and sends a notification to a Telegram chat if it's less than 10%. I don't like setting up email notifications for it is more cumbersome, but it's easy to add a new API client you need in the `pkg/notifier` package.

Build the script with `make build-freespace` or a custom build command, the artifact will be located in the `./bin` directory. Upload it to a server and run periodically using cron, e.g.:

```cron
0 * * * * cd /where/you/uploaded/it && TELEGRAM_TOKEN=<your bot token> TELEGRAM_CHAT_ID=<...> ./freespace 2>&1 | /usr/bin/logger -t freespace
```

## pgdumpdoc

Runs pg_dump inside a specified Docker container and outputs the backup in a file on the host filesystem. Encountered errors will be sent to a Telegram chat.

Build the script with `make build-pgdumpdoc` or a custom build command, the artifact will be located in the `./bin` directory. A cron job for it might look like:

```
0 3 * * * cd /where/you/uploaded/it && CONTAINER=my_postgres USER=postgres DATABASE=my_db OUTFILE=/home/user/backups/my_db.pgdata TELEGRAM_TOKEN=<your bot token> TELEGRAM_CHAT_ID=<...> ./pgdumpdoc 2>&1 | /usr/bin/logger -t pgdumpdoc
```

To restore from a backup run:

```bash
docker exec -i my_pg pg_restore --dbname=my_db --verbose --clean < /tmp/pg_dump.pgdata
```
