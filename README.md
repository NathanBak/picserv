# picserv

picserv is a picture server written in Go. It connects to a Dropbox folder containing images and serves them up one at a time via a simple HTTP API.

## How It Works

On startup, picserv fetches the list of all files in the configured Dropbox folder, shuffles them, and then serves them sequentially. When all images have been served, the file list is refreshed from Dropbox and reshuffled before the cycle repeats. If the refresh fails or returns no files, the existing list is reused.

## Prerequisites

- A [Dropbox app](https://www.dropbox.com/developers/apps) with `files.content.read` and `files.metadata.read` permissions
- A refresh token for the app (see [Authentication](#authentication) below)

## Configuration

Configuration is via environment variables. A `.env` file in the working directory is loaded automatically.

| Variable | Required | Default | Description |
|---|---|---|---|
| `DROPBOX_APP_KEY` | yes | | Dropbox app key |
| `DROPBOX_APP_SECRET` | yes | | Dropbox app secret |
| `DROPBOX_USER_REFRESH_TOKEN` | yes | | OAuth2 refresh token |
| `DROPBOX_SEED` | no | `0` | Random seed for shuffling. `0` = use current time, negative = no shuffle |
| `PORT` | no | `8080` | HTTP port to listen on |
| `READ_TIMEOUT` | no | `3s` | HTTP read timeout |
| `WRITE_TIMEOUT` | no | `3s` | HTTP write timeout |
| `PICTURE_LIFE` | no | `9s` | How long each image is cached before fetching the next |

## Authentication

To obtain a refresh token:

1. Call `dbpics.MakeUserAuthenticationURL(appKey)` to get an authorization URL
2. Visit the URL, authorize the app, and copy the code shown
3. Call `dbpics.GetOfflineAccessToken(appKey, appSecret, code)` to exchange the code for a refresh token
4. Store the `RefreshToken` field value as `DROPBOX_USER_REFRESH_TOKEN`

## Running

```sh
go run .
```

## Docker

### Build

```sh
docker build -t picserv .
```

### Run

The image defaults to port `8080` internally. Use `--env-file` to pass credentials and `-e PORT=8080` to ensure the port setting is not overridden by the env file. Adjust the left side of `-p` to map to a different host port.

```sh
docker run --env-file .env -e PORT=8080 -p 8083:8080 picserv
```

### Test

With the container running, open a browser and navigate to http://localhost:8083/next to verify the service is working.

## API

### `GET /next`

Returns the next image as `application/octet-stream`.

## License

MIT — see [LICENSE](LICENSE)
