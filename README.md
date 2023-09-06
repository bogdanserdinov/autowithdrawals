# Auto withdrawal tool for Kaminari

Tool to automate withdraw on-chain BTC from Kaminari to given address. All code shown as example how to implement auto withdraw, feel free to change and use it.

## Usage:

To start CLI tool firstly you need to fill `.env` file.

### Command to start service:

`go run main.go by-amount` - command that will check account balance every hour, and if needed will initiate on-chain tx to your 'back-up' address.

`go run main.go by-date` - command that will initiate payment every given interval.
