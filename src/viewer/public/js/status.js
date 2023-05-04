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
