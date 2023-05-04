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

// Send request on 'enter'
// function enterKeyHandler(e, button) {
// 	e = e || window.event;
// 	if (e.key == 'Enter') {
// 		document.getElementById(button).click();
// 	}
// }

// Send calculation to backend using websocket
function sendEquation() {
	fname = "equation";
	var eq = document.getElementById(fname);
	if (eq != null) {
		console.log(eq.value);
	}
	document.getElementById(fname).value = '';

	return false;
}
