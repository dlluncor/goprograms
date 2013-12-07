
// Ajax with some printing.
var Jax = {};


Jax.ajax = function(url, doneCallback) {
  var aDiv = function(text) {
    return $('<div>' + text + '</div>');
  };

  var div = $('#rpcs');
  var before = new Date();
  var urlToDisp = url.substring(0, 11);
  var b4msg = 'Sent ' + urlToDisp + ' at ' + before.toLocaleString();
  div.append(aDiv(b4msg));
  $.ajax(url)
    .done(function(data) {
      	var after = new Date();
      	var seconds = (after - before) / 1000;
      	var afterMsg = 'Rec. ' + urlToDisp + ' at ' + after.toLocaleString() +
      	  ' (' + seconds + ' secs)';
      	div.append(aDiv(afterMsg));
      	doneCallback(data);
  });
};

// Connection to the backend and its info.
var backend = {};

// Returns to the callback a list of words.
backend.getAllWords = function(callback) {
  var before = new Date();

  var doneCallback = function(data) {
     var after = new Date();
     // Server sends back text a words separated by a comma.
     var words = data.split(',');
     callback(words);
  };

  Jax.ajax('/getallwords', doneCallback);
}

backend.solvePuzzle = function(answersCb, board, length) {

  var parseAnswersData = function(linesAsText) {
        var words = linesAsText.split(',');
        var answers = [];
        for (var i = words.length -1; i >= 0; i--) {
                var word = words[i];
                answers.push(word);
    }
        return answers;
  };

  var doneCallback = function(data) {
    var answers = parseAnswersData(data);
        window.console.log('success');
        window.console.log(data);
        answersCb(answers);
  };
  Jax.ajax('/wordracer_json?board=' + board + '&length=' + length,
      doneCallback);
};

BoardSolver = function(text) {
        this.text = text;
};

BoardSolver.prototype.solve = function(answersCb) {
  var text = this.text;
  var lines = text.split('\n');
  // Validate the board works.
  var length = lines[0].length;
  window.console.log(this.text);
  backend.solvePuzzle(answersCb, text, length);
};

// Round object to control keeping track of the
// round numbers and time left.
Round = function(boardC) {
  this.roundNumEl = $('#roundNum');
  this.timeLeftTextEl = $('#timeLeftText');
  this.timeLeftEl = $('#timeLeft');
  this.roundTimerInterval = null;

  // Coming from the logic controller.
  this.curRound = null;
  this.boardC = boardC;
};

Round.prototype.init = function() {
  this.boardC.init();
  this.curRound = 1;
};

// Method called when a round starts (delegates to board controllers and
// such).
Round.prototype.roundBegins = function(opt_secsLeft) {
  this.timeLeftTextEl.text('Time: ');
  this.startRoundTimer(opt_secsLeft);
  window.console.log('Go start game!');
  this.boardC.roundStart(this.curRound);
};

// 60 -> 1:00. 50 -> 0:50. 130 -> 2:10
Round.renderTime = function(secs) {
  var mins = Math.floor(secs/60.0);
  var secsLeft = secs - (mins * 60);
  var zeroPad = '';
  if (secsLeft < 10) {
    zeroPad = '0';
  }
  return mins + ':' + zeroPad + secsLeft;
};

// opt_secsToStart - number of seconds to start the round (for fast-forwarding
// users to the right state).
Round.prototype.startRound = function(roundNum, opt_secsToStart) {
  this.roundNumEl.text('' + roundNum);
  this.timeLeftTextEl.text('Starts in: ');

  var start = ctrl.table.config.betweenRound;
  if (opt_secsToStart) {
    start = opt_secsToStart;
  }
  this.boardC.getReadyForRound(this.curRound);
  var updateTime = function() {
    this.timeLeftEl.text('' + Round.renderTime(start));
    if (start != 0) {
      window.setTimeout(updateTime, 1000);
    } else {
      this.roundBegins();
    }
    if (!ctrl.STOP_TIMERS) {
      start--;
    }
  }.bind(this);
  updateTime();
};

// Start a timer that runs for 2 minutes before shutting down
// the round and advancing to the next one.
// opt_secsLeft seconds left in round for fast-forwarders.
Round.prototype.startRoundTimer = function(opt_secsLeft) {
  var left = ctrl.table.config.eachRound;
  if (opt_secsLeft) {
    left = opt_secsLeft;
  }
  var roundTimer = function() {
    this.timeLeftEl.text('' + Round.renderTime(left));
    if (left == 0) {
      this.roundOver();
    } else {
      window.setTimeout(roundTimer, 1000);
    }
    if (!ctrl.STOP_TIMERS) {
      left--;
    }
    return false;
  }.bind(this);
  roundTimer();
};

Round.prototype.roundOver = function() {
  this.curRound++;
  this.boardC.startInBetween();
  if (this.curRound == 5) {
  	this.timeLeftTextEl.text('Game over.');
  	this.timeLeftEl.text('');
    multi.sendMessage('gameOver');
  } else {
    this.startRound(this.curRound);
  }
};

BoardC = function(board, solvedWordHandler) {
  this.solvedWordHandler = solvedWordHandler;
  this.board = board;
  // Valid states: 
  //   'IN_BETWEEN' - during the 10 second period of waiting time.
  //   'NEW_BOARD' - new board just shown.
  this.state = '';

  this.curAnswers = {}; // Map of word to true. Only valid words are in map.

  this.solvedWords = {}; // Map of word to true. Only solved words are in this map.
  // Used to show the words which were not found at the end.
  
  // Only need to fetch this once...
  this.allPossible = {}; // Map of word to true. All possible words in English.
};

// Called once when the entire game starts.
BoardC.prototype.init = function() {

  backend.getAllWords(function(words) {
    for (var i = 0; i < words.length; i++) {
      this.allPossible[words[i]] = true;
    }
  }.bind(this));

	this.solvedWordHandler.addDiscoverer({
      word: '<b>Word</b>',
      points: '<b>Pts</b>',
      user: '<b>Discoverer</b>'
  	});
}

// Updates the UI and sets the state to in between.
BoardC.prototype.startInBetween = function() {
  this.state = 'IN_BETWEEN';
  this.updateUi();
};

// Prepare everything for the 10 seconds you have leading up
// to the round starting and save your state.
BoardC.prototype.getReadyForRound = function(curRound) {
  window.console.log('Round about to start in 10 seconds.');
  multi.sendMessage('getRoundInfo', {'r': curRound});
};

// Gets the board from the server, and then asks the server yet
// again to solve this board. Silly, but only way I can get this to
// work....
BoardC.prototype.useSolutions = function(tableInfo, opt_afterCb) {
  // Solve the board and store the results locally for now...
  this.board.resetBoard(tableInfo.getLines());
  this.curAnswers = {};
  var b = new BoardSolver(this.board.asStringToSolve());
  this.curAnswers = {};
  var answersCb = function(answers) {
    // Store the words locally.
    for (var i = 0; i < answers.length; i++) {
      this.curAnswers[answers[i]] = true;
    }
    this.fillSolution();
    if (opt_afterCb) {
      opt_afterCb();  // Call this callback once we've gotten all the solutions.
    }
  }.bind(this);
  b.solve(answersCb);
};

BoardC.prototype.roundStart = function(curRound) {
  this.state = 'NEW_BOARD';
  this.updateUi();
  this.board.destroy();
  // Now we can render the table b/c we are ready with the
  // information...
  this.board.renderBoard();
};

BoardC.prototype.updateUi = function() {
  if (this.state == 'NEW_BOARD') {
  	$('#discovererList').find('tr:gt(0)').remove();
    this.clearDevelConsole();
  } else if (this.state == 'IN_BETWEEN') {
  	// Show the list of words which were not solved.
  	for (var word in this.curAnswers) {
      if (!(word in this.solvedWords)) {
      	// User has not found these words but list them anyway.
      	this.solvedWordHandler.addDiscoverer({
          user: '',
          word: word,
          points: Word.getPoints(word)
      	}, {noScroll: true});
      }
  	}
    this.solvedWordHandler.scrollWords(100); 
  }
};

BoardC.prototype.fillSolution = function() {
  var html = '';
  for (word in this.curAnswers) {
    html += '<div>' + word + '</div>';
  }
  $('#answers').html(html);
};

BoardC.prototype.clearDevelConsole = function() {
  //$('#answers').html('');
};

BoardC.prototype.wordUpdate = function(wordUpdateObj) {
  var msgEl = $('#msgAfterWordEntry');
  var quote = function(val) {
    return "'" + val + "'";
  };
  var wordIsSeen = wordUpdateObj.TotalPoints == -1;
  var word = wordUpdateObj.Word;
  var user = wordUpdateObj.User;
  var totalPoints = wordUpdateObj.TotalPoints;
  if (wordIsSeen) {
    // Already found.
    var span = $('<span>' + word + ' is already found.' + '</span>');
    span.addClass('redText');
    msgEl.append(span);
  } else {
    // If this is my update, notify myself that I got points.
    if (user == ctrl.getUser()) {
      msgEl.html(Word.getPoints(word) + ' points for finding ' + quote(word));
    }
    this.solvedWordHandler.addWord(user, word, totalPoints);
    this.solvedWords[word] = true;
  }
};

BoardC.prototype.submitWord = function(word) {
  var quote = function(val) {
    return "'" + val + "'";
  };
  var msgEl = $('#msgAfterWordEntry');
  var clearMsg = function() {
    msgEl.html('');
  };

  clearMsg();

  var wordIsValid = word in this.curAnswers;
  if (wordIsValid) {
    var params = {
      'word': word,
      'points': Word.getPoints(word),
    };
    multi.sendMessage('submitWord', params);
    // Need to wait to find out whether this submission was valid.
  } else {
  	var wordIsEnglish = word in this.allPossible;
  	if (wordIsEnglish) {
      var span = $('<span>' + quote(word) + ' is not in the puzzle.' + '</span>');
      span.addClass('redText');
      msgEl.append(span);
  	} else {
     var span = $('<span>' + quote(word) + ' is not a word.' + '</span>');
      // Word is not valid!
      span.addClass('redText');
      msgEl.append(span);
    }
  }

  window.setTimeout(clearMsg, 2000);
};

var Word = {};

Word.getPoints = function(word) {
  return word.length * 10;
};

// Keeps track of user leader boards for all users.
UsersHandler = function() {
  this.reset();
}

// Resets the state of the UI to its original state.
UsersHandler.prototype.reset = function() {
  this.usersTable = new UsersTable($('#usersList'));
  this.usersTable.reset();
  this.usersToPoints = {};
  this.usersOrder = []; // A list of users in the order they are currently displayed.
};

/*
UsersHandler.prototype.resetPoints = function() {
  this.usersToPoints = {};
};
*/

UsersTable = function(tableEl) {
  this.tableEl = tableEl;
};

UsersTable.prototype.reset = function() {
  this.tableEl.find('tr:gt(0)').remove();
};

UsersTable.prototype.register = function(user, points) {
  var userList = this.tableEl;
  // Add a row.
  var row = $('<tr></tr>');
  if (ctrl.isMe(user)) {
    row.addClass('me-user-row');
  }
  var td0 = $('<td>' + user + '</td>');
  var td1 = $('<td>' + points + '</td>');
  td1.attr('id', 'userPoints' + user);
  row.append(td0);
  row.append(td1);
  userList.append(row);
};

UsersTable.prototype.hasUser = function(user) {
  var pointsEl = this.tableEl.find('#userPoints' + user);
  if (pointsEl[0]) {
    return true;
  }
  return false;
};

UsersTable.prototype.update = function(user, points) {
  var pointsEl = this.tableEl.find('#userPoints' + user);
  // Don't redraw if we have the same number of points.
  var curPoints = pointsEl.html();
  if (curPoints != points) {
    pointsEl.html(points);
  }
};

UsersTable.prototype.cloneTable = function() {
  return this.tableEl.clone();
};

UsersTable.prototype.get = function() {
  return this.tableEl;
};

// Add a user to the list with a certain number of points.
UsersHandler.prototype.register = function(user, points) {
  this.usersTable.register(user, points);
  this.usersToPoints[user] = points;
  this.usersOrder.push(user);
};

// update('dlluncor', 20) -> david score shows 20 points.
UsersHandler.prototype.update = function(user, points) {
  if (this.usersTable.hasUser(user)) {
    this.usersTable.update(user, points);
    this.usersToPoints[user] = points;
  } else {
    this.register(user, points);
  }
};


var animCompleted = true;

var animOptions = {
    duration: [500,60,300,60,500],
    onComplete: function(){
        animCompleted = true;
        //change update status if required...
    },
    animationSettings: {
        up: {
            left: 0,
            backgroundColor: '#CCFFCC'
        },
        down: {
            left: 0,
            backgroundColor: '#FFCCCC' //the same red as 'down'
        },
        fresh: {
            left: 0,
            backgroundColor: '#CCFFCC' //the same green as 'up'
        },
        drop: {
            left: 0,
            backgroundColor: '#FFCCCC'
        }
    }
};

Dict = function() {
};

// Orders a dictionary by value, order ascending.
Dict.byValue = function(obj) {
  var tuples = [];

  for (var key in obj) tuples.push([key, obj[key]]);

  tuples.sort(function(a, b) {
      a = a[1];
      b = b[1];

      return a < b ? -1 : (a > b ? 1 : 0);
  });
  return tuples;
};

// obj returned from computeNewTable.
UsersHandler.prototype.shouldPerformReorder = function(obj) {
  // Compare current user order with previous user order and make sure there is no
  // difference.
  // Get the current points for the current order and make sure all the numbers
  // are ascending.
  var pointsArr = [];
  for (var i = 0; i < this.usersOrder.length; i++) {
    var user = this.usersOrder[i];
    var points = this.usersToPoints[user];
    pointsArr.push(points);
  }
  // Given the current ordering of our users, do the points decrease for sure???
  var prevVal = null;
  for (var p = 0; p < pointsArr.length; p++) {
    var curVal = pointsArr[p];
    if (prevVal == null) {
      prevVal = curVal;
      continue;
    }
    if (curVal > prevVal) {
      return true;
    }
  }
  return false;
};

UsersHandler.prototype.computeNewTable = function() {
  // Check if users are already in correct order??

  // Create a new ordering for the table based on the current state of
  // the users and their points.
  var newTable = new UsersTable(this.usersTable.cloneTable());
  newTable.reset();

  // Get the users sorted by points.
  var usersAndPoints = Dict.byValue(this.usersToPoints);
  var newOrder = [];
  for (var i = 0; i < usersAndPoints.length; i++) {
    var j = usersAndPoints.length - 1 - i;
    var userAndPoint = usersAndPoints[j];
    var user = userAndPoint[0];
    var points = userAndPoint[1];
    newTable.register(user, points);
    newOrder.push(user);
  }
  return {
    newUsersTable: newTable,
    newUsersOrder: newOrder
  };
};

// Going to rearrange the list of users according to their points.
UsersHandler.prototype.reorder = function() {
  window.console.log('Reordering might take place.');
  if (!animCompleted) {
    window.console.log('Animation is not complete.');
    return;
  }

  var oldTable = this.usersTable.get();
  var obj = this.computeNewTable();
  if (!this.shouldPerformReorder(obj)) {
    window.console.log('Do not need to complete reordering.');
    return;
  }
  var newUsersTable = obj.newUsersTable;
  $(oldTable).rankingTableUpdate(newUsersTable.get(), animOptions);
  animCompleted = false;
  // Then need to replace the
  this.usersTable = newUsersTable;
  this.usersOrder = obj.newUsersOrder;
};

// Handles the update when a new word is found.
WordHandler = function(usersHandler) {
  this.usersHandler = usersHandler;
  this.w = 0; // For scrolling every other.
};

WordHandler.prototype.addDiscoverer = function(inf, opt_options) {
  var aDiv = function(val, width) {
    var div = $('<div>' + val + '</div>');
    div.css('width', width + 'px');
    div.addClass('noOverflow');
    return div;
  };

  var aTd = function(el) {
    var td = $('<td></td>');
    td.append(el);
    return td;
  };
  
  // Draw entry to discoverers board.
  var row = $('<tr></tr>');
  row.append(aTd(aDiv(inf.word, 60)));
  row.append(aTd(aDiv(inf.user, 80)));
  row.append(aTd(aDiv(inf.points, 30)));
  if (inf.user == '') {
  	row.addClass('greyText');
  }
  $('#discovererList').append(row);

  var shouldScroll = this.w % 2 == 0;
  if (opt_options) {
    if ('noScroll' in opt_options) {
      shouldScroll = false;
    }
  }

  if (shouldScroll) {
  // For the scrolling effect. 
    this.scrollWords(40);
  }

  this.w += 1;
};

WordHandler.prototype.scrollWords = function(numPixels) {
  var container = $('.discoverersContainer');
  var curTop = container.scrollTop();
  var nextPos = curTop + numPixels;
  container.scrollTop(nextPos);
};

WordHandler.prototype.addWord = function(user, word, totalPoints) {
  var points = Word.getPoints(word);

  this.addDiscoverer({
  	word: word,
  	points: points,
  	user: user
  });

  // Update the points for the user who scored.
  this.usersHandler.update(user, totalPoints);
};

Console = function() {

};

Console.prototype.init = function() {
  // Hide dev console.
  $('#answers').hide();
  $('#rpcs').hide();
  $('#multiPostRpcs').hide();

  // Handlers.
  $('#showSolutionBtn').click(function(e) {
    $('#answers').toggle();
  });
  $('#showRpcsBtn').click(function(e) {
    $('#rpcs').toggle();
  });
  $('#showMultiRpcsBtn').click(function(e) {
    $('#multiPostRpcs').toggle();
  });

  // Control what gets shown to the user.
  if (ctrl.isLocal()) {
    $('#develConsole').show();
  }
};

Console.prototype.multiPrint = function(str) {
  $('#multiPostRpcs').append('<div>Handle: ' + str + '</div>');
};

var ctrl = {
  STOP_TIMERS: false,
  table: null, // Current table user is part of.
  console: null // devel console to print to.
};

ctrl.stopTimers = function() {
  var text = 'Stop timers';
  if (!ctrl.STOP_TIMERS) {
    text = 'Start timers';
  }
  $('#stopTimerBtn').val(text);
  ctrl.STOP_TIMERS = !ctrl.STOP_TIMERS;
};

// We need to know who THIS user is.
ctrl.getUser = function() {
  return ctrl.table.user;
};

Table = function(curUser, table, token) {
  this.user = curUser; // current user.
  this.table = table; // table name.
  this.token = token; // token to connect with stream.

  this.rounder = null;
  this.usersHandler = null;
  this.boardC = null;
  this.solveWordHandler = null;

  // Game config.
  this.config = {
    betweenRound: 10, // Seconds between rounds.
    eachRound: 90  // Each round is this many seconds.
  };
};

Table.prototype.startGame = function() {
  multi.sendMessage('startGame');
};

Table.prototype.startRound = function() {
    // Now everyone can request the round 1 puzzle and all solvable info
    // needed.
    this.rounder.startRound(1);
};

Table.prototype.startButtonDisabled = function(disabled) {
  $('#startGameBtn').prop('disabled', disabled);
};

Table.prototype.fastForwardUi = function(gameM) {
    // For example we can give the user all known words as well as part
    // of this payload.

    var updateSolvedWords = function() {
      // Show all the words users have found on the right.
      var curWordObjs = gameM.getCurWordObjs();
      for (var wo = 0; wo < curWordObjs.length; wo++) {
        var word = curWordObjs[wo].word; 
        this.solvedWordHandler.addDiscoverer({
          word: word,
          user: 'anOpponent',
          points: Word.getPoints(word)
        });
      }
    }.bind(this);

    // Fast forward to the appropriate round and amount of time left
    // in the round. (server keeps track of when game started and what
    // time it is when the user gets a response?)
    var round = gameM.obj.CurRound;
    var tableInfo = new TableInfo({
      Table: gameM.obj.CurTable
    });

    // Update the round's UI and what round its pointing to.
    this.curRound = round;
    this.rounder.roundNumEl.text('' + round);

    var timeLeft = gameM.getTimeLeft(); // -4; // Relative to a round, how much time do we have left in the race.
    // timeLeft == -4 means that its 4 seconds after round 2 ended (so we are at round 3.)

    if (timeLeft < 0) {
      var curSecsToWait = -timeLeft;
      this.rounder.startRound(round, curSecsToWait);
    }

    var before = new Date();
    var afterCb = function() {
      // Now update the UI once we've gotten all answers the user
      // can engage with the board.
      window.console.log('Show the damn board already!!!');
      var after = new Date();
      var secsDelay = Math.floor((after - before) / 1000.0);
      if (timeLeft > 0) {
        // We can also show the timer finally with the board since
        // we have the board and everything.
        this.rounder.roundBegins(timeLeft - secsDelay);
      }
      updateSolvedWords();
    }.bind(this);
    this.boardC.useSolutions(tableInfo, afterCb);
};

// Updates the UI based on a game model passed from the server.
Table.prototype.updateUi = function(gameM) {
    // Update the UI given the game state.
    this.startButtonDisabled(gameM.isStarted());

    // TODO(dlluncor): This logic is FOOBAR since I'm trying to keep track of 
    // the current position of users and whether we need to reorder things.
    // Left part (users and their total points).
    this.usersHandler.reset();
    var userInfoMap = gameM.getUsersInfo();
    // Need to reset all info we know about the users and completely replace it.
    for (var user in userInfoMap) {
      var userInfo = userInfoMap[user];
      this.usersHandler.register(user, userInfo.points);
    }
};

// Creates and sets up the table with a user name and a table id.
Table.prototype.create = function() {
   $('#entireGameArena').show();

   // Wait for the user to click join table.

  window.console.log("ready for damage");
  this.usersHandler = new UsersHandler();
  this.solvedWordHandler = new WordHandler(this.usersHandler);
	
  // Couple components to this game.
	var board = new Board($('#wordRacerBoard'));
	this.boardC = new BoardC(board, this.solvedWordHandler); // board controller.
  this.rounder = new Round(this.boardC);
  this.rounder.init();

  multi.initConnection(this.user, this.table, this.token);

  $('#startGameBtn').click(function(e) {
    this.startGame();
  }.bind(this));

  var clearWord = function() {
    $('#submissionText').val('');
  };

  $('#clearWordBtn').click(function(e) {
    clearWord();
  });

  var submitWord = function(e) {
    var word = $('#submissionText').val();
    this.boardC.submitWord(word);
    clearWord();
    return false;
  }.bind(this);

  $('#submissionText').keyup(function(e) {
    board.clearPaths();
    if (e.which == 13) {
    	submitWord();
    } else {
    	// Draw the path up until this point in the UI.
    	var word = $('#submissionText').val();
      board.drawPaths(word);
    }
  });
  $('#submitWordBtn').click(function(e) {
    submitWord();
  });

  // Arrange the order of the users based on their points.
  var reorderUsersOnLeft = function() {
    if (!ctrl.STOP_TIMERS) {
      this.usersHandler.reorder();
    }
  }.bind(this);

  window.setInterval(reorderUsersOnLeft, 5000);
}

// Url: localhost:8081/match?hi=cheese&bye=my. qs('hi') -> 'cheese'
function qs(key) {
    key = key.replace(/[*+?^$.\[\]{}()|\\\/]/g, "\\$&"); // escape RegEx meta chars
    var match = location.search.match(new RegExp("[?&]"+key+"=([^&]+)(&|$)"));
    return match && decodeURIComponent(match[1].replace(/\+/g, " "));
};

ctrl.getUserName = function() {
  var user = localStorage.getItem('wr_username');
  if (user == null) {
    user = 'anonymous' + Math.floor(Math.random() * 1000);
  }
  return user;
};

ctrl.isMe = function(userName) {
  return ctrl.table.user == userName; 
};

ctrl.init_ = function() {
    $('#entireGameArena').hide();

    // Must happen first.
    multi.initState();
    ctrl.console = new Console();
    ctrl.console.init();

    var token = $('#userToken').val();

    ctrl.table = new Table(ctrl.getUserName(), qs('t'), token);
    ctrl.table.create();
};

ctrl.isLocal = function() {
  return document.location.hostname == 'localhost';
};

$(document).ready(ctrl.init_);
