#!/bin/sh

URL=$1

echo "Sending data to $URL"
echo "Sending configs..."
curl -X POST -F '_action=1' -F 'username=user1' -F 'password=password1' $URL/v1/config
curl -X POST -F '_action=2' -F 'username=user2' -F 'password=password2' $URL/v1/config
curl -X POST -F '_action=3' -F 'username=user3' -F 'password=password3' $URL/v1/config
curl -X POST -F '_action=4' -F 'username=user4' -F 'password=password4' $URL/v1/config

echo "\nSending schedules..."
curl -X POST -F 'label=schedule1' -F 'crontab=8 * * * *' $URL/v1/schedule
curl -X POST -F 'label=schedule2' -F 'crontab=8 * * * *' $URL/v1/schedule
curl -X POST -F 'label=schedule3' -F 'crontab=8 * * * *' $URL/v1/schedule
curl -X POST -F 'label=schedule4' -F 'crontab=8 * * * *' $URL/v1/schedule

echo "\nSending probes..."
curl -X POST -F 'config_id=1' -F 'schedule_id=1' $URL/v1/probe
curl -X POST -F 'config_id=2' -F 'schedule_id=2' $URL/v1/probe
curl -X POST -F 'config_id=3' -F 'schedule_id=3' $URL/v1/probe
curl -X POST -F 'config_id=4' -F 'schedule_id=4' $URL/v1/probe

echo "\nSending agents..."
curl -X POST -F 'label=Test\ Agent' -F 'ip=127.0.0.1' $URL/v1/agent

echo "\nSending deploy..."
curl -X POST -F 'probe_id=1' -F 'agent_id=1' $URL/v1/deploy
curl -X POST -F 'probe_id=2' -F 'agent_id=2' $URL/v1/deploy
curl -X POST -F 'probe_id=3' -F 'agent_id=3' $URL/v1/deploy
