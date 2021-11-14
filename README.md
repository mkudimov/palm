## _Palm is a fast token reader_

## Overview
CLI tool for generating and fast reading tokens.

## Features
Palm command _**generate**_ generates random tokens using lowercase letters _**a-z**_ and dumps them to file (one token per line).
Command _**read**_ reads and handles tokens from a file. User may want to dump unique tokens to a storage, so this tool allows him to write tokens to postgresql database. To set database connection, the user should use configuration file. 

#### Generate
This command has 3 flags:
* out (shorthand o) - required flag that describes output file path for generated tokens;
* number (shorthand n) - not required flag describing the number tokens needed to generate;
* length (shorthand l) - not required flag that describes the length of each token.

Example of usage
```
palm generate -n 3 -l 7 -o ~/home/user/desktop/output-sample.txt
```
Result (output-example.txt content)
```
ftyolck
ycmslxb
znuskor
```

#### Read
Has 2 flags:
* in (shorthand i) - required flag describing input file path with tokens;
* out (shorthand o) - required flag describing output file path for non-unique tokens and their frequencies.

Example of usage
```
palm read -i ~/home/user/desktop/output-sample.txt -o ~/home/user/desktop/non-unique_tokens.txt
```
Result (non-unique_tokens.txt content)
```
{"ftyolck":3, "ycmslxb":2}
```

## Database dumping
As it mentioned above palm allows to dump unique tokens to postgresql storage. To take advantage of this feature an user need to place his configuration file next to palm's execution file. The configuration file must be called as _**.palm.yaml**_ and looks like
```
postgresql: postgres://d4l:d4l@localhost:5432/d4l
```

Within the postgresql database the user needs to create a table using the following script
```
create table tokens (token varchar not null unique);
```

Within _**db.script**_ you can find commands for creating postgresql role and database.

## Design decision
Palm uses [cobra](https://github.com/spf13/cobra) package to work as a CLI tool. 

Package token_generator implements the functionality for generating tokens and dumping them to a file.
Package token_handler implements the functionality for reading and handling tokens from a file. There is struct _Handler_ which incapsulates an instance of interface _Reader_ from package token_handler/readers within the package. This approach has an advantage allows to use different implementations of the _Reader_. 
There are two implementations of _Reader_ interface in token_handler/readers package 
* DefaultReader - reads the whole file and then token by token handles them. 
* FastReader - an asynchronous reader that runs a number of goroutines equal to the number of logical processors on the computer and each goroutine reads a chunk from the file and then processes the read chunk.

There is also interface _Hook_ which can be used for dumping unique tokens to a storage. There is only one implementation of this interface - _PostgreSQLHook_.