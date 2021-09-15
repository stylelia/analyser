# analyser

> A GitHub App built with [Probot](https://github.com/probot/probot) that Check latest master commit sha &amp; cookbook version against cache

## Setup

```sh
# Install dependencies
npm install

# Run the bot
npm start
```

## Docker

```sh
# 1. Build container
docker build -t analyser .

# 2. Start container
docker run -e APP_ID=<app-id> -e PRIVATE_KEY=<pem-value> analyser
```

## Contributing

If you have suggestions for how analyser could be improved, or want to report a bug, open an issue! We'd love all and any contributions.

For more, check out the [Contributing Guide](CONTRIBUTING.md).

## License

[ISC](LICENSE) © 2021 Artur Kondas <devilyoushy@gmail.com>
