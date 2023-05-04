/* Index page javascript */

// Check if the browser supports WebSocket
window.onload = function() {
	if (window["WebSocket"]) {
		console.log("Browser supports websockets OK");

		// Connect to websocket
		conn = new WebSocket("ws://" + document.location.host + "/ws");

	} else {
		console.log("Browser does not support websockets FAIL");
		alert("Browser does not support websockets FAIL");
	}
};

// Send calculation to backend using websocket
function sendEquation() {
	fname = "equation";

	var eq = document.getElementById(fname);
	if (eq != null) {
		console.log(eq.value);
		conn.send(eq.value)
	}

	document.getElementById(fname).value = '';
}
