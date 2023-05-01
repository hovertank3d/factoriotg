# factorio telegram bridge

## building

as simple as 

```bash
go mod tidy
go build -o ftg
```

## installing

### factorio server

if you want messages without server prefix , you can install plugin with /say command.

## running

it must be piped from factorio with set envvars. so, preferred way looks like:
```bash
export RCON_PORT=27015
export RCON_PASS="amongus super secret password"
export TELEGRAM_TOKEN="blah blah blah your tg token must be here"
export TELEGRAM_GROUP="chatid for group"

./bin/x64/factorio --rcon-port=$RCON_PORT --rcon-password=$RCON_PASS --start-server=save.zip | ./ftg
```

you also can use other bridges, for example [irc brigde](https://github.com/mickael9/factoirc), just piping things :)