# wowaddon
Manager for World of Warcraft addons

## Features

* Install addons from Curse and TukUI
* Uninstall addons
* Update addons
* List installed addons, showing which are "out of date"
* Show addon metadata
* Show which directories are used by which addons

## Installation

Download the file from the [github release page](https://github.com/wttw/wowaddon/releases/tag/v0.1.0),
unzip it and put it somewhere on your path.

Run `wowaddon environment`. If it can't find your World of Warcraft
installation in the normal places you can set the environment variable
WOWDIR to override it.

## Usage

To install [Fishing Buddy](https://mods.curse.com/addons/wow/fishingbuddy)
run `wowaddons install fishingbuddy`. The name needed is the one used for
the addon at the source you're fetching the addon from. For addons installed
from Curse, it's the name you see in the URL of the page, such as
`deadly-boss-mods` or `weakauras-2`.

To list the installed addons (that `wowaddon` is aware of - if you've
manually unzipped addons it doesn't know about them) run `wowaddons list`.

To update everything that needs to be updated run `wowaddon update`.

```
NAME:
   wowaddon - Install WoW addons

USAGE:
   wowaddon [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR(S):
   Steve Atkins <steve@blighty.com>

COMMANDS:
     install            Install addon `NAME`
     uninstall          Uninstall addon `NAME`
     reinstall          Reinstall all addons
     update             Update all addons
     checkupdate        List addons that can be updated
     folders, list, ls  List addons and their folders
     blame              Show which addon created a folder
     environment, env   Show environment
     info               Show information about installed addons
     fullinfo           Show toc metadata about installed addons
     dlurl              Find a download URL

GLOBAL OPTIONS:
   --wowdir value, --dir value, -d value  WoW base directory [$WOWDIR]
   --config value, -c value               Use an alternate configuration file [$WOW_ADDON_CONFIG]
   --cache value                          Use an alternate cache directory [$WOW_ADDON_CACHE]
   --help, -h                             show help
   --version, -v                          print the version
```

## Compilation

Install pre-requisites and build
```
go get github.com/urfave/cli
go get github.com/fatih/color
go get github.com/kardianos/osext
go build
```

This will create a single binary, `wowaddon` or `wowaddon.exe` that can
be copied to any machine and run with no other prerequisites.

## Inspiration
This tool was inspired by [wow-cli](https://github.com/zekesonxx/wow-cli),
a similar tool written in Javascript/node.

The user interface is very similar, and the configuration file is nearly
identical - if you're using `wow-cli` you can copy your `.addons.json`
to `addons.json`, remove any addons sourced from anywhere other than
Curse and it should work with `wowaddon`.

## Issues
Coloured text works on Windows 10, OS X and (probably) Linux. It doesn't
work in a vanilla Windows 7 command prompt. It's still perfectly usable,
but not as pretty.

An obvious need is to be able to search for addons, but that data isn't
made easily available anywhere.
## License
[Two-clause BSD](LICENSE)
