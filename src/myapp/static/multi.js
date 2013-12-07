var multi = {
  state: null
};

multi.initState = function() {
  multi.state = {
    table: '',
    token: '',
    username: '',
    round: 1,
    socket: null, // Current socket used to connect to the server.
    uiInitiated: false // Whether the UI has been initiated yet.
   };
};

GameModel = function(json) {
    this.obj = json; // JSON which underlies the game object.
};

GameModel.inArr = function(arr, el) {
    for (var i = 0; i < arr.length; i++) {
      if (arr[i] == el) {
      	return true;
      }	
    }
    return false;
};

GameModel.prototype.isStarted = function() {
  return GameModel.inArr(this.obj.States, 'justStarted');
};

// Return a map of user -> information about them, such as
// {points: number}
GameModel.prototype.getUsersInfo = function() {
  var len = this.obj.Users.length;
  var infoMap = {};
  for (var i = 0; i < len; i++) {
    var user = this.obj.Users[i];
    var info = {
      points: this.obj.Points[i]
    };
    infoMap[user] = info;
  }
  return infoMap;
};

// We know when the last round started (server time).
// We know when we just got this message (server time)
// How much time is left is a delta of that.
GameModel.prototype.getTimeLeft = function() {
  var roundFetched = this.obj.LastRoundFetched; // Everyone fetches 10 seconds before the next round starts. 
  var now = this.obj.Now;

  var delta = Math.floor(now - roundFetched);
  var config = ctrl.table.config;
  if (delta < config.betweenRound) {
  	// We haven't started the round yet!
  	var timeTilNextRound = config.betweenRound - delta;
  	return -timeTilNextRound;
  } else {
  	// This many seconds left in the round.
  	// Scenario.
  	// roundFetched: 20. now: 50. Each round is 120 seconds. 10 secs in between. 100 seconds left.
  	// 120 - 30 + 10 = 100
    var timeLeft = config.eachRound - delta + config.betweenRound;
  	return timeLeft;
  }
};

// Returns the current words found thus far in this puzzle.
// An array of objects with one field called words.
GameModel.prototype.getCurWordObjs = function() {
  if (!this.obj.CurWords) {
  	return [];
  }
  var wordObjs = [];
  this.obj.CurWords.forEach(function(word) {
    var wordObj = {
      word: word
    };
    wordObjs.push(wordObj);
  });
  return wordObjs;
};

TableInfo = function(obj) {
  this.obj = obj;
};

TableInfo.prototype.getLines = function() {
  return BoardGen.unjoin(this.obj.Table);
};

 // Logic that deals with responding to user requests.
multi.handleMessage = function(resp) {
    if (resp.Action == 'join') {
      // This is the case when you need an entire table refresh such as you joining
      // the table, or another using leaving.
      ctrl.console.multiPrint('Set up state of the table for the user.');
      var gameM = new GameModel(resp.Payload);
      ctrl.table.updateUi(gameM);
      if (gameM.isStarted() && multi.state.uiInitiated == false) {
      	// User is jumping into the game right in the middle, we need to fast-forward
      	// his state to the current correct state!
      	ctrl.console.multiPrint('Fast forward user to running game.');
      	ctrl.table.fastForwardUi(gameM);
      }
      multi.state.uiInitiated = true;  // UI has been initiated at least once.
    }
    else if (resp.Action == 'startGame') {
      ctrl.table.startButtonDisabled(true);
      ctrl.console.multiPrint('Game about to start once we generate tables.');
      // This guy can generate the tables for everyone for each round for now
      // even though the backend should be doing that, and then the server
      // can notify everyone when to start the round 1 (everyone should
      // be synchronized at that point).
      var params = {
      	'table1': BoardGen.join(BoardGen.generateBoard(1)),
      	'table2': BoardGen.join(BoardGen.generateBoard(2)),
      	'table3': BoardGen.join(BoardGen.generateBoard(3)),
      	'table4': BoardGen.join(BoardGen.generateBoard(4))
      };
      multi.sendMessage('sendTables', params);
    }
    else if (resp.Action == 'startTimers') {
      ctrl.console.multiPrint('Starting my timers.');
      ctrl.table.startRound();
    }
    else if (resp.Action == 'aboutToStartRound') {
      var info = resp.Payload;
      ctrl.console.multiPrint('Got info for next round');
      // The timer is still counting down to 10...but now we have all
      // the solutions for the puzzle, what are complete words.
      ctrl.table.boardC.useSolutions(new TableInfo(info));
    }
    //else if (resp.Action == 'endRound') {
      // The round just ended for someone, so we all end the round at the
      // same time and then start our counters at 10.
    //}
    else if (resp.Action == 'wordUpdate') {
      // Here we can handle situations such as:
      // Who is the update for (user - string).
      // What happened - "rejected" means your word got rejected.
      // Points - this is the delta for what the consequence of this action
      // is.
      // Cur points - the points that the user as a result of this action
      // (better than returning deltas in case clients don't get all
      //  the messages.)
      // So we can handle other users giving us updates or we can handle
      // just us getting these updates.
      ctrl.table.boardC.wordUpdate(resp.Payload);
    } else if (resp.Action == 'gameEnded') {
      // Enable the "Start game" button again and provide some message
      // of how the users did.
      ctrl.table.startButtonDisabled(false);
    }
};

multi.sendMessage = function(path, opt_params) {
    path += '?g=' + multi.state.table;
    path += '&t=' + multi.state.token;
    path += '&u=' + multi.state.username;
    if (opt_params) {
      for (var param in opt_params) {
        path += '&' + param + '=' + opt_params[param];	
      }
    }
    var xhr = new XMLHttpRequest();
    xhr.open('POST', path, true);
    xhr.send();
};

multi.onOpened = function() {
	// Let the server notify the user of its connected state.
  	window.console.log('Client channel opened.');
    multi.sendMessage('/opened');
};
  
multi.onMessage = function(m) {
    msg = JSON.parse(m.data);
    window.console.log(msg);
    multi.handleMessage(msg);
};

multi.onClose = function(opt_msg) {
  var msg = '';
  if (opt_msg) {
  	msg = opt_msg;
  }
  window.console.log('Channel closing because ' + msg);
  // Notify server that I am exiting this table so remove me from the list of
  // users that are playing.
  multi.state.socket.close();
}

multi.onError = function() {
  window.console.log('Error with the channel.');
};

multi.openChannel = function(token) {
	multi.state.token = token;
	var channel = new goog.appengine.Channel(token);
	var handler = {
	  'onopen': multi.onOpened,
	  'onmessage': multi.onMessage,
	  'onerror': multi.onError,
	  'onclose': multi.onClose
	};
	var socket = channel.open(handler);
	socket.onopen = multi.onOpened;
	socket.onmessage = multi.onMessage;
	// Enable if sessions ever expire or got past 2 hours??
	socket.onerror = multi.onError;
	multi.state.socket = socket;
	//socket.onclose = multi.onClose;
};

multi.initConnection = function(user, table, token) {
  multi.state.username = user;
  multi.state.table = table;
  multi.state.uiInitiated = false;
  multi.openChannel(token);
};

window.onbeforeunload = function() {
    websocket.onclose = function () {}; // disable onclose handler first
    multi.onClose('browser exiting');
};