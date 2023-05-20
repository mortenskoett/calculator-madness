// Event types used for requests to backend.
const RequestEventType = {
	START_CALCULATION: "start_calculation",
};

// Event types used for received events.
const ResponseEventType = {
	NEW_CALCULATION: "new_calculation",
	ENDED_CALCULATION: "ended_calculation",
};

class Event {
	constructor(type, contents) {
		this.type = type;
		this.contents = contents;
	}
}

class StartCalculationRequest {
	constructor(equation) {
		this.equation = equation;
	}
}

// Response received from backend when a new calculation must be shown
class NewCalculationResponse {
	constructor(id, created_time, equation, progress) {
		this.id = id;
		this.created_time = created_time;
		this.equation = equation;
		this.progress = progress;
	}
}

// Response received from backend when a calculation is ended.
class EndCalculationResponse {
	constructor(id, result) {
		this.id = id;
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

		console.log("Event recieved:", evt.type)

		switch (evt.type) {
			case ResponseEventType.NEW_CALCULATION:
				const newcalc = Object.assign(new NewCalculationResponse, evt.contents);
				ui.prependCalculation(newcalc)
				break;
			case ResponseEventType.ENDED_CALCULATION:
				const endcalc = Object.assign(new EndCalculationResponse, evt.contents);
				ui.endCalculation(endcalc)
				break;
			default:
				console.log("Unsupported event type received:", evt)
				break;
		}
	},

	// sendEvent ships an event to the backend using websocket
	sendEvent: (type, contents) => {
		const event = new Event(type, contents);
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
			routing.sendEvent(RequestEventType.START_CALCULATION, outEvent)
		}
		document.getElementById(fname).value = '';
	},

	prependCalculation: (calc) => {
		console.log("Prepending calculation");
		var calcElem = document.createElement("div");
		calcElem.className = "result-elem";
		calcElem.id = calc.id;

		// The actual calculation element in the UI
		calcElem.innerHTML = `
		<div class="calc-result-div">
			<p id="calc-eq">${calc.equation}</p> = <p class="calc-result">?</p>
		</div>
		<div class="progress-bar">
			<img class="progress-icon" src="/public/images/cog.png" alt="Cog icon." width="20" height="20">
			<progress class="progress-indicator" value=${calc.progress.current} max=${calc.progress.outof}></progress>
		</div>
		`;

		document.getElementById("ongoing").prepend(calcElem);
	},

	// Update a calculation by moving it into ended and updating its look
	endCalculation: (calc) => {
		console.log("Ending calculation");

		var calcElem = document.getElementById(calc.id);
		calcElem.getElementsByClassName("calc-result")[0].innerText = `${calc.result}`;

		calcElem.getElementsByClassName("progress-icon")[0].src = "/public/images/done.png";
		calcElem.getElementsByClassName("progress-icon")[0].alt = "Checkmark icon.";

		calcElem.getElementsByClassName("progress-indicator")[0].value = 1;
		calcElem.getElementsByClassName("progress-indicator")[0].max = 1;

		document.getElementById("ended").prepend(calcElem);
	},
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

