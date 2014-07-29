ttygif
======

Generate an animated GIF by playing ttyrecord (by `ttyrec`) and taking screenshots of your terminal.

![](https://cloud.githubusercontent.com/assets/80381/3628176/6569016e-0e91-11e4-9b0d-6bbfd46a6d32.gif)


About ttyrec
------------

see [http://0xcc.net/ttyrec/](http://0xcc.net/ttyrec/).


Installation
------------

    go get github.com/sugyan/ttygif

or download binaries from [Releases](https://github.com/sugyan/ttygif/releases).


Usage
-----

    ttygif -in <input file> -out <output file> -s <speed> -noloop <bool>

* `in`: ttyrec file (default: `"ttyrecord"`)
* `out`: output animated GIF file name (default: `"tty.gif"`)
* `s`: play speed (default: `1.0`)
* `noloop`: play only once (default: `false`)
