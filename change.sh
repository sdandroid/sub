#!/bin/bash

array=($(echo $ServerList | tr ";" "\n"))
idx=1
c=0
cp=0

for(( i=0;i<${#array[@]};i++))
do
    nidx=$(($idx+$i))
    old=$(uci get xray_fw3.@servers[$nidx].server)
    new=${array[i]}
    if [[ "$old" != "$new" ]];then
      uci set xray_fw3.@servers[$nidx].server=$new
      c=1
    fi
done;

for(( i=0;i<${#array[@]};i++))
do
    nidx=$(($idx+$i))
    old=$(uci get xray_fw3.@servers[$nidx].server_port)
    if [[ "$old" != "$Port" ]];then
      uci set xray_fw3.@servers[$nidx].server_port=$Port
      cp=1
    fi
done;

if [[ $c == 1  ||  $cp == 1 ]];then
  uci commit xray_fw3
  /etc/init.d/xray_fw3 restart
  echo "end"
fi
