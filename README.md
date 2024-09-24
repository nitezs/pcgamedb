# pcgamedb

pcgamedb is a powerful command-line tool designed to scrape and manage repack game data from various online sources. With support for multiple data sources and the ability to provide a RESTful API.

## Features

- **Data Sources**:

  - Fitgirl
  - KaOSKrew
  - DODI
  - ~~FreeGOG~~
  - GOGGames
  - OnlineFix
  - Xatab
  - ~~ARMGDDN~~
  - SteamRIP
  - Chovka

- **Database**:

  - Stores game data in MongoDB
  - Supports Redis for caching to improve performance

- **RESTful API**:
  - Provides an API for external access to the game data

## Usage

run `go run . help`.

## Configuration

Edit the `config.json` file to set up your environment or set system environment variables.

Read `/config/config.go` for more details.

## Api Doc

Read `http://127.0.0.1:<port>/swagger/index.html` for more details.
