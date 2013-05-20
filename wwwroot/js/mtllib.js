
function mtllib_parse(data) {
	
	var mtllib = {};
	var curr_mtl_name = null;
		
	var lines = data.split('\n');
	var len = lines.length;
	var line_num = 0;
	for (var i = 0; i < len; ++i) {
		++line_num;
		var line = lines[i];
		
		// See str.js
		var newmtl = "newmtl ";
		if (line.startsWith(newmtl)) {
			curr_mtl_name = line.slice(newmtl.length);
			console.log("mtllib_parse: new material: " + curr_mtl_name);
			continue;
		}
		
		var prefix_map_Kd = "map_Kd ";
		if (line.startsWith(prefix_map_Kd)) {
			if (curr_mtl_name == null) {
				console.log("mtllib_parse: line=" + $line_num + ": map_Kd found for undefined material");
				continue;
			}
			var kd = line.slice(prefix_map_Kd.length);
			mtllib[curr_mtl_name] = {
				map_Kd: kd
			};
			console.log("mtllib_parse: " + curr_mtl_name + " map_Kd = " + kd);
			continue;
		}
	}

	return mtllib;
}