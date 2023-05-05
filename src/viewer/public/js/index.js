/*
	JS used to interact with server using websocket.
*/

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
	conn.onmessage = function(e) {
		console.log(e);
		const data = JSON.parse(e.data);
		const event = Object.assign(new Event, data);
		routeEvent(event);
	}
};

/*
	Handling of ingoing and outgoing websocket messages.
*/

class Event {
	constructor(type, payload) {
		this.type = type;
		this.payload = payload;
	}
}

// routeEvent handles the incoming event properly.
function routeEvent(event) {
	if (!isSomething(event.type)) {
		console.log("failed to route event because type is empty:", event);
		return;
	}

	switch (event.type) {
		case "new_equation":
			console.log("new equation")
			break;
		default:
			console.log("unsupported type")
			break;
	}
}

// sendEvent ships an event to the backend using websocket
function sendEvent(type, payload) {
	const event = new Event(type, payload);
	conn.send(JSON.stringify(event));
	console.log("event sent to server: ", type)
}

// sendEquation sends a new equation to the server to be shown in the status page.
function sendEquation() {
	var fname = "equation";

	var eq = document.getElementById(fname);
	if (isSomething(eq.value)) {
		conn.send(eq.value)
	}

	document.getElementById(fname).value = '';
}

// isSometing returns true when a var is not null|undefined|""
const isSomething = (str) => str ? true : false

