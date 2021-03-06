
function fetchFile(url, handler, opaque) {
	var client = new XMLHttpRequest();
	
	client.processHandler = handler;
	client.processOpaque = opaque;
	
	client.onreadystatechange = onFetchHandler;
	client.open("GET", url);
	client.send();
}

function onFetchHandler() {
	if (this.readyState == this.DONE) {
		if (this.status == 200 && this.responseText != null) {
			// fetch ok
			this.processHandler(this.processOpaque, this.responseText);
			return;
		}
		// fetch failure
		this.processHandler(this.processOpaque, null);
	}
}
