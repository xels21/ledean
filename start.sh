#!/usr/bin/env bash

cd `dirname "$0"`
./ledean -gpio_button=GPIO21 -led_count=92 -led_rows=2 -reverse_rows=1,1 -path2frontend="frontend" -direct_start -port=2211 -address="0.0.0.0"