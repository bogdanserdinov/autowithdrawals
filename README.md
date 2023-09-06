# Auto withdrawal tool for Kaminari

Tool to automate withdraw on-chain BTC from Kaminari to given address. All code shown as example how to implement auto withdraw, feel free to change and use it.

## Usage:

To start CLI tool you need to fill `.env` file, where you should specify max account balance, in case when your account balance > expected auto withdrawals will be initiated to specified address.

Command to start service:

`go run main.go by-amount` - command that will check your account balance every hour, and if needed will initiate on-chain tx to your 'back-up' address.

`go run main.go by-date` - command that will initiate payment every given interval.
