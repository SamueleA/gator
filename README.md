# gator

Gator is an open-source RSS reader built in go and running on the command line.

# Installing

Run `go install` to install the program

Postgres and go needs to be installed on the machine prior to installation.

# Running

The program can be run with `go run . COMMAND`

# Commands

Before using the program an account must be created with the register command.

`gator register my_user_name`

Then add feeds with the command `addFeed`

`gator addFeed name url`

Finally, run the command `agg` to have new posts regularly fetched each time interval.

`gator agg 5m`

## Generate Go DB queries

To generate the Go DB queries run the following command:
`sqlc generate`
