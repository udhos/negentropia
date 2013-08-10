
@rem assign zone
\redisbin\redis-cli hset everton.marques@gmail.com location z:0

@rem add zone
\redisbin\redis-cli hset z:0 backfaceCulling true
\redisbin\redis-cli hset z:0 skyboxURL /skybox/skybox_galaxy.json
\redisbin\redis-cli hset z:0 programName p:simpleTexturizer
\redisbin\redis-cli hset p:simpleTexturizer vertexShader /shader/simpleTex_vs.txt
\redisbin\redis-cli hset p:simpleTexturizer fragmentShader /shader/simpleTex_fs.txt

@rem add instance list to zone
\redisbin\redis-cli hset z:0 instanceList l:0
\redisbin\redis-cli sadd l:0 m:0 m:1 m:2

@rem create instance m:0
\redisbin\redis-cli hset m:0 programName p:simpleTexturizer
\redisbin\redis-cli hset m:0 obj /obj/airship.obj
\redisbin\redis-cli hset m:0 coord 0.0,0.0,0.0
\redisbin\redis-cli hset m:0 scale 1.0

@rem create instance m:1
\redisbin\redis-cli hset m:1 programName p:simpleTexturizer
\redisbin\redis-cli hset m:1 obj /obj/airship.obj
\redisbin\redis-cli hset m:1 coord 0.0,4.0,0.0
\redisbin\redis-cli hset m:1 scale .5

@rem eof

