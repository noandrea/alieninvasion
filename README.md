Alien Invasion
==============
Mad​ aliens​ are​ about​ to​ invade​ the​ earth​ and​ you​ are​ tasked​ with​ simulating​ the
invasion.


You​ are​ given​ a map​ containing​ the​ names​ of​ cities​ in​ the​ non-existent​ world​ of
X.​ The​ map​ is​ in​ a file,​ with​ one​ city​ per​ line.​ The​ city​ name​ is​ first,
followed​ by​ 1-4​ directions​ (north,​ south,​ east,​ or​ west).​ Each​ one​ represents​ a
road​ to​ another​ city​ that​ lies​ in​ that​ direction.
For​ example:

```
Foo​ north=Bar​ west=Baz​ south=Qu-ux
Bar​ south=Foo​ west=Bee
```

The​ city​ and​ each​ of​ the​ pairs​ are​ separated​ by​ a single​ space,​ and​ the
directions​ are​ separated​ from​ their​ respective​ cities​ with​ an​ equals​ (=)​ sign.
You​ should​ create​ N aliens,​ where​ N is​ specified​ as​ a command-line​ argument.
These​ aliens​ start​ out​ at​ random​ places​ on​ the​ map,​ and​ wander​ around​ randomly,
following​ links.​ Each​ iteration,​ the​ aliens​ can​ travel​ in​ any​ of​ the​ directions
leading​ out​ of​ a city.​ In​ our​ example​ above,​ an​ alien​ that​ starts​ at​ Foo​ can​ go
north​ to​ Bar,​ west​ to​ Baz,​ or​ south​ to​ Qu-ux.

When​ two​ aliens​ end​ up​ in​ the​ same​ place,​ they​ fight,​ and​ in​ the​ process​ kill
each​ other​ and​ destroy​ the​ city.​ When​ a city​ is​ destroyed,​ it​ is​ removed​ from
the​ map,​ and​ so​ are​ any​ roads​ that​ lead​ into​ or​ out​ of​ it.

In​ our​ example​ above,​ if​ Bar​ were​ destroyed​ the​ map​ would​ now​ be​ something
like:

```
Foo​ west=Baz​ south=Qu-ux
```

Once​ a city​ is​ destroyed,​ aliens​ can​ no​ longer​ travel​ to​ or​ through​ it.​ This
may​ lead​ to​ aliens​ getting​ "trapped".

You​ should​ create​ a program​ that​ reads​ in​ the​ world​ map,​ creates​ N aliens,​ and
unleashes​ them.​ The​ program​ should​ run​ until​ all​ the​ aliens​ have​ been
destroyed,​ or​ each​ alien​ has​ moved​ at​ least​ 10,000​ times.​ When​ two​ aliens
fight,​ print​ out​ a message​ like:Bar​ has​ been​ destroyed​ by​ alien​ 10​ and​ alien​ 34!
(If​ you​ want​ to​ give​ them​ names,​ you​ may,​ but​ it​ is​ not​ required.)​ Once​ the
program​ has​ finished,​ it​ should​ print​ out​ whatever​ is​ left​ of​ the​ world​ in​ the
same​ format​ as​ the​ input​ file.

Feel​ free​ to​ make​ assumptions​ (for​ example,​ that​ the​ city​ names​ will​ never
contain​ numeric​ characters),​ but​ please​ add​ comments​ or​ assertions​ describing
the​ assumptions​ you​ are​ making.

## Preview

**Normal Mode**

[![asciicast](https://asciinema.org/a/349792.png)](https://asciinema.org/a/349792)


**Debug mode**

[![asciicast](https://asciinema.org/a/349793.svg)](https://asciinema.org/a/349793)

## Usage

For an example input file check [this one](https://github.com/noandrea/alieninvasion/blob/master/land/testdata/numpad.txt).

Check the help from the command line for how to use it, here an example

```
dist |> ./alieninvasion 
Simulate an alien invasion

Usage:
  alieninvasion [command]

Available Commands:
  help        Help about any command
  run         A brief description of your command
  version     print the program version

Flags:
      --debug     Enable debug output
  -h, --help      help for alieninvasion

Use "alieninvasion [command] --help" for more information about a command.

```

## Development

to build the project, from a terminal, run:

```
git clone https://github.com/noandrea/alieninvasion.git
cd alieninvasion
```

Run the test using `make test`

Run the build using `make build`, the binary files are stored in the `dist` folder buy default 
