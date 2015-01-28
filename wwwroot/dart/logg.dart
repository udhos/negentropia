library log;

final String _compile_time_env_debug = const String.fromEnvironment('DEBUG');

void logg_init() {
  log("compile time environment: DEBUG=$_compile_time_env_debug");
}

void log(String msg) {
  print(msg);
}

void warn(String msg) {
  log("WARN: $msg");
}

void err(String msg) {
  log("ERR: $msg");
}

void debug(String msg) {
  if (_compile_time_env_debug != null) {
    log("DEBUG: $msg");
  }
}

void fixme(String msg) {
  log("FIXME: $msg");
}
