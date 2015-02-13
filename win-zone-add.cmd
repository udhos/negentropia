@rem send commands from redis-zone-add.txt to redis server thru redis-cli

find /v "#" < \tmp\devel\negentropia\redis-zone-add.txt | \redis\redis-cli -x

@rem eof
