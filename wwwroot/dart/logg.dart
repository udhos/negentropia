library logg;

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
  if (const String.fromEnvironment('DEBUG') != null) {
    log("DEBUG: $msg");
  }
}

void fixme(String msg) {
  debug("FIXME: $msg");
}
