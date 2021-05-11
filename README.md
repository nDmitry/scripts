## freespace

Checks available disk space and sends a notification to a Telegram chat if it's less than 10%.

Build the script with `make build-freespace` or a custom build command, the artifact will be located in the `./bin` directory. Upload it to a server and run periodically using cron, e.g.:

```cron
0 * * * * cd /where/you/uploaded/it && TELEGRAM_TOKEN=<your bot token> TELEGRAM_CHAT_ID=<...> ./freespace  2>&1 | /usr/bin/logger -t freespace
```