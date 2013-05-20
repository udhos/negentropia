// Thanks to CMS for the startsWith function
// http://stackoverflow.com/questions/646628/javascript-startswith/646643#646643

if (typeof String.prototype.startsWith !== 'function') {
  String.prototype.startsWith = function (str){
    return this.slice(0, str.length) === str;
  };
}
