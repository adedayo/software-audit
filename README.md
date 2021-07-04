# Software Audit Service and Tool

Software Audit is a simple service for checking what a piece of software is based on its sha1 checksum. An accompanying tool can scan the filesystem for executable files, checking the shasum against the service, to determine what software it is.

The service uses the NIST [National Software Reference Library (NSRL)](https://www.nist.gov/itl/ssd/software-quality-group/national-software-reference-library-nsrl) as the source of the checksums. On first run of the server, the full NSRL ISO is automatically downloaded by default and will not require intervention.

A command line tool is also provided to search the file system for executables to check against the server.

To run the service, invoke (_Be patient on first run, it will take a long while to download and extract the NSRL data_)

```bash
audit serve
```

Once the server finishes loading the data, to search for applications and other executable binaries on your file system issue

```bash
audit apps [directories-to-search]
```

The server need not be on the same machine. Check the command-line help to set use a remote service

```bash
Usage:
  audit [command]

Available Commands:
  apps        Audit applications.
  help        Help about any command
  serve       Run software audit API service

Flags:
      --config string   config file (default is $HOME/.softaudit.yaml)
  -h, --help            help for audit
  -t, --toggle          Help message for toggle

Use "audit [command] --help" for more information about a command.
```
