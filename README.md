# Tensho
> 🐦 Directory watching daemon which copies downloaded beatmaps into the osu! songs directory

Tensho is a simple file/directory watching daemon, which moves recently downloaded
beatmaps into osu!'s songs directory. For people that are lazy and don't want to
drag and drop ([or can't](#Rationale)).

## Rationale
I have recently got into playing osu! again (I'm garbage), but im on Linux. Though
with a handy script, I got it setup easily and performs better than lazer on my bad
laptop. Thing is, drag and drop did not work so I could not import skins/beatmaps
in that way. This is why Tensho was born. It automatically puts recently created/
downloaded osu! beatmaps into osu!'s song directory. Then to import, just press F5
in beatmaps list. Or when you go back to it, it'll import.

# Installation
Getting Tensho is pretty easy. All you need is [Go](https://go.dev) downloaded (and git, yknow).

Type these commands in a terminal:  
```
git clone https://github.com/TorchedSammy/Tensho
cd Tensho
go get -d
go build
```  
And Tensho will be compiled, at which you can copy/move to any bin directory.

# Usage
> ⚠️ Tensho has only been tested on Linux (and is really the only intended platform),
beware!
> It also may not move already left beatmaps, read on how
[Tensho watches for changes](#Detection)

Simply run `tensho &`, which will spawn Tensho as a background job in your shell.  
It'll log any files that were moved.

# Detection
Tensho detects "new" beatmaps (files with an `osz` extension) by watching the
configured directory to see if any files have been chmod'd, which is usually what
happens when a while is created (and also when you, well, chmod a file).
This means that beatmaps that were already in the configured directory or moved in
won't automatically be moved to the songs directory. Fear not, as Tensho checks
every 5mins if there are files in the current directory that have the `osz` extension!
If they do, it'll move them.

# License
Tensho is licensed under the BSD 3-Clause license.  
[Read here](LICENSE) for more info.

