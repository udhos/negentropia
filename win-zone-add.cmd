@set REDIS_BIN=\redis\redis-cli

@rem assign zone to user
%REDIS_BIN% hset everton.marques@gmail.com location z:0
%REDIS_BIN% hset udhos0@gmail.com          location z:1

@rem add shader program
%REDIS_BIN% hset p:simpleTexturizer vertexShader /shader/simpleTex_vs.txt
%REDIS_BIN% hset p:simpleTexturizer fragmentShader /shader/simpleTex_fs.txt

@rem add zone ------------------------------------------------------------

%REDIS_BIN% hset z:0 backfaceCulling true
@rem %REDIS_BIN% hset z:0 skyboxURL /skybox/skybox_galaxy.json
@rem %REDIS_BIN% hset z:0 skyboxURL /skybox/skybox_alien.json
%REDIS_BIN% hset z:0 skyboxURL /skybox/skybox_sky30.json
%REDIS_BIN% hset z:0 programName p:simpleTexturizer
%REDIS_BIN% hset z:0 cameraCoord 0.0,0.0,90.0

@rem add instance list to zone
%REDIS_BIN% hset z:0 instanceList l:0
%REDIS_BIN% del l:0
%REDIS_BIN% sadd l:0 m:0 m:1 m:2 m:3 m:4 m:5 m:6

@rem add object/model o:airship
%REDIS_BIN% hset o:airship objURL /obj/airship.obj
%REDIS_BIN% hset o:airship programName p:simpleTexturizer
%REDIS_BIN% hset o:airship modelFront 5.0,0.0,0.0
%REDIS_BIN% hset o:airship modelUp 0.0,5.0,0.0

@rem add object/model o:old_house
%REDIS_BIN% hset o:old_house objURL /obj/old_house.obj
%REDIS_BIN% hset o:old_house programName p:simpleTexturizer
%REDIS_BIN% hset o:old_house modelFront 40.0,0.0,0.0
%REDIS_BIN% hset o:old_house modelUp 0.0,40.0,0.0

@rem add object/model o:mars
%REDIS_BIN% hset o:mars objURL /obj/MarsPlanet.obj
%REDIS_BIN% hset o:mars programName p:simpleTexturizer
%REDIS_BIN% hset o:mars modelFront 20.0,0.0,0.0
%REDIS_BIN% hset o:mars modelUp 0.0,20.0,0.0

@rem add object/model o:bigearth
%REDIS_BIN% hset o:bigearth globeRadius 200.0
%REDIS_BIN% hset o:bigearth globeTextureURL /texture/earthmap1k.jpg
%REDIS_BIN% hset o:bigearth programName p:simpleTexturizer
%REDIS_BIN% hset o:bigearth modelFront 200.0,0.0,0.0
%REDIS_BIN% hset o:bigearth modelUp 0.0,200.0,0.0

@rem create instance
%REDIS_BIN% hset m:0 obj o:airship
%REDIS_BIN% hset m:0 coord 0.0,0.0,0.0
%REDIS_BIN% hset m:0 scale 1.0
%REDIS_BIN% hset m:0 mission rotateYaw
%REDIS_BIN% hset m:0 team alpha0
%REDIS_BIN% hset m:0 owner udhos0@gmail.com

@rem create instance
%REDIS_BIN% hset m:1 obj o:airship
%REDIS_BIN% hset m:1 coord 0.0,7.0,5.0
%REDIS_BIN% hset m:1 scale .5
%REDIS_BIN% hset m:1 team alpha1
%REDIS_BIN% hset m:1 mission hunt
%REDIS_BIN% hset m:1 owner everton.marques@gmail.com

@rem create instance
%REDIS_BIN% hset m:2 obj o:airship
%REDIS_BIN% hset m:2 coord 0.0,7.0,0.0
%REDIS_BIN% hset m:2 scale .5
%REDIS_BIN% hset m:2 team alpha1
%REDIS_BIN% hset m:2 mission hunt
%REDIS_BIN% hset m:2 owner everton.marques@gmail.com

@rem create instance
%REDIS_BIN% hset m:3 obj o:airship
%REDIS_BIN% hset m:3 coord 0.0,7.0,-5.0
%REDIS_BIN% hset m:3 scale .5
%REDIS_BIN% hset m:3 team alpha1

@rem create instance
%REDIS_BIN% hset m:4 obj o:old_house
%REDIS_BIN% hset m:4 coord -50.0,0.0,0.0
%REDIS_BIN% hset m:4 scale 1.0

@rem create instance
%REDIS_BIN% hset m:5 obj o:mars
%REDIS_BIN% hset m:5 coord 20.0,20.0,-20.0
%REDIS_BIN% hset m:5 scale 1.0

@rem create instance
%REDIS_BIN% hset m:6 obj o:bigearth
%REDIS_BIN% hset m:6 coord 500.0,500.0,-500.0
%REDIS_BIN% hset m:6 scale 1.0

@rem add zone ------------------------------------------------------------

%REDIS_BIN% hset z:1 backfaceCulling true
@rem %REDIS_BIN% hset z:1 skyboxURL /skybox/skybox_galaxy.json
%REDIS_BIN% hset z:1 skyboxURL /skybox/skybox_alien.json
@rem %REDIS_BIN% hset z:1 skyboxURL /skybox/skybox_sky30.json
%REDIS_BIN% hset z:1 programName p:simpleTexturizer
%REDIS_BIN% hset z:1 cameraCoord 0.0,0.0,90.0

@rem add instance list to zone
%REDIS_BIN% hset z:1 instanceList l:1
%REDIS_BIN% del l:1
%REDIS_BIN% sadd l:1 m:4

@rem eof
