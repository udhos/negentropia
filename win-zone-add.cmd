@rem send commands from redis-zone-add.txt to redis server thru redis-cli

more \tmp\devel\negentropia\redis-zone-add.txt | find /v "#" | \redis\redis-cli -x

@rem eof
