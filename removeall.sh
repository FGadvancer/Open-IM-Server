#!/usr/bin/env bash

# Copyright © 2023 OpenIM. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

make stop
if [ $? -eq 0 ]; then
    docker-compose -f docker-compose-a.yml down && \
    rm -rf ./components/kafka/* &&\
    rm -rf ./components/redis/* &&\
    rm -rf ./components/mongodb/* && \
    rm ./logs/* && \
    docker-compose -f docker-compose-a.yml up -d && \
    make start

    nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-push --port 30171 -c /data/Open-IM-Server/config/ --prometheus_port 40171 >> openIM.log 2>&1 &
     nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-push --port 30172 -c /data/Open-IM-Server/config/ --prometheus_port 40172 >> openIM.log 2>&1 &
       nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-push --port 30173 -c /data/Open-IM-Server/config/ --prometheus_port 40173 >> openIM.log 2>&1 &
          nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-push --port 30174 -c /data/Open-IM-Server/config/ --prometheus_port 40174 >> openIM.log 2>&1 &
              nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-push --port 30175 -c /data/Open-IM-Server/config/ --prometheus_port 40175 >> openIM.log 2>&1 &
               nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-push --port 30176 -c /data/Open-IM-Server/config/ --prometheus_port 40176 >> openIM.log 2>&1 &
                 nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-push --port 30177 -c /data/Open-IM-Server/config/ --prometheus_port 40177 >> openIM.log 2>&1 &
#          nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-rpc-msg --port 30131 -c /data/Open-IM-Server/config --prometheus_port 40131 >> openIM.log 2>&1 &
#        nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-rpc-msg --port 30132 -c /data/Open-IM-Server/config --prometheus_port 40132 >> openIM.log 2>&1 &
#        nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-rpc-msg --port 30133 -c /data/Open-IM-Server/config --prometheus_port 40133 >> openIM.log 2>&1 &
#        nohup /data/Open-IM-Server/_output/bin/platforms/linux/amd64/openim-rpc-msg --port 30134 -c /data/Open-IM-Server/config --prometheus_port 40134 >> openIM.log 2>&1 &
fi