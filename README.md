
![Dark theme preview](http://invidget.switchblade.xyz/967058504403808356)

# Kelbi - MHFZ Server

Kelbi is a fork of [Erupe](#whats-erupe) to customize our server (Invasor de Fronteiras).

### Table of Contents

- [What's Erupe?](#whats-erupe)
- [Contributing](#contributing)
- [Local development](#local-development)
- [License](#license)

## What's Erupe?

Europe is a Monster Hunter Frontier Z Server developed by several people!

#### Here are the main repositories maintained by the community that this code is based on

- [ZeruLight/Erupe](https://github.com/ZeruLight/Erupe) - Main repository maintained by the CE (Community edition)
- [xl3lackout/Erupe](https://github.com/xl3lackout/Erupe) - Old repository of CE.

## Local Development

TODO

## Contributing

If your contribution is an advance for Erupe, prefer to send it to the [main community repository](#here-are-the-main-repositories-maintained-by-the-community-that-this-code-is-based-on), but otherwise, feel the advantage for any contributions! If you have any questions, don't hesitate to join [our community]() server and ask as many questions as you like.

## License

Licensed under [MIT](/LICENSE).

<!-- 
## WARNING

This project is in its infancy and has no reliable active developer, no documentation, and no support.

# General info

Currently allows a JP MHF client (with GameGuard removed) to:

* Login and register an account (registration is automatic if account doesn't exist)
* Create a character
* Get ingame to the main city
* See other players walk around
* Do quests
* Use chat*

# Installation

## Server

1. Start database with docker

```shell
docker-compose up -d
```

2. Run migrations

```shell
cd database
npm i
npm run db:push
```

5. Edit the config.json

    Namely:
    * Update the database username and password
    * Update the `host_ip` and `ip` fields (there are multiple) to your external IP if you are hosting for multiple clients.

6. Place quest/scenario binaries.

    The quest and scenario binary files should be placed in `bin/quests/` and `bin/scenarios` respectively.

## Launcher

Erupe ships with a rudimentary custom launcher, so you don't need to obtain the original TW/JP files to simply get ingame. However, it does still support using the original files if you choose to. To set this up, place a copy of the original launcher html/js/css in `./www/tw/`, and `/www/jp/` for the TW and JP files respectively.

Then, modify the the `/launcher/js/launcher.js` file as such:

* Find the call to `startUpdateProcess();` in a case statement and replace it with `finishUpdateProcess();`. (This disables the file check and updating)
* (JP ONLY): replace all uses of "https://" with "http://" in the file.

Finally, edit the config.json and set `UseOriginalLauncherFiles` to `true` under the launcher settings.

# Usage

### Note: If you are switching to/from the custom launcher html, you will have to clear your IE cache @ `C:\Users\<user>\AppData\Local\Microsoft\Windows\INetCache`

## Server

```
cd Erupe
go run .
```

## Client

Add to hosts:

```
127.0.0.1 mhfg.capcom.com.tw
127.0.0.1 mhf-n.capcom.com.tw
127.0.0.1 cog-members.mhf-z.jp
127.0.0.1 www.capcom-onlinegames.jp
127.0.0.1 srv-mhf.capcom-networks.jp
```

Run mhf.exe normally (with locale emulator or appropriate timezone). -->
