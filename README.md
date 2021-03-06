# Countdown

<p align="center"><img src="https://user-images.githubusercontent.com/141232/54696023-9ed03e00-4b5d-11e9-9c7b-d6f67691e70c.gif" width="450" alt="Screen shot"></p>

## Usage

Specify duration in go format `1h2m3s`.

```bash
countdown 25s
```

Add command with `&&` to run after countdown.

```bash
countdown 1m30s && tput bel
```

## Key binding

- `p` or `P`: To pause the countdown.
- `c` or  `C`: To resume the countdown.
- `Space`: To pause/resume the countdown.
- `Esc` or `Ctrl+C`: To stop the countdown without running next command.
- `b` or `B`: To turn on/off the bell.

## Install

```bash
go get github.com/huyiwen/countdown
```

## License

MIT
