#!/usr/bin/bash
# Raspberry Pi Controller Install Script

SERVICE_PATH="/etc/systemd/system/"
SERVICE_NAME="raspc.service"
APP_PATH="/opt/raspc/"
APP_NAME="raspc"
BUILD_PATH="../cmd/$APP_NAME"

function make_build {
    # Entrar no diretório de construção
    cd "$BUILD_PATH" || exit

    # Construir para Linux ARM7
    GOOS=linux GOARCH=arm GOARM=7 go build -o "$APP_NAME"

    # Retornar ao diretório do pacote
    cd - || exit
}

function install_app {
    # Entrar no diretório de construção
    cd "$BUILD_PATH" || exit

    # Criar diretório do aplicativo
    mkdir -p "$APP_PATH"

    # Copiar arquivos do aplicativo
    cp "$APP_NAME" "$APP_PATH"
    cp "conf.yml" "$APP_PATH"

    # Retornar ao diretório do pacote
    cd - || exit

    # Copiar arquivo de serviço
    cp "./$SERVICE_NAME" "$SERVICE_PATH"
}

function install_service {
    # Criar arquivo de serviço
    cat <<EOT >"$SERVICE_PATH$SERVICE_NAME"
[Unit]
Description=Serviço RaspController
After=network.target

[Service]
ExecStart=/opt/raspc/raspc
Restart=always
RestartSec=3
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOT
}

case "$1" in
"make")
    make_build
    install_app
    install_service
    ;;
"install")
    install_app
    ;;
"build")
    make_build
    ;;
*)
    echo -e "RaspController Installer\n"
    echo -e "use: $0 make to build and install"
    echo -e "use: $0 build to build files"
    echo -e "use: $0 install to install"
    ;;
esac
