```bash
#!/bin/bash

echo "ServerList:"$ServerList
echo "$Service:"$Service
echo "Port:"$Port

array=($(echo $ServerList | tr ";" "\n"))
idx=1
c=0
cp=0

for(( i=0;i<${#array[@]};i++))
do
    nidx=$(($idx+$i))
    old=$(uci get $Service.@servers[$nidx].server)
    new=${array[i]}
    if [[ "$old" != "$new" ]];then
      uci set $Service.@servers[$nidx].server=$new
      c=1
    fi
done;

for(( i=0;i<${#array[@]};i++))
do
    nidx=$(($idx+$i))
    old=$(uci get $Service.@servers[$nidx].server_port)
    if [[ "$old" != "$Port" ]];then
      uci set $Service.@servers[$nidx].server_port=$Port
      cp=1
    fi
done;

if [[ $c == 1  ||  $cp == 1 ]];then
  uci commit $Service
  /etc/init.d/$Service restart
  echo "end"
fi


