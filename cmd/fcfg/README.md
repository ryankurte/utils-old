# FCFG (File ConFiGurator)

FCFG is a very small utility that parses template files with replacements from the command line or the environment.

Yes, you can achieve this with a pile of bash and/or sed. No, it's not a good time for anyone.

## Usage

`./fcfg -i=[INPUT] -o=[OUTPUT] [-v=KEY:VALUE] [-k=KEY] [--verbose]`

- Use go/mustache standard template fields (`{{.NAME}}`) in your config files
- Specify the input template with `-i` or `--input`
- Specify the output file with `-o` or `--output`
- Specify keys to load from the environment (`-k=key`) as required
- Specify key:value pairs (`-v=key:value`) from the command line
- (Optionally) add `--verbose` to see what's getting loaded

For example:  

`./fcfg -i=cmd/fcfg/test.tmpl.yml -o=test.yml -v=KEY_ONE:lol -v=KEY_TWO:boop -k=PWD`

