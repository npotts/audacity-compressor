# audacity-compressor
A simple Audacity Recursive Project Compressor

# What it does
This is tool I threw together in a couple hours to recurse over a directory, and 
pack all the discovered [Audacity](https://www.audacityteam.org) projects into a tar.gz file for archival.

No optimzation was attempted.

# Usage:
There really isnt much here:

```sh
λ ./audacity-compressor --help
usage: audacity-compresser [<flags>] <root>

A stupid tool to locate and compress audacity project data

Flags:
      --help       Show context-sensitive help (also try --help-long and
                   --help-man).
  -o, --output=""  Output dir. Empty means place in the same directory as the
                   project

Args:
  <root>  Root Path to parse

```