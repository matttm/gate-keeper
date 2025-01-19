# gate-keeper

## Description

This a very simple program I wrote to automate a common task at work--opening and closing application gates to mimic an application cycle flow. I tried writing this generically, so that some user with a different database schema can find this helpful. Specifying this schema is done in a `config.json`.

To help you visualize what this program does. see the below timeline:
```
       aaaaa       bbbbb       ccccc       ddddd       eeeee     
------------------------------------------------------------------------
```
This represents an application cycle, as you can see, the cycle has periods a through e.

This program will configure the gates for a given year such that you will be virtually teleported.

Lets say, given the database of the above system, and as program inputs I provided: `year 2025`, `at gate c`, `place me before it`

TThese inputs, will update the start datetime and end datetime of the application gates so that you are placed on the timeline as shown below. (indicated by `me`)
```
       aaaaa       bbbbb  me   ccccc       ddddd       eeeee     
------------------------------------------------------------------------
```
In other words, the open and close dates of the gates are configured so that the current moment is before the opening of gate `c`.

## Getting started

To begin using this program, it rrequires two files for proper configuration--a `config.json` which is holds a mapping of the database schema, and an environment file, `.env`. In practice, the `.env` is not needed, we just need the contents to be in the environment of the one running the app.

These files look like:

config.json
```
{
	"Dbname": "DB_NAME",
	"TableName": "GATE_TABLE",
	"GateNameKey": "GATE_CODE",
	"GateYearKey": "GATE-YEAR",
	"StartKey": "GATE_STRT_DT",
	"EndKey": "GATE_END_ST"
}
```

.env
```
export DB_HOST="database.url.com"
export DB_USERNAME="USER"
export DB_PASSWORD="CocoNuttz"
export DB_PORT="3306" # Default for MySql is 3306
```

## Authors

-   Matt Maloney : matttm

## Contribute

If you want to contribute, just send me a message.
