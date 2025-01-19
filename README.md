# gate-keeper

## Description

This a very simple program I wrote to automate a common task at work--opening and closing application gates to mimic an application cycle flow. I tried writing this generically, so that some user with a different database schema can find this helpful. Specifying this schema is done in a `config.json`.

To help you visualize what this program does. see the below timeline:
```
    aaaaaaaaa       bbbbbbbb        ccccc     ddddd          eeeee
------------------------------------------------------------------------
```
This represents an application cycle, as you can see, the cycle has periods a through e.

This program will configure the gates for a given year such that you will be virtually teleported.

Lets say, given the database of the above system, and as program inputs I provided: `year 2025`, `at gate c`, `place me before it`

TThese inputs, will update the start datetime and end datetime of the application gates so that you are placed on the timeline as shown below. (indicated by `me`)
```

    aaaaaaaaa       bbbbbbbb   me   ccccc     ddddd          eeeee
------------------------------------------------------------------------
```
