#!/bin/bash
set -e
BINARY_NAME="beamstack"

download_binary() {
    echo "Downloading the binary..."
    curl -L -o /tmp/${BINARY_NAME} "https://raw.githubusercontent.com/beamstackproj/beamstack-cli/main/releases/latest/beamstack-linux-amd64"
    if [ ! -f /tmp/${BINARY_NAME} ]; then
        echo "Error: Binary download failed."
        exit 1
    fi
}

install_binary() {
    echo "Creating directories..."
    mkdir -p ~/."$BINARY_NAME"
    mkdir -p ~/."$BINARY_NAME"/config
    mkdir -p ~/."$BINARY_NAME"/profiles

    echo "Creating config file..."
    if [ ! -f ~/."$BINARY_NAME"/config/config.json ]; then
        cat <<EOF > ~/."$BINARY_NAME"/config/config.json
{
    "version": "1.0.0",
    "contexts": {}
}
EOF
        echo "Created new config file: ${CONFIG_FILE}"
    fi

# TODO: UPDATE version
    echo "Moving binary to /usr/local/bin..."
    sudo mv /tmp/${BINARY_NAME} /usr/local/bin/${BINARY_NAME}
    sudo chmod +x /usr/local/bin/${BINARY_NAME}

    echo "Setting permissions..."
    sudo chmod -R 777 ~/."$BINARY_NAME"

    echo "✅ ${BINARY_NAME} installed..!"
}

download_binary
install_binary
