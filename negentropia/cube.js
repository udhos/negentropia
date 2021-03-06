
  var cubeVerticesCoord = [
       // Front face
       1.0,  1.0,  1.0, //v0
      -1.0,  1.0,  1.0, //v1
      -1.0, -1.0,  1.0, //v2
       1.0, -1.0,  1.0, //v3

       // Back face
       1.0,  1.0, -1.0, //v4
      -1.0,  1.0, -1.0, //v5
      -1.0, -1.0, -1.0, //v6
       1.0, -1.0, -1.0, //v7
       
       // Left face
      -1.0,  1.0,  1.0, //v8
      -1.0,  1.0, -1.0, //v9
      -1.0, -1.0, -1.0, //v10
      -1.0, -1.0,  1.0, //v11
       
       // Right face
       1.0,  1.0,  1.0, //12
       1.0, -1.0,  1.0, //13
       1.0, -1.0, -1.0, //14
       1.0,  1.0, -1.0, //15
       
        // Top face
        1.0,  1.0,  1.0, //v16
        1.0,  1.0, -1.0, //v17
       -1.0,  1.0, -1.0, //v18
       -1.0,  1.0,  1.0, //v19
       
        // Bottom face
        1.0, -1.0,  1.0, //v20
        1.0, -1.0, -1.0, //v21
       -1.0, -1.0, -1.0, //v22
       -1.0, -1.0,  1.0, //v23
  ];
  
var cubeColors = [
    // Front face
	1.0, 0.0, 0.0, 1.0, //v0 red
	1.0, 0.0, 0.0, 1.0, //v1 red
	1.0, 0.0, 0.0, 1.0, //v2 red
	1.0, 0.0, 0.0, 1.0, //v3 red

    // Back face
	0.0, 1.0, 0.0, 1.0, //v4 green
	0.0, 1.0, 0.0, 1.0, //v5 green
	0.0, 1.0, 0.0, 1.0, //v6 green
	0.0, 1.0, 0.0, 1.0, //v7 green

	// Left face
	0.0, 0.0, 1.0, 1.0,  //v8 blue
	0.0, 0.0, 1.0, 1.0,  //v9 blue
	0.0, 0.0, 1.0, 1.0,  //v10 blue
	0.0, 0.0, 1.0, 1.0,  //v11 blue

    // Right face	
	1.0, 1.0, 1.0, 1.0,  //v12 white
	1.0, 1.0, 1.0, 1.0,  //v13 white
	1.0, 1.0, 1.0, 1.0,  //v14 white
	1.0, 1.0, 1.0, 1.0,  //v15 white

	// Top face
	0.5, 0.5, 0.5, 1.0,  //v16 gray
	0.5, 0.5, 0.5, 1.0,  //v17 gray
	0.5, 0.5, 0.5, 1.0,  //v18 gray
	0.5, 0.5, 0.5, 1.0,  //v19 gray
	
    // Bottom face
	0.05, 0.05, 0.05, 0.05,  //v20 black
	0.05, 0.05, 0.05, 0.05,  //v21 black
	0.05, 0.05, 0.05, 0.05,  //v22 black
	0.05, 0.05, 0.05, 0.05,  //v23 black
	];
  
var cubeTextureCoord = [
    // Front face
    0.0, 0.0, //v0
    1.0, 0.0, //v1
    1.0, 1.0, //v2
    0.0, 1.0, //v3
    
    // Back face
    0.0, 1.0, //v4
    1.0, 1.0, //v5
    1.0, 0.0, //v6
    0.0, 0.0, //v7
    
    // Left face
    0.0, 1.0, //v8
    1.0, 1.0, //v9
    1.0, 0.0, //v10
    0.0, 0.0, //v11
    
    // Right face
    0.0, 1.0, //v12
    1.0, 1.0, //v13
    1.0, 0.0, //v14
    0.0, 0.0, //v15
    
    // Top face
    0.0, 1.0, //v16
    1.0, 1.0, //v17
    1.0, 0.0, //v18
    0.0, 0.0, //v19
    
    // Bottom face
    0.0, 1.0, //v20
    1.0, 1.0, //v21
    1.0, 0.0, //v22
    0.0, 0.0, //v23
  ];
  
  var cubeVertexIndices = [
	0, 1, 2,      0, 2, 3,    // Front face
    4, 6, 5,      4, 7, 6,    // Back face
    8, 9, 10,     8, 10, 11,  // Left face
    12, 13, 14,   12, 14, 15, // Right face
    16, 17, 18,   16, 18, 19, // Top face
    20, 22, 21,   20, 23, 22  // Bottom face
  ];

  /*
var cubeTextures = [
	"wood_128x128.jpg",
	"wood_floor_256.jpg",
	"wicker_256.jpg"
];
*/
