#!/bin/bash

cp /etc/configs/get-master.sh /sample.sh
chmod +x /sample.sh
/bin/bash -c /mysql_proxy
