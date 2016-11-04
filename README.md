# wowaddon

Manager for World of Warcraft addons

## Features

* Install addons from Curse and TukUI
* Uninstall addons
* Update addons
* Search for new addons
* List installed addons, showing which are "out of date"
* Show addon metadata
* Show which directories are used by which addons
* No complex setup, it discovers your current set of addons

## Installation

Download the file from the [github release page](https://github.com/wttw/wowaddon/releases/latest),
for your operating system, unzip it and put it somewhere on your path. (If
you're on Windows you can open a command prompt, cd to the directory where
you unzipped it and run it from there.)

Run `wowaddon environment`. If it can't find your World of Warcraft
installation in the normal places you can set the environment variable
WOWDIR to override it.

## Usage

To install [Fishing Buddy](https://mods.curse.com/addons/wow/fishingbuddy)
run `wowaddon install fishingbuddy`. The name needed is the one used for
the addon at the source you're fetching the addon from or returned by
`wowaddon search`. For addons installed from Curse, it's the name you see
in the URL of the page, such as `deadly-boss-mods` or `weakauras-2`.

To list the installed addons run `wowaddon list`.

To update everything that needs to be updated run `wowaddon update`.

To sync with new addons if they've been installed manually, or to create
your first configuration file run `wowaddon bootstrap`.

To search for new addons for your druid to install run `wowaddon search druid`.

```
NAME:
   wowaddon - Install WoW addons

USAGE:
   wowaddon [global options] command [command options] [arguments...]
   
VERSION:
   0.7.0
   
AUTHOR(S):
   Steve Atkins <steve@blighty.com> 
   
COMMANDS:
     install, i         Install addon `NAME`
     update, u          Update all addons
     search, s          Search for new addons
     uninstall          Uninstall addon `NAME`
     reinstall          Reinstall all addons
     checkupdate        List addons that can be updated
     folders, list, ls  List addons and their folders
     blame              Show which addon created a folder
     environment, env   Show environment
     info               Show information about installed addons
     fullinfo           Show toc metadata about installed addons
     bootstrap          Create a configuration file from existing addons
     lock               Lock an addon so it won't be automatically updated
     unlock             Unlock a locked addon, so it will be automatically updated
     help, h            Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --wowdir value, --dir value, -d value  WoW base directory [$WOWDIR]
   --config value, -c value               Use an alternate configuration file [$WOW_ADDON_CONFIG]
   --cache value                          Use an alternate cache directory [$WOW_ADDON_CACHE]
   --useragent value                      Use this useragent for http requests [$WOW_ADDON_USERAGENT]
   --color, --colour                      Use coloured output [$WOW_ADDON_COLOUR]
   --help, -h                             show help
   --version, -v                          print the version
   

```

## Configuration files

All configuration is stored in your World of Warcraft base directory.

The addons installed and managed by `wowaddon` are stored in the
`addons.json` file, with a backup copy stored in `addons.json.bak`

The catalog of known addons (used for bootstrapping and searching) is
stored in `addoncatalog.json.zip` and `addoncatalog.json`. This will
be fetched and updated automatically.

## Compilation

Install pre-requisites and build
```
go get github.com/wttw/wowaddon
go install
```

This will create a single binary, `wowaddon` or `wowaddon.exe` that can
be copied to any machine and run with no other prerequisites.

## Inspiration

This tool was inspired by [wow-cli](https://github.com/zekesonxx/wow-cli),
a similar tool written in Javascript/node.

## Issues

While it's tested on OS X, it isn't tested with a real WoW installation (I
don't have the disk space for that).

## License

[Two-clause BSD](LICENSE)
