# RaspController

Control your Raspberry Pi with a user-friendly web interface.

## Overview

RaspController lets you manage and monitor your Raspberry Pi remotely. Key features include:

* **File Sharing:** Easily share files from a designated public folder.
* **Hardware and Software Monitoring:** Track RAM, CPU, disk usage, and running processes.
* **GPIO Control:** Configure and manipulate GPIO pins directly through the interface.
* **Service Discovery:** Automatic device detection on your network using mDNS.

## Technologies Used

* **RoseDB** ([https://github.com/rosedblabs/rosedb](https://github.com/rosedblabs/rosedb)): A fast and efficient embedded database for storing historical data, charts, and pin management.
* **go-gpiocdev** ([https://github.com/warthog618/go-gpiocdev](https://github.com/warthog618/go-gpiocdev)): Provides Go bindings to interact with Raspberry Pi GPIO pins.
* **mdns** ([https://github.com/hashicorp/mdns](https://github.com/hashicorp/mdns)): Implements multicast DNS (mDNS) for effortless service discovery.
* **Fiber** ([https://gofiber.io/](https://gofiber.io/)): A high-performance web framework for Go, powering the web interface.

## Requirements

* **Raspberry Pi:** This project was initially designed for the Raspberry Pi, but it's versatile enough to run on other compatible devices.
* **Operating System:** Tested on Arch Linux ARM and Raspberry Pi OS, but should work on other supported distributions.
* **vcgencmd** ([https://github.com/raspberrypi/utils](https://github.com/raspberrypi/utils)): Required for certain functionalities.

## Configuration

Customize RaspController's behavior using the `conf.yml` file:

```yaml
AUTH_TOKEN: "your_strong_secret_key" # Replace with a secure key 
DB_DIR: "/tmp/rosedb"                # Path to store the RoseDB database
PORT: 8080                          # Port for the web server
SHARE_DIR: "/home/rasp/public"      # Directory for shared files
```

## API Routes

RaspController exposes a RESTful API (all routes prefixed with `/api`):

**Information**

* **`/api/info`:** Retrieve general system information (RAM, CPU, disk, etc.).
* **`/api/info/ps`:** List running processes.

**GPIO**

* **`/api/gpio`:**  List used GPIO pins.
* **`/api/gpio/all`:**  List available GPIO pins.
* **`/api/gpio/:id`:** Get details about a specific GPIO pin.
* **`/api/gpio/:id` (PATCH):** Update GPIO configuration (example JSON body):
   ```json
   {
       "direction": "out", 
       "value": 1
   }
   ```

## Installation

1. **Create a project directory:** e.g., `/opt/raspc`.
2. **Transfer files:** Place the RaspController binary and `conf.yml` in the created directory.
3. **Set up systemd service:**  Copy the provided service file to `/etc/systemd/system/` and enable it: 
   ```bash
   systemctl enable --now raspc.service
   ```
4. **View Logs:** Access logs using `journalctl -u raspc`.

## Additional Notes

* Consider security improvements like HTTPS and proxy authentication when exposing RaspController to the Internet.
* For detailed usage and customization, explore the project's source code repository. 
