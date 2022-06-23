# Tong

<img src="tong.png" alt="tong.png" width="150"/>

Tong is a command-line tool facilitating development of gotong-based application.


## Installation

```sh
go install github.com/Tuingking/tong@latest
```

## Basic commands

Tong provides a variety of commands which can be helpful at various stages of development. The top level commands include:

```
    version     Prints the current tong version
    kafka       Start kafka zookeeper and broker
    sql         MySQL related utility command
    split-file  Split huge `csv` or `tsv` file into multiple files

    # tong sql
    migrate     Runs database migrations
    ddl         Generate MySQL DDL query from go struct
    find-field  Find field in database table        
```
