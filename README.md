# Mori
> 🌲 Automatically put osu! related archives in their places.

Mori is a simple file/directory watching daemon, which moves recently downloaded
osu! archives (osz/osk, beatmaps/skins) into their relevant directories.
For people that are lazy and don't want to drag and drop ([or can't](#Rationale)).

## Rationale
I have recently got into playing osu! again (I'm garbage), but im on Linux. Though
with a handy script, I got it setup easily and performs better than lazer on my bad
laptop. Thing is, drag and drop did not work so I could not import skins/beatmaps
in that way. This is why Mori was born. It automatically puts recently created/
downloaded osu! beatmaps into osu!'s song directory, and skins in the skins
directory. Then to import, just press F5 in beatmaps list. Or when you go back to it,
it'll import.

# Installation
Getting Mori is pretty easy. All you need is [Go](https://go.dev) downloaded (and git, yknow).

Type these commands in a terminal:  
```
git clone https://github.com/TorchedSammy/Mori
cd Mori
go get -d
make build
sudo make install
```  

# Usage
> ⚠️ Mori has only been tested on Linux (and is really the only intended platform),
beware!
> It also may not move already left beatmaps, read on how
[Mori watches for changes](#Detection)

Simply run `mori &`, which will spawn Mori as a background job in your shell.  
It'll log any files that were moved.

## Configuration
Mori uses a JSON config, located in `~/.config/mori/mori.json`.
It does not create this default config, but as a reference it looks like:  
```json
{
	"osuDir": "~/.local/share/osu-wine/OSU",
	"sourceDir": "~/Downloads"
}
```
`sourceDir` is where Mori will watch, and copy from.  
`osuDir` is the osu! data directory, where all the files are.

## Detection
Mori detects "new" osu! archives by watching the
configured directory to see if any files have been chmod'd, which is usually what
happens when a while is created (and also when you, well, chmod a file).
This means that beatmaps/skins that were already in the configured directory or moved in
won't automatically be moved to their direcory. <!-- Fear not, as Mori checks
every 5mins if there are files in the current directory that have the `osz` extension!
If they do, it'll move them. -->

# License
Mori is licensed under the BSD 3-Clause license.  
[Read here](LICENSE) for more info.

