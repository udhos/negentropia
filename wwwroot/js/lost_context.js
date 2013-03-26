
function stopAnimation() {
	//console.log("stopAnimation()");
	if ('reqId' in neg) {
		window.cancelAnimationFrame(neg.reqId); // from Brandon Jones game-shim.js
		delete neg.reqId;
	}
}
	
function cancelImageLoads() {
	// Ignore all ongoing image loads by removing their onload handler
	for (var i = 0; i < neg.ongoingImageLoads.length; i++) {
		neg.ongoingImageLoads[i].onload = undefined;
		neg.ongoingImageLoads[i].onerror = undefined;
		neg.ongoingImageLoads[i].onabort = undefined;
	}
	neg.ongoingImageLoads = [];		
}
	
function handleLostContext(event) {
	console.log("context LOST");
	event.preventDefault();
	stopAnimation();
	cancelImageLoads();
	
	//neg.textureTable = {}; // reload textures
}
	
function handleRestoredContext() {
	console.log("context RESTORED");
	// re-setup all your WebGL state and re-create all your WebGL resources when the context is restored.
	initContext();
}
	
function simulateLostContext(canvas) {
	console.log("lost context: simulate");
	canvas.loseContext();
}

function simulateRestoredContext(canvas) {
	canvas.restoreContext();
}

function toggleAutoRestore(canvas, autoRestoreCheck) {
	neg.autoRestore = autoRestoreCheck.checked;
	
	if (neg.autoRestore) {
		// Turn on automatic recovery
		console.log("auto restore: on");
		canvas.setRestoreTimeout(0);
	}
	else {
		// Turn off automatic recovery
		console.log("auto restore: off");
		canvas.setRestoreTimeout(-1);
	}
}

function initDebugLostContext(canvas) {

	if (!neg.debugLostContext) {
		return;
	}
		
	var control = document.getElementById("control");
	if (!control.appendChild) {
		return;
	}
	
	var autoRestoreCheck = document.createElement('input');
	autoRestoreCheck.type = "checkbox";
	//autoRestoreCheck.name = "auto restore";
	//autoRestoreCheck.value = "Auto Restore";
	autoRestoreCheck.id = "autoRestore";
	autoRestoreCheck.checked = true;
    autoRestoreCheck.onclick = function() { toggleAutoRestore(canvas, autoRestoreCheck); };
	var label = document.createElement('label')
	label.htmlFor = "autoRestore";
	label.appendChild(document.createTextNode('auto restore'));
	control.appendChild(autoRestoreCheck);
	control.appendChild(label);
	
    var loseContextButton = document.createElement("input");
    loseContextButton.type = 'button';
    loseContextButton.value = 'lose context';
    //loseContextButton.name = 'Lose Context';
    loseContextButton.onclick = function() { simulateLostContext(canvas); };
	control.appendChild(loseContextButton);

    var restoreContextButton = document.createElement("input");
    restoreContextButton.type = 'button';
    restoreContextButton.value = 'restore context';
    //restoreContextButton.name = 'Restore Context';
    restoreContextButton.onclick = function() { simulateRestoredContext(canvas); };
	control.appendChild(restoreContextButton);
	
	canvas.addEventListener("webglcontextlost", function(event) { handleLostContext(event); }, false);
	canvas.addEventListener("webglcontextrestored", handleRestoredContext, false);
}
