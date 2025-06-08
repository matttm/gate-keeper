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

<img width="519" alt="Screenshot 2025-06-08 at 2 08 19 PM" src="https://github.com/user-attachments/assets/9b1a2f84-7216-49d6-8683-57eb815d72c7" />


Once you select a year, you can now see the current gate statuses and get real time updates every second

<img width="509" alt="Screenshot 2025-06-08 at 2 14 05 PM" src="https://github.com/user-attachments/assets/cdcb735b-8ea8-4abf-bdcf-401d3399dad8" />
<img width="504" alt="Screenshot 2025-06-08 at 2 10 06 PM" src="https://github.com/user-attachments/assets/322c8f5c-90fa-4a20-868a-34543570fb74" />


There is now also a gate health check, which basically causes all the gates to turn yellow, indicating that they are out of order compared to the sort order in the database.

<img width="522" alt="Screenshot 2025-06-08 at 2 08 52 PM" src="https://github.com/user-attachments/assets/49c5c40d-c8be-4901-98a3-ccc27aa4f8d7" />



.
## Changelog
- 06/08/2025 - gate status table with a gate health check

## Authors

-   Matt Maloney : matttm

## Contribute

If you want to contribute, just send me a message.
