library string;

bool stringIsTrue(String str) {
  if (str == null) return false;

  str = str.trim();

  if (str.isEmpty) return false;

  str = str.toLowerCase();

  // f* (false)
  if (str.startsWith("f")) return false;

  // of* (off is false, on is true)
  if (str.startsWith("of")) return false;

  int i;

  try {
    i = int.parse(str);
  } catch (e) {
    // not an integer
    return true;
  }

  assert(i != null);

  return i != 0;
}
