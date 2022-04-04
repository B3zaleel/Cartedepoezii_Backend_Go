# Cartedepoezii_Backend_Go

A rewrite of the Cartedepoezii backend in Go.

## Requirements

### Applications

+ **PostgreSQL** 12+
+ **Go** 1.17+

### APIs

+ A Google API should be created with at least an email sending scope and this API server's root URL should be one of the redirect URIs. The `credentials.json` and `token.json` files should be stored in the root directory of this project.
+ An Imagekit.IO API should be created. The public key, private key, and URL endpoint should be stored in the `.env.local` file according to the requirements mentioned below.

### Environment Variables

The required environment variables should be stored in a file named `.env` and each line should have the format `Name: Value`. The table below lists the environment variables that will be used by this server:

| Name | Description |
|:-|:-|
| GIN_MODE | `production` if running for production else `debug`. |
| DB_URL | The URL of the PostgreSQL database to connect to. |
| APP_MAX_SIGNIN_TRIES | The maximum number of sign in attempts a user can make in succession. |
| IMG_CDN_PUB_KEY | Imagekit.io public key. |
| IMG_CDN_PRI_KEY | Imagekit.io private key. |
| IMG_CDN_URL_EPT | Imagekit.io url endpoint. |
| GOOGLE_MAIL_SENDER | The email address of the account responsible for sending emails to users. |
| WEB_CLIENT_DOMAIN | The domain name of the web client. |
| APP_SECRET_KEY | The secret key for this application. |

## Installation

+ Create the environment variables mentioned above.
+ Install `libpq-dev`.
+ Install the Go packages using `go get github.com/B3zaleel/Cartedepoezii_Backend_Go`.
+ Start the database service using `sudo service postgresql start`.
+ Open an interactive shell with the `postgres` user by running `sudo -s -u postgres`.
+ Run `psql -f data/DBSetup.sql` to initialize the database in the interactive shell and exit.
+ Open an interactive shell with the `cartedepoezii_dev` user by running `sudo -s -u cartedepoezii_dev`. Remember to add the `cartedepoezii_dev` user (if it doesn't exist) before running this command.
+ Run `psql -f src/db/DBInit.sql -d cartedepoezii_dev_db` to initialize the database entitites and exit.

## Usage

Run the server using `./run.bash`.

## Related Projects

+ [Cartedepoezii's FastAPI API server](https://github.com/B3zaleel/Cartedepoezii/tree/main/backend)
