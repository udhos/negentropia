library shader;

import 'dart:html';
import 'dart:web_gl';
import 'dart:convert';
import 'dart:typed_data';
import 'dart:math' as math;

import 'package:vector_math/vector_math.dart';
import 'package:vector_math/vector_math_geometry.dart';
import 'package:vector_math/vector_math_lists.dart';
import 'package:game_loop/game_loop_html.dart';
import 'package:obj/obj.dart';

import 'camera.dart';
import 'texture.dart';
import 'asset.dart';
import 'logg.dart';
import 'selection.dart';

part 'buffer.dart';
part 'shader_tex.dart';
part 'picker.dart';
part 'solid.dart';

class ShaderProgram {
  Program program;
  RenderingContext gl;
  String programName;
  int a_Position;
  UniformLocation u_MV;
  UniformLocation u_P;
  bool shaderReady = false;

  List<Model> modelList = new List<Model>();

  ShaderProgram(this.gl, this.programName);

  /*
  void initContext(RenderingContext gl, Map<String,Texture> textureTable) {
  }
  */

  void getLocations() {
    a_Position = gl.getAttribLocation(program, "a_Position");
    u_MV = gl.getUniformLocation(program, "u_MV");
    u_P = gl.getUniformLocation(program, "u_P");
  }

  void fetch(Map<String, Shader> shaderCache, String vertexShaderURL,
      String fragmentShaderURL) {
    Shader compileShader(
        String shaderURL, String shaderSource, int shaderType) {
      Shader shader = gl.createShader(shaderType);
      gl.shaderSource(shader, shaderSource);
      gl.compileShader(shader);
      bool parameter =
          gl.getShaderParameter(shader, RenderingContext.COMPILE_STATUS);
      if (!parameter) {
        String infoLog = gl.getShaderInfoLog(shader);
        err("compileShader: compilation FAILURE: $shaderURL: info=$infoLog");
        if (gl.isContextLost()) {
          err("compileShader: compilation FAILURE: $shaderURL: info=$infoLog: context is lost");
        }
        return null;
      }

      if (shaderCache[shaderURL] != null) {
        err("compileShader: " +
            shaderURL +
            ": FIXME: overwriting shader cache");
      }
      shaderCache[shaderURL] = shader;

      return shader;
    }

    Shader vertexShader, fragmentShader;

    void tryLink() {
      if (vertexShader == null || fragmentShader == null) {
        // not ready to link
        return;
      }

      Program p = gl.createProgram();
      gl.attachShader(p, vertexShader);
      gl.attachShader(p, fragmentShader);
      gl.linkProgram(p);
      bool parameter = gl.getProgramParameter(p, RenderingContext.LINK_STATUS);
      if (!parameter) {
        String infoLog = gl.getProgramInfoLog(p);
        err("tryLink: shader program link FAILURE: $infoLog");
        if (gl.isContextLost()) {
          err("tryLink: shader program link FAILURE: $infoLog: context is lost");
        }
        return;
      }

      this.program = p;

      getLocations();

      shaderReady = true;
    }

    void fetchVertexShader() {
      String url = vertexShaderURL;

      var requestVert = new HttpRequest();
      requestVert.open("GET", url);
      requestVert.onLoad.listen((ProgressEvent e) {
        String response = requestVert.responseText;
        if (requestVert.status != 200) {
          err("vertexShader: url=$url: error: [$response]");
          return;
        }
        vertexShader =
            compileShader(url, response, RenderingContext.VERTEX_SHADER);
        tryLink();
      });
      requestVert.onError.listen((e) {
        err("vertexShader: url=$url: error: [$e]");
      });
      requestVert.send();
    }

    void fetchFragmentShader() {
      String url = fragmentShaderURL;

      var requestFrag = new HttpRequest();
      requestFrag.open("GET", url);
      requestFrag.onLoad.listen((ProgressEvent e) {
        String response = requestFrag.responseText;
        if (requestFrag.status != 200) {
          err("fragmentShader: url=$url: error: [$response]");
          return;
        }
        fragmentShader =
            compileShader(url, response, RenderingContext.FRAGMENT_SHADER);
        tryLink();
      });
      requestFrag.onError.listen((e) {
        err("fragmentShader: url=$url: error: [$e]");
      });
      requestFrag.send();
    }

    vertexShader = shaderCache[vertexShaderURL];
    if (vertexShader == null) {
      fetchVertexShader();
    }

    fragmentShader = shaderCache[fragmentShaderURL];
    if (fragmentShader == null) {
      fetchFragmentShader();
    }

    tryLink();
    // needed when both vertexShader and fragmentShader are found in cache
  }

  void addModel(Model newModel) {
    Model m = findModelByName(newModel.modelName);
    if (m != null) {
      err("Model.addModel: existing model modelName=${m.modelName}");
      return;
    }
    this.modelList.add(newModel);
  }

  Model findModelByName(String name) {
    Model mod;
    try {
      mod = modelList.firstWhere((m) {
        return m.modelName == name;
      });
    } on StateError {
      assert(mod == null); // not found
    }
    return mod;
  }

  Instance findInstance(String id) {
    Instance i;
    try {
      modelList.firstWhere((m) {
        i = m.findInstance(id);
        return i != null;
      });
    } on StateError {
      assert(i == null); // not found
    }
    return i;
  }

  void drawModels(GameLoopHtml gameLoop, Camera cam, Matrix4 pMatrix) {
    if (!shaderReady) {
      return;
    }

    gl.useProgram(program);
    gl.enableVertexAttribArray(a_Position);

    // send perspective projection matrix uniform
    gl.uniformMatrix4fv(this.u_P, false, pMatrix.storage);

    modelList.forEach((Model m) => m.drawInstances(gameLoop, this, cam));

    // clean up
    gl.bindBuffer(RenderingContext.ARRAY_BUFFER, null);
    gl.bindBuffer(RenderingContext.ELEMENT_ARRAY_BUFFER, null);

    //gl.disableVertexAttribArray(a_Position); // needed ??
  }

  void update(GameLoopHtml gameLoop) {
    modelList.forEach((m) => m.update(gameLoop));
  }
}
