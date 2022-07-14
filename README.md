# STOREHELPER
Project intended for doing an easy thing the hard way.
## Description
storehelper is a configurable command-line tool used for scripted management and cleanup of multiple file storages.
In short: it takes FS path and a set of instructions as an input and performs actions on given files. Originally it was intended to be run on schedule on Windows servers with a config file to perform routine tasks which I was too tired and lazy to do manually. (and it's also easily automated, so...)

Input is formatted like this (subject to change):
```
f n:10 w:1 k
```
which means to select in given path files created on mondays for last 10 days, keep them and delete the rest.
Or
```
f e:10 m:\\new-path\
```
which means to take files older than 10 days in given path and move them to their new destination.

Exact up to day manual is available in sample `conf` file.

## Why not just make that into a Python or Shell script?
I don't want to 🤷‍♀️
I also wanted it to be cross-platform and easily portable, so no.

## Features
### Implemented:
Filters: newest, oldest, day of week, day of month, age (days)

Operations: remove, keep, copy

### In works:
Filters: -

Operations: simple rename, email

### Some day:
Filters: by complex date mask

Operations: rename by complex mask, proper reports