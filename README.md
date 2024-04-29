# RaspController

Control your Raspberry Pi with a web interface.

## Tecnologias Utilizadas

- [RoseDB](https://github.com/rosedblabs/rosedb)
- [go-gpiocdev](https://github.com/warthog618/go-gpiocdev)
- [mdns](https://github.com/hashicorp/mdns)
- [Fiber](https://gofiber.io/)

## Obersevações

- O projeto foi desenvolvido para ser executado em um Raspberry Pi.
- O projeto usa RoseDB para armazenar os dados para hitórico e gráficos.
- As configurações devem ser especificadas no arquivo `config.json`.
- O projeto usa mDns para descoberta de serviços.

## Configurações

```yml
JWT_SECRET: "secretkey"
DB_DIR: "./db"
PORT: 8080
JWT_EXPIRES_IN: 3600
SHARE_DIR: "/home/rasp/public"
```

## Rotas

Todas as rotas começam com /api

### Informações

```http request
GET /info
```

### GPIO

#### Listar todos os GPIOs

```http request
GET /gpio
```

#### Obter informações de um GPIO

```http request
GET /gpio/:id
```

#### Atualizar um GPIO

```http request
PATCH /gpio/:id

```
