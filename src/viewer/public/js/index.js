const EventType = {
	START_CALCULATION: "start_calculation",
	NEW_CALCULATION: "new_calculation",
};

class Event {
	constructor(type, content) {
		this.type = type;
		this.content = content;
	}
}

class StartCalculationRequest {
	constructor(equation) {
		this.equation = equation;
	}
}

// Response received from backend when a new calculation must be shown
class StartCalculationResponse {
	constructor(id, created_time, equation, progress, result) {
		this.id = id;
		this.created_time = created_time;
		this.equation = equation;
		this.progress = progress;
		this.result = result;
	}
}

var routing = {
	// routeEvent handles the incoming events.
	routeEvent: (evt) => {
		if (!utils.isSomething(evt.type)) {
			console.log("Failed to route event because type is empty:", evt);
			return;
		}

		switch (evt.type) {
			case EventType.NEW_CALCULATION:
				console.log("New calculation event type recieved")
				const calc = Object.assign(new StartCalculationResponse, evt.content);
				ui.appendCalculation(calc)
				break;
			default:
				console.log("Unsupported event type received:", evt)
				break;
		}
	},

	// sendEvent ships an event to the backend using websocket
	sendEvent: (type, content) => {
		const event = new Event(type, content);
		websocket.conn.send(JSON.stringify(event));
		console.log("Event sent to server:", type)
	}
}

var ui = {
	// createNewCalculation sends a new equation to the server and starts the calculation.
	startCalculation: (evt) => {
		evt.preventDefault();
		var fname = "eq-text";
		var eq = document.getElementById(fname);
		if (utils.isSomething(eq.value)) {
			let outEvent = new StartCalculationRequest(eq.value)
			routing.sendEvent(EventType.START_CALCULATION, outEvent)
		}
		document.getElementById(fname).value = '';
	},

	appendCalculation: (calc) => {
		console.log("nice");
		console.log(calc);
		// TODO: Insert js code to create a calc in the board
	}
}

var utils = {
	// isSometing returns true when a var is not null|undefined|""
	isSomething: (str) => str ? true : false
}

var websocket = {
	// Access the websocket connection. 'connect' must be called first.
	conn: undefined,

	// Check if the browser supports websocket
	isBrowserSupported: () => {
		if (window["WebSocket"]) {
			return true;
		}
		return false;
	},

	connect: () => {
		if (!websocket.isBrowserSupported(window)) {
			console.log("Browser does not support websockets FAIL");
			return;
		}
		console.log("Browser supports websockets OK");

		var conn = new WebSocket("ws://" + document.location.host + "/ws");
		websocket.conn = conn;

		/* Register handlers of incoming events */

		conn.onopen = function() {
			console.log("Websocket connection established")
		}

		conn.onclose = function() {
			console.log("Websocket connection closed")
		}

		conn.onerror = function(ev) {
			console.log("Websocket connection error happened:", ev)
		}

		conn.onmessage = function(ev) {
			const data = JSON.parse(ev.data);
			const event = Object.assign(new Event, data);
			routing.routeEvent(event);
		}
	}
}

// Initial calls when page is loaded
window.onload = function() {
	websocket.connect();
	document.getElementById("eq-form").addEventListener('submit', ui.startCalculation)
};

