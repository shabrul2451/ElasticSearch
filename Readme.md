# 🚀 Try Elasticsearch and Kibana locally

## 🏃‍♀️‍➡️ Getting started

### Setup

Run the `start-local` script using [curl](https://curl.se/):

```bash
curl -fsSL https://elastic.co/start-local | sh
```

This script creates an `elastic-start-local` folder containing:

- `docker-compose.yml`: Docker Compose configuration for Elasticsearch and Kibana
- `.env`: Environment settings, including the Elasticsearch password
- `start.sh` and `stop.sh`: Scripts to start and stop Elasticsearch and Kibana
- `uninstall.sh`: The script to uninstall Elasticsearch and Kibana

### 🌐 Endpoints

After running the script:

- Elasticsearch will be running at <http://localhost:9200>
- Kibana will be running at <http://localhost:5601>

The script generates a random password for the `elastic` user, displayed at the end of the installation and stored in the `.env` file.

### Local User and Password:

Username: elastic

Password: 7FAW0rS2

### Local Api Key
`b3IxMGlKUUIxc3VPOHRMNk9RMTc6bzQwSGlsNXRRV096Q1ZHckVnNG85Zw==`

### To Start
``
./start.sh
``

### To Stop
``
./stop.sh
``

## ⚙️ Customizing settings

To change settings (e.g., Elasticsearch password), edit the `.env` file. Example contents:

```bash
ES_LOCAL_VERSION=8.15.2
ES_LOCAL_URL=http://localhost:9200
ES_LOCAL_CONTAINER_NAME=es-local-dev
ES_LOCAL_DOCKER_NETWORK=elastic-net
ES_LOCAL_PASSWORD=hOalVFrN
ES_LOCAL_PORT=9200
KIBANA_LOCAL_CONTAINER_NAME=kibana-local-dev
KIBANA_LOCAL_PORT=5601
KIBANA_LOCAL_PASSWORD=YJFbhLJL
ES_LOCAL_API_KEY=df34grtk...==
```

> [!IMPORTANT]
> After changing the `.env` file, restart the services using `stop` and `start`:
>
> ```bash
> cd elastic-start-local
> ./stop.sh
> ./start.sh
> ```


# For Code  format
```gofumpt -w . && golines -w .```