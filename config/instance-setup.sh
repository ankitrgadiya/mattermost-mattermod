#!/bin/bash

cd /home/ubuntu
echo "ALTER USER mmuser WITH PASSWORD 'mostest';" | sudo -u postgres psql
echo "DROP DATABASE mattermost;" | sudo -u postgres psql
echo "CREATE DATABASE mattermost;" | sudo -u postgres psql
wget https://example.com/mattermost-linux-amd64.tar.gz
tar -zxvf mattermost-linux-amd64.tar.gz
chown -R ubuntu mattermost
cd mattermost
sed -i'.bak1' 's|"DataSource": "[^"]*"|"DataSource": "postgres://mmuser:mostest@localhost:5432/mattermost?sslmode=disable\&connect_timeout=10"|g' config/s.Config.json
sed -i'.bak2' 's|"DriverName": "mysql"|"DriverName": "postgres"|g' config/s.Config.json
sed -i'.bak3' 's|"EnableDeveloper": false|"EnableDeveloper": true|g' config/s.Config.json
./bin/platform sampledata
rm -f ./logs/mattermost.log # Required because of permissions issue
start mattermost
