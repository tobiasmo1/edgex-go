#!/bin/bash
#
# Copyright (c) 2018
# Mainflux
#
# SPDX-License-Identifier: Apache-2.0
#

###
# Launches all EdgeX Go binaries (must be previously built).
#
# Expects that Consul and MongoDB are already installed and running.
#
###

DIR=$PWD
CMD=../cmd

#
# TJM: BEGIN Expressly assure that Consul and MongoDB are running.
cd ~/dev/edgexfoundry/developer-scripts/compose-files
# I've renamed docker-compose-delhi-0.7.1.yml to docker-compose.yml and modified mongodb port
EDGEX_COMPOSE_FILE=docker-compose.yml
# IF this needs updating, I either do it manually or can run launch script which updates /tmp folder, etc.
if [ -z $EDGEX_COMPOSE_FILE ]; then
  COMPOSE_FILENAME=docker-compose-delhi-0.7.1.yml
  COMPOSE_FILE=/tmp/${COMPOSE_FILENAME}
  COMPOSE_URL=https://raw.githubusercontent.com/edgexfoundry/developer-scripts/master/compose-files/${COMPOSE_FILENAME}
  echo "Pulling latest compose file..."
  curl -o $COMPOSE_FILE $COMPOSE_URL
else
  COMPOSE_FILE=$EDGEX_COMPOSE_FILE
fi

EDGEX_CORE_DB=${EDGEX_CORE_DB:-"mongo"}

echo "Starting Mongo"
docker-compose -f $COMPOSE_FILE up -d mongo

if [ ${EDGEX_CORE_DB} != mongo ]; then
  echo "Starting $EDGEX_CORE_DB for Core Data Services"
  docker-compose -f $COMPOSE_FILE up -d $EDGEX_CORE_DB
fi

echo "Starting consul"
docker-compose -f $COMPOSE_FILE up -d consul
echo "Populating configuration"
docker-compose -f $COMPOSE_FILE up -d config-seed

echo "Sleeping before launching remaining services"
sleep 15
# TJM: END Expressly assure that Consul and MongoDB are running



#docker-compose run consul &
#NOTE: I run mongo container on 27018, confirm runtime not populating my local mongo (27017)
#docker-compose run mongo &
cd $DIR
# TJM: END Expressly assure that Consul and MongoDB are running.

# Kill all edgex-* stuff
function cleanup {
	pkill edgex
}

###
# Support logging
###
cd $CMD/support-logging
# Add `edgex-` prefix on start, so we can find the process family
exec -a edgex-support-logging ./support-logging &
cd $DIR

###
# Core Command
###
cd $CMD/core-command
# Add `edgex-` prefix on start, so we can find the process family
exec -a edgex-core-command ./core-command &
cd $DIR

###
# Core Data
###
cd $CMD/core-data
exec -a edgex-core-data ./core-data &
cd $DIR

###
# Core Metadata
###
cd $CMD/core-metadata
exec -a edgex-core-metadata ./core-metadata &
cd $DIR

###
# Export Client
###
cd $CMD/export-client
exec -a edgex-export-client ./export-client &
cd $DIR

###
# Export Distro
###
cd $CMD/export-distro
exec -a edgex-export-distro ./export-distro &
cd $DIR

###
# Support Notifications
###
cd $CMD/support-notifications
# Add `edgex-` prefix on start, so we can find the process family
exec -a edgex-support-notifications ./support-notifications &
cd $DIR

###
# System Management Agent
###
cd $CMD/sys-mgmt-agent
# Add `edgex-` prefix on start, so we can find the process family
exec -a edgex-sys-mgmt-agent ./sys-mgmt-agent &
cd $DIR

# Support Scheduler
###
cd $CMD/support-scheduler
# Add `edgex-` prefix on start, so we can find the process family
exec -a edgex-support-scheduler ./support-scheduler &
cd $DIR

trap cleanup EXIT

while : ; do sleep 1 ; done