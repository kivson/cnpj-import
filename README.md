# CNPJ Importer

<!--toc:start-->
- [CNPJ Importer](#cnpj-importer)
  - [Usage](#usage)
  - [Performance](#performance)
<!--toc:end-->

Tool to import all brazilian company tax ids (CNPJs) open data into a sqlite database

The executable has to commands:

- `download` to download all CNPJs zip files from [receita federal](https://dadosabertos.rfb.gov.br/CNPJ/) to a folder.
- `import` to import all the zips from a folder to a SQLite database.

## Usage

To download all zips:

```bash
./cnpj-import download -u "https://dadosabertos.rfb.gov.br/CNPJ/" -f "./zips/"
```

To import to sqlite:

```bash
./cnpj-import import -p "./zips/" -d "database.db" 


```

The database dsn string `-d` is the same used by the GORM library to connect to SQLite.

## Performance

The application is able to download at Receita's maximum allowed speed (~11 MB/s)

The database insertion speed will depends on the disk, but you can the dsn string `-d` with the modificators `database.db?_journal_mode=MEMORY&_sync=OFF`
to improve the insertion performance. It will open the database with the journal in MEMORY and disble syncs, be aware that this can cause data loss if
you have an abrupt interruption during the database use. But, considering this tool use case this should not be a problem.
