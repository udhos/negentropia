
@rem assign zone
\redisbin\redis-cli hset everton.marques@gmail.com location z:0
\redisbin\redis-cli hset udhos0@gmail.com          location z:1

@rem add shader program
\redisbin\redis-cli hset p:simpleTexturizer vertexShader /shader/simpleTex_vs.txt
\redisbin\redis-cli hset p:simpleTexturizer fragmentShader /shader/simpleTex_fs.txt

@rem add zone ------------------------------------------------------------

\redisbin\redis-cli hset z:0 backfaceCulling true
@rem \redisbin\redis-cli hset z:0 skyboxURL /skybox/skybox_galaxy.json
@rem \redisbin\redis-cli hset z:0 skyboxURL /skybox/skybox_alien.json
\redisbin\redis-cli hset z:0 skyboxURL /skybox/skybox_sky30.json
\redisbin\redis-cli hset z:0 programName p:simpleTexturizer
\redisbin\redis-cli hset z:0 cameraCoord 0.0,0.0,90.0

@rem add instance list to zone
\redisbin\redis-cli hset z:0 instanceList l:0
\redisbin\redis-cli del l:0
\redisbin\redis-cli sadd l:0 m:0 m:1 m:2

@rem add object/model o:airship
\redisbin\redis-cli hset o:airship objURL /obj/airship.obj
\redisbin\redis-cli hset o:airship programName p:simpleTexturizer
\redisbin\redis-cli hset o:airship directionFront 5.0,0.0,0.0
\redisbin\redis-cli hset o:airship directionUp 0.0,5.0,0.0

@rem add object/model o:old_house
\redisbin\redis-cli hset o:old_house objURL /obj/old_house.obj
\redisbin\redis-cli hset o:old_house programName p:simpleTexturizer
\redisbin\redis-cli hset o:old_house directionFront 40.0,0.0,0.0
\redisbin\redis-cli hset o:old_house directionUp 0.0,40.0,0.0

@rem create instance m:0
\redisbin\redis-cli hset m:0 obj o:airship
\redisbin\redis-cli hset m:0 coord 0.0,0.0,0.0
\redisbin\redis-cli hset m:0 scale 1.0

@rem create instance m:1
\redisbin\redis-cli hset m:1 obj o:airship
\redisbin\redis-cli hset m:1 coord 0.0,7.0,0.0
\redisbin\redis-cli hset m:1 scale .5

@rem create instance m:2
\redisbin\redis-cli hset m:2 obj o:old_house
\redisbin\redis-cli hset m:2 coord -50.0,0.0,0.0
\redisbin\redis-cli hset m:2 scale 1.0

@rem add zone ------------------------------------------------------------

\redisbin\redis-cli hset z:1 backfaceCulling true
@rem \redisbin\redis-cli hset z:1 skyboxURL /skybox/skybox_galaxy.json
\redisbin\redis-cli hset z:1 skyboxURL /skybox/skybox_alien.json
@rem \redisbin\redis-cli hset z:1 skyboxURL /skybox/skybox_sky30.json
\redisbin\redis-cli hset z:1 programName p:simpleTexturizer
\redisbin\redis-cli hset z:1 cameraCoord 0.0,0.0,90.0

@rem add instance list to zone
\redisbin\redis-cli hset z:1 instanceList l:1
\redisbin\redis-cli del l:1
\redisbin\redis-cli sadd l:1 m:2

@rem eof
