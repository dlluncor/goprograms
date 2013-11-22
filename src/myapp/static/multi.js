var multi = {
  state: null
};

multi.initState = function() {
  multi.state = {
    table: '',
    token: '',
    username: '',
    round: 1
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

SolutionModel = function(obj) {
  this.obj = obj;
};

SolutionModel.prototype.getLines = function() {
  return BoardGen.unjoin(this.obj.Table);
};

SolutionModel.prototype.getAnswers = function() {
  return this.obj.Answers;
};

 // Logic that deals with responding to user requests.
multi.handleMessage = function(resp) {
    if (resp.Action == 'join') {
      ctrl.console.multiPrint('Set up initial state of the table for the user.');
      var gameM = new GameModel(resp.Payload);
      ctrl.table.updateUi(gameM);
    }
    else if (resp.Action == 'startGame') {
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
      ctrl.table.boardC.useSolutions(new SolutionModel(info));
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
    } else if (resp.Action == 'gameEnded') {
      // Enable the "Start game" button again and provide some message
      // of how the users did.
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
  	window.console.log('Client channel opened.');
    multi.sendMessage('/opened');
  };
  
multi.onMessage = function(m) {
    msg = JSON.parse(m.data);
    window.console.log(msg);
    multi.handleMessage(msg);
};

multi.openChannel = function(token) {
	multi.state.token = token;
	var channel = new goog.appengine.Channel(token);
	var handler = {
	  'onopen': multi.onOpened,
	  'onmessage': multi.onMessage,
	  'onerror': function() {
	  	window.console.log('Error with the channel.');
	  },
	  'onclose': function() {
	  	window.console.log('Channel closing.');
	  }
	};
	var socket = channel.open(handler);
	socket.onopen = multi.onOpened;
	socket.onmessage = multi.onMessage;
};

multi.initConnection = function(user, table) {
  multi.state.username = user;
  multi.state.table = table;
  var params = '?';
  params += 'u=' + multi.state.username;
  params += '?g=' + multi.state.table;
  $.ajax('/getToken' + params).done(function(data) {
    multi.openChannel(data);
  });
};