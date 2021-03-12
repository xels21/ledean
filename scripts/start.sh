#!/bin/bash

# Absolute path to this script, e.g. /home/user/bin/foo.sh
SCRIPT=$(readlink -f "$0")
# Absolute path this script is in, thus /home/user/bin
SCRIPTPATH=$(dirname "$SCRIPT")

pkill ledean
$SCRIPTPATH/ledean -gpio_button=GPIO17 -led_count=92 -led_rows=2 -reverse_rows=0,1 -log_level=log -address=0.0.0.0 -path2frontend=$SCRIPTPATH/frontend
