# RaspController

Control your Raspberry Pi with a web interface.

## Technologies Used

- [RoseDB](https://github.com/rosedblabs/rosedb)
- [go-gpiocdev](https://github.com/warthog618/go-gpiocdev)
- [mdns](https://github.com/hashicorp/mdns)
- [Fiber](https://gofiber.io/)

## Observations

- This project is developed to run on a Raspberry Pi.
- RoseDB is used to store data for history, charts, and manage pins.
- Configuration should be specified in the `conf.yml` file.
- mDns is used for service discovery.
- Some functionalities require [vcgencmd](https://github.com/raspberrypi/utils).
- Tested on Arch Linux ARM but should work on other distributions.

## Configurations

```yaml
JWT_SECRET: "secretkey"
DB_DIR: "./db"
PORT: 8080
JWT_EXPIRES_IN: 3600
SHARE_DIR: "/home/rasp/public"
```

## Routes

All routes begin with /api

### Information

#### List all informations
```http
GET /api/info
```
#### List all process
```http
GET /api/info/ps
```

### GPIO

#### List all GPIOs

```http
GET /api/gpio
```

#### Get information about a GPIO

```http
GET /api/gpio/:id
```

#### Update a GPIO

```http
PATCH /api/gpio/:id
Content-Type: application/json

{
    "direction": "out",
    "value": 1
}
```