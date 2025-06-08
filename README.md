# gate-keeper

<!-- [![Go Coverage](https://github.com/matttm/gate-keeper/wiki/coverage.svg)](https://raw.githack.com/wiki/matttm/gate-keeper/coverage.html) -->

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

These inputs, will update the start datetime and end datetime of the application gates so that you are placed on the timeline as shown below. (indicated by `me`)
```
       aaaaa       bbbbb  me   ccccc       ddddd       eeeee     
------------------------------------------------------------------------
```
In other words, the open and close dates of the gates are configured so that the current moment is before the opening of gate `c`.

For the third input, other options are `after` and `inside`, which wouls appear as:

inside:
```
                                me
       aaaaa       bbbbb       ccccc       ddddd       eeeee     
------------------------------------------------------------------------
```
after
```
       aaaaa       bbbbb       ccccc   me  ddddd       eeeee     
------------------------------------------------------------------------
```

---
NOTE: To be clear and specific, this program will update all gates that are part of an application cycle, such that the selected gate is today's date, with a relative offset, which is set based on the program's third parameter.
--

## Getting started

To begin using this program, it rrequires one file for configuration.

This file look like:

config.json
```
{
	"Credentials": {
		"User": "USER",
		"Pass": "xxxx",
		"Host": "url",
		"Port": "3306"
	},
	"GateConfig": {
		"Dbname": "dbname",
		"TableName": "table",
		"GateNameKey": "code",
		"GateYearKey": "year",
		"GateOrderKey": "order",
		"GateIsApplicableFlag": "active",
		"StartKey": "start",
		"EndKey": "end"
	}
}

```

Once you have this file in your working directory, let's download the dependencies and build this thing.
```
go mod download

go build
```
and run the binary:
```
./gate-keeper
```
Then just choose the selects and press 'Set Gates'.

<img width="568" alt="Screenshot 2025-01-19 at 1 45 29 PM" src="https://github.com/user-attachments/assets/83ecfe93-1c49-437a-8bec-5c73b3efbb67" />

Once you select a year, you can now see the current gate statuses and get real time updates every second
<img width="751" alt="Screenshot 2025-06-08 at 9 09 20 AM" src="https://github.com/user-attachments/assets/632aac03-3187-4ae0-97c0-8ceab281d7d2" />
<img width="838" alt="Screenshot 2025-06-08 at 9 09 01 AM" src="https://github.com/user-attachments/assets/2136115f-d56f-4a4f-8809-9072a4144485" />

There is now also a gate health check, which basically causes all the gates to turn yellow, indicating that they are out of order compared to the sword order in the database.
<img width="775" alt="Screenshot 2025-06-08 at 11 17 31 AM" src="https://github.com/user-attachments/assets/d4f8e195-69d4-428d-b887-8b56c3920d19" />


.
## Changelog
- 06/08/2025 - gate status table with a gate health check
- 
## Authors

-   Matt Maloney : matttm

## Contribute

If you want to contribute, just send me a message.
