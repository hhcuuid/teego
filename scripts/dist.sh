#!/bin/bash

# Deploy
# Build with make dist
# Copy the executable and the market.sh script
# Create a crontab:
#   crontab -e
#   # copy the following
#   0 20 * * * /home/holden_uuid/market/market.sh
# Stream logs to cloud logging:
#   copy dest/config.yaml to /etc/google-cloud-ops-agent/config.yaml
#   restart the agent sudo systemctl restart google-cloud-ops-agent

DIST=dist
BASE=/home/holden_uuid/market
GOOS=linux 
GOARCH=amd64

go build -o ${DIST}/market-linux-amd64 main.go

cat << EOF > "${DIST}/market.sh"
#!/bin/bash
export FINNUB_API_KEY=${FINNUB_API_KEY}
export GH_TOKEN=${GH_TOKEN}
nohup ${BASE}/market-linux-amd64 &
EOF

cat << EOF > "${DIST}/config.yaml"
logging:
  receivers:
    market:
      type: files
      include_paths:
      - ${BASE}/nohup.out
    service:
      pipelines:
      default_pipeline:
      receivers: [market]
EOF