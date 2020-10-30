# Helix Github CLI

## TL;DR

[Download](https://github.com/helix-collective/hxgh/releases/latest) or locally build `hxgh` (`go get -v github.com/helix-collective/hxgh`)
```
  export GITHUB_TOKEN=xxx
hxgh events csv --timezone-hours-offset 11 username > event.csv
```

## Details

<!--tmpl,code=bash:./hxgh -h -->
``` bash 

  Usage: hxgh [options] <command>

  Options:
  --gh-token, -g   env GITHUB_TOKEN
  --version, -v    display version
  --help, -h       display help

  Completion options:
  --install, -i    install bash-completion
  --uninstall, -u  uninstall bash-completion

  Commands:
  · events

  Version:
    dev

```
<!--/tmpl-->

<!--tmpl,code=bash:./hxgh events -h -->
``` bash 

  Usage: hxgh events [options] <command>

  Options:
  --help, -h  display help

  Commands:
  · csv

```
<!--/tmpl-->

<!--tmpl,code=bash:./hxgh events csv -h -->
``` bash 

  Usage: hxgh events csv [options] <username>

  Options:
  --timezone-hours-offset, -t
  --help, -h                   display help

```
<!--/tmpl-->

## Update the Readme

```
go get -v github.com/jpillora/md-tmpl
md-tmpl -w Readme.md
```