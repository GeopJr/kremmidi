<h1 align="center">🧅 Kremmidi 🧅</h1>
<h4 align="center">Experimental GetTor backend with multiple frontends.</h4>
<p align="center">
  <br />
    <a href="./CODE_OF_CONDUCT.md"><img src="https://img.shields.io/badge/Contributor%20Covenant-v2.1-262329.svg?style=for-the-badge&labelColor=7f77bc" alt="COC Contributor Covenant v2.1" /></a>
    <a href="./LICENSE"><img src="https://img.shields.io/badge/LICENSE-BSD--2--Clause-262329.svg?style=for-the-badge&labelColor=7f77bc" alt="BSD-2-Clause" /></a>
</p>

<p align="center">

```mermaid
flowchart TB
    Kremmidi --> |Part & Binary\nmetadata| MongoDB
    Kremmidi --> |GET| Mirrors[TorProject, Calyx, EFF]
    Mirrors --> |Binaries| Kremmidi
    Kremmidi --> |Save & Split| FileSystem
    FileSystem --> |Path| Kremmidi
    FE --> |Parts| User
    FE --> |Find| MongoDB
    FE --> |Read| FileSystem
    FileSystem --> |Parts| FE
    MongoDB --> |Part & Binary\nmetadata| FE
    User --> |Parts| onion-rings
    onion-rings --> |Binary| User
    subgraph Backend
        Kremmidi
        Mirrors
        FileSystem
    end
    subgraph Database
        MongoDB
    end
    subgraph Frontend
        FE[Discord, Slack, Matrix...]
    end
```

</p>

# What is this?

One of [GetTor](https://gettor.torproject.org/)'s main roadblocks is the upload limits of today's platforms - it's very rare for platforms to allow uploading files that are >=10 MB in size. Kremmidi aims to solve that through a 2-step workaround.

**Kremmidi** is responsible for **downloading** the binaries (both desktop and Android) from multiple mirrors (TorProject, Calyx, EFF), **splitting** them into parts of 5MB each (configurable) and **inserting** their paths (and other relative metadata) to the database.

Then **frontends** consume the database and send all parts to the user - one at a time - "bypassing" their upload limits.

After the user downloads all the parts, they have to use **onion-rings** to compile them back into one.

**onion-rings** is a webpage that the user receives, downloads and runs in their browser where they **drag-n-drop** (or "open") all the parts downloaded previously. **onion-rings** then puts them in order, **"compiles"** them back into a single/complete binary and prompts the user to download it.

# Building

Dependencies: `go`

```
$ make
```

# Docker

```bash
$ docker compose -f docker-compose.yml up -d
```

# Frontends

You can find some example frontends with more info at [`./examples`](./examples/).

# Example

```bash
kremmidi -h
NAME:
   kremmidi - Experimental GetTor backend with multiple frontends.

USAGE:
   kremmidi [global options] [command]

VERSION:
   v0.2.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --database value, -d value, --db value               Database url.
   --drop-db                                            Whether to drop all collections on start. (default: false)
   --help, -h                                           show help (default: false)
   --limit value, -l value                              Amount of *MB* per part. (default: 5)
   --locales value                                      Locales to download seperated by a comma or 'ALL'. (default: en-us)
   --mirror value, -m value                             Use a mirror. Available: TOR_PROJECT, EFF, CALYX (default: TOR_PROJECT)
   --no-android                                         Disable grabbing Android builds. (default: false)
   --no-desktop                                         Disable grabbing Desktop builds. (default: false)
   --no-split                                           Disable splitting binaries into parts. (default: false)
   --output value, -o value, --out value, --dest value  Binary destination. (default: ./data)
   --overwrite                                          Whether to overwrite builds if they already exist. (default: false)
   --version, -v                                        print the version (default: false)
```

```bash
$ kremmidi -o data -m eff --db=mongodb://user:pass@host:27017
```

```bash
$ KREMMIDI_DB=mongodb://user:pass@host:27017 kremmidi  --no-desktop

2022/10/02 11:36:45 Started downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-aarch64-multi.apk
2022/10/02 11:37:42 Finished downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-aarch64-multi.apk
2022/10/02 11:37:42 Started downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-aarch64-multi.apk.asc
2022/10/02 11:37:43 Finished downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-aarch64-multi.apk.asc
2022/10/02 11:37:43 Saving to database...
2022/10/02 11:37:43 Finished saving to database.
2022/10/02 11:37:43 Started downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-armv7-multi.apk
2022/10/02 11:38:36 Finished downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-armv7-multi.apk
2022/10/02 11:38:36 Started downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-armv7-multi.apk.asc
2022/10/02 11:38:36 Finished downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-armv7-multi.apk.asc
2022/10/02 11:38:36 Saving to database...
2022/10/02 11:38:36 Finished saving to database.
2022/10/02 11:38:36 Started downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-x86_64-multi.apk
2022/10/02 11:39:37 Finished downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-x86_64-multi.apk
2022/10/02 11:39:37 Started downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-x86_64-multi.apk.asc
2022/10/02 11:39:37 Finished downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-x86_64-multi.apk.asc
2022/10/02 11:39:37 Saving to database...
2022/10/02 11:39:37 Finished saving to database.
2022/10/02 11:39:37 Started downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-x86-multi.apk
2022/10/02 11:40:38 Finished downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-x86-multi.apk
2022/10/02 11:40:38 Started downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-x86-multi.apk.asc
2022/10/02 11:40:38 Finished downloading https://dist.torproject.org/torbrowser/11.5.3/tor-browser-11.5.3-android-x86-multi.apk.asc
2022/10/02 11:40:38 Saving to database...
2022/10/02 11:40:38 Finished saving to database.
```
