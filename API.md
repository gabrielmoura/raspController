# API Documentation

## Endpoints

### `/api/gpio`

- **Description:** Returns the status of all configured GPIO pins.
- **Method:** GET
- **Response:**

 ```json  
  {
  "pins": [
    {
      "pin_number": 1,
      "status": "HIGH"
    },
    {
      "pin_number": 2,
      "status": "LOW"
    }
  ]
}  
 ```  

### `/api/gpio/all`

- **Description:** Returns all GPIO pins from the GPIO chip.
- **Method:** GET
- **Response:**

 ```json  
  {
  "pins": [
    {
      "pin_number": 1,
      "name": "GPIO1"
    },
    {
      "pin_number": 2,
      "name": "GPIO2"
    }
  ]
}  
 ```  

### `/api/info`

- **Description:** Returns system information.
- **Method:** GET
- **Response:**

 ```json  
  {
  "cpu": "/api/info/cpu",
  "disk": "/api/info/disk",
  "gpio": "/api/info/gpio",
  "memory": "/api/info/mem",
  "network": "/api/info/net",
  "process": "/api/info/ps",
  "usb": "/api/info/usb"
}  
 ```  

### `/api/info/cpu`

- **Description:** Returns CPU information.
- **Method:** GET
- **Response:**

 ```json  
  {
  "model": "Intel Core i7",
  "cores": 4,
  "speed": "3.5 GHz"
}  
 ```  

### `/api/info/disk`

- **Description:** Returns disk information.
- **Method:** GET
- **Response:**

 ```json  

{
  "disks": {
    "percent": {
      "boot": 0.141898568375577,
      "home": 0.955070369719774,
      "root": 0.955070369719774
    },
    "total": {
      "boot": 534736896,
      "home": 239518875648,
      "root": 239518875648
    },
    "usage": {
      "boot": 75878400,
      "home": 228757381120,
      "root": 228757381120
    }
  },
  "reading_date": "2024-09-09 18:07:11"
}

```  

### `/api/info/gpio`

- **Description:** Returns list of available GPIOs.
- **Method:** GET
- **Response:**

 ```json 
 {
  "available_pins": [
    {
      "pin_number": 1,
      "name": "GPIO1"
    },
    {
      "pin_number": 2,
      "name": "GPIO2"
    }
  ]
} 
```

### `/api/info/mem`

- **Description: ** Returns memory information.
- **Method: ** GET
- **Response: **
```json

{
"memory": {
"free": 1638060032,
"total": 16646524928,
"used": 15008464896
},
"reading_date": "2024-09-09 18:06:39"
}

```  

### `/api/info/net`

- **Description:** Returns network information.
- **Method:** GET
- **Response:**

 ```json
  {
  "network": [
    {
      "interface": "enp2s0",
      "mac": "0a:e0:af:xx:xx:xx",
      "rx_bytes": 243957752,
      "tx_bytes": 101033968
    },
    {
      "interface": "lo",
      "mac": "00:00:00:00:00:00",
      "rx_bytes": 1916770,
      "tx_bytes": 1916770
    }
  ],
  "reading_date": "2024-09-09 18:04:37"
}
 ```  

### `/api/info/ps`

- **Description:** Returns process information.
- **Method:** GET
- **Response:**

 ```json  
  {
  "processes": [
    {
      "pid": 1234,
      "name": "nginx",
      "status": "running"
    },
    {
      "pid": 5678,
      "name": "mysql",
      "status": "running"
    }
  ]
}  
 ```  

### `/api/info/usb`

- **Description:** Returns list of USB devices.
- **Method:** GET
- **Response:**

 ```json  
  {
  "usb": [
    {
      "Bus": "001",
      "Device": "001",
      "ID": "1d6b:0002",
      "Vendor": "1d6b",
      "Product": "0002 Linux Foundation 2.0 root hub"
    }
  ]
}

```  

### `/api/share`

- **Description:** Returns a list of files contained in the sharing directory.
- **Method:** GET
- **Response:**  
  ```json { "files": [ "file1.txt", "file2.jpg", "file3.pdf" ] } ```

## Error Handling

- **Error Response Example:**

 ```json  
  {
  "error": "Invalid request",
  "message": "The provided parameters are incorrect."
}  
 ```