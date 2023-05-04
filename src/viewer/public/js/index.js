/* Index page javascript */

window.onload = function() {
	// Check if the browser supports WebSocket
	if (!window["WebSocket"]) {
		console.log("Browser does not support websockets FAIL");
		alert("Browser does not support websockets FAIL");
		return;
	}
	console.log("Browser supports websockets OK");

	// Connect to websocket
	conn = new WebSocket("ws://" + document.location.host + "/ws");

	// Add a listener to the onmessage event
	conn.onmessage = function(evt) {
		console.log(evt);
	}
};

// sendEquation sends the equation to the backend using websocket
function sendEquation() {
	var fname = "equation";

	var eq = document.getElementById(fname);
	if (isSomething(eq.value)) {
		conn.send(eq.value)
	}

	// console.log(eq.value)
	// console.log(isSomething(eq.value))
	// console.log(conn)

	document.getElementById(fname).value = '';
}

// isSometing checks whether a var is not null|undefined|""
const isSomething = (str) => str ? true : false
