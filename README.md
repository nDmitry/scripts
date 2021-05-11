Scripts I'm using for my pet-projects where some full-scale monitoring solutions feel like an overkill. All of them are written in Go.

## freespace

Checks available disk space and sends a notification to a Telegram chat if it's less than 10%. I don't like setting up email notifications for it is more cumbersome, but it's easy to add a new API client you need in the `pkg/notifier` package.

Build the script with `make build-freespace` or a custom build command, the artifact will be located in the `./bin` directory. Upload it to a server and run periodically using cron, e.g.:

```cron
0 * * * * cd /where/you/uploaded/it && TELEGRAM_TOKEN=<your bot token> TELEGRAM_CHAT_ID=<...> ./freespace  2>&1 | /usr/bin/logger -t freespace
```
