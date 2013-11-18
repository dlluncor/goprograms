
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

  // Game config.
  this.config = {
  	betweenRound: 10, // Seconds between rounds.
    eachRound: 10  // Each round is this many seconds.
  };
};

Round.prototype.start = function() {
  this.boardC.start();
  this.curRound = 1;
  this.startRound(this.curRound);
};

Round.prototype.startRound = function(roundNum) {
  this.roundNumEl.text('' + roundNum);
  this.timeLeftTextEl.text('Starts in: ');

  var start = this.config.betweenRound;
  this.boardC.getReadyForRound(this.curRound);
  var updateTime = function() {
    this.timeLeftEl.text('' + start);
    if (start != 0) {
      window.setTimeout(updateTime, 1000);
    } else {
      this.timeLeftTextEl.text('Time: ');
      this.startRoundTimer();
      window.console.log('Go start game!');
      this.boardC.roundStart(this.curRound);
    }
    if (!ctrl.STOP_TIMERS) {
      start--;
    }
  }.bind(this);
  updateTime();
};

// Start a timer that runs for 2 minutes before shutting down
// the round and advancing to the next one.
Round.prototype.startRoundTimer = function() {
  var left = this.config.eachRound;
  var roundTimer = function() {
    this.timeLeftEl.text('' + left);
    if (left == 0) {
      this.roundOver();
    } else {
      window.setTimeout(roundTimer, 1000);
    }
    if (!ctrl.STOP_TIMERS) {
      left--;
    }
  }.bind(this);
  roundTimer();
};

Round.prototype.roundOver = function() {
  this.curRound++;
  this.boardC.startInBetween();
  if (this.curRound == 5) {
  	this.timeLeftTextEl.text('Game over.');
  	this.timeLeftEl.text('');
  } else {
    this.startRound(this.curRound);
  }
};

RndLetter = function() {
};

RndLetter.randLetter = function() {
  var ind = Math.floor((Math.random() * 25));
  return RndLetter.letters[ind];
};

RndLetter.emptySpace = 'X';
RndLetter.letters = [
  'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h',
  'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
  'Q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'
 ];

BoardC = function(boardEl, solvedWordHandler) {
  this.solvedWordHandler = solvedWordHandler;
  this.boardEl = boardEl;
  // Valid states: 
  //   'IN_BETWEEN' - during the 10 second period of waiting time.
  //   'NEW_BOARD' - new board just shown.
  this.state = '';

  this.curBoard = null; // 'ab\nXc' Cur board as a string joined by \n lines.
  this.lines = null; // [['a', 'b'], ['X', 'c']]An array of the board where each letter is an element.
  this.curAnswers = {}; // Map of word to true. Only valid words are in map.

  // Lots of this stuff will have to be distributed...
  this.solvedWords = {}; // Map of word to true. Only solved words are in this map.
  
  // Only need to fetch this once...
  this.allPossible = {}; // Map of word to true. All possible words in English.
};

BoardC.prototype.createLine = function(emptyIndices, width) {
  var emptyMap = {};
  emptyIndices.forEach(function(index) {
    emptyMap[index] = true;
  });

  var letters = [];
  for (var i = 0; i < width; i++) {
    if (i in emptyMap) {
    	letters.push(RndLetter.emptySpace);
    } else {
      letters.push(RndLetter.randLetter());
    }
  }
  return letters;
};

// Returns the board as an array of characters, e.g.
// [['a', 'b', 'X'], ['f', 'e', 'd']]
BoardC.prototype.generateBoard = function(curRound) {
  window.console.log('About to generate the board.');
  var lines = [];
  if (curRound == 1) {
  	// Need a 4 by 4 grid no empty spaces.
  	var width = 4;
  	lines = [
      this.createLine([], width),
      this.createLine([], width),
      this.createLine([], width),
      this.createLine([], width)
  	];
  } else if (curRound == 2) {
  	var width = 6;
  	lines = [
  	  this.createLine([0, 1, 4, 5], width),
  	  this.createLine([0, 5], width),
  	  this.createLine([], width),
  	  this.createLine([], width),
  	  this.createLine([0, 5], width),
  	  this.createLine([0, 1, 4, 5], width)
  	];
  } else if (curRound == 3) {
  	var width = 6;
    lines = [
  	  this.createLine([4, 5], width),
  	  this.createLine([4, 5], width),
  	  this.createLine([], width),
  	  this.createLine([], width),
  	  this.createLine([0, 1], width),
  	  this.createLine([0, 1], width)
  	];
  } else if (curRound == 4) {
  	var width = 6;
    lines = [
  	  this.createLine([], width),
  	  this.createLine([], width),
  	  this.createLine([2, 3], width),
  	  this.createLine([2, 3], width),
  	  this.createLine([], width),
  	  this.createLine([], width)
  	];
  }

  // Create board and write it to the display.
  var boardTextLines = [];
  for (var j = 0; j < lines.length; j++) {
  	var lineArr = lines[j];
  	var lineText = lineArr.join('');
    boardTextLines.push(lineText);
  }
  this.curBoard = boardTextLines.join('\n');
  return lines;
};

BoardC.prototype.renderBoard = function(lines) {
  // Create board and write it to the display.
  var boardTextLines = [];
  var table = $('<table class="boardTable"></table>');
  for (var j = 0; j < lines.length; j++) {
  	var lineArr = lines[j];
  	var lineText = lineArr.join('');
  	var row = $('<tr></tr>');
  	for (var c = 0; c < lineArr.length; c++) {
  	  var character = lineArr[c];
  	  if (character == RndLetter.emptySpace) {
  	  	character = ' ';
  	  }
      var td = $('<td><div>' + character + '</div></td>');
      row.append(td);
  	}
  	table.append(row);
    boardTextLines.push(lineText);
  }
  this.boardEl.append(table);
};

// Called once when the entire game starts.
BoardC.prototype.start = function() {

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
  var lines = this.generateBoard(curRound);  // renders as well.
  // Solve the board and store the results locally for now...
  var b = new BoardSolver(this.curBoard);
  this.curAnswers = {};
  this.lines = lines;
  var answersCb = function(answers) {
    // Store the words locally.
    for (var i = 0; i < answers.length; i++) {
      this.curAnswers[answers[i]] = true;
    }
    this.fillSolution();
  }.bind(this);
  b.solve(answersCb);
};

BoardC.prototype.roundStart = function(curRound) {
  this.state = 'NEW_BOARD';
  this.updateUi();
  this.boardEl.html('');
  // Now we can render the table b/c we are ready with the
  // information...
  this.renderBoard(this.lines);
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
      	});
      }
  	} 
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
    var wordIsSeen = word in this.solvedWords;
    // Distributed.
    if (wordIsSeen) {
      // Already found.
      msgEl.html(word + ' is already found.');
    } else {
      // Give this guy some points...
      msgEl.html(Word.getPoints(word) + 'points for finding ' + quote(word));
      this.solvedWords[word] = true;
      this.solvedWordHandler.addWord(word);
    }
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

// Keeps track of user leader boards and this current user.
UsersHandler = function(curUser) {
  this.curUser = curUser;
  this.register('<b>User</b>', '<b>Points</b>');
}

// Add a user to the list with zero points.
UsersHandler.prototype.register = function(user, points) {
  var userList = $('#usersList');

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

  // Add a row.
  var row = $('<tr></tr>');
  var td0 = $('<td>' + user + '</td>');
  row.append(aTd(aDiv(user, 90)));
  var div1 = aDiv(points, 50);
  div1.attr('id', 'userPoints' + user);
  row.append(aTd(div1));
  userList.append(row);
};

// update('dlluncor', 20) -> david get's twenty more points.
UsersHandler.prototype.update = function(user, points) {
  var pointsEl = $('#userPoints' + user);
  var curPoints = parseInt(pointsEl.html());
  pointsEl.html(curPoints + points);
};

// Handles the update when a new word is found.
WordHandler = function(usersHandler) {
  this.usersHandler = usersHandler;
};

WordHandler.prototype.addDiscoverer = function(inf) {
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
};

WordHandler.prototype.addWord = function(word) {
  var points = Word.getPoints(word);
  var user = this.usersHandler.curUser;

  this.addDiscoverer({
  	word: word,
  	points: points,
  	user: user
  });

  // Update the points for the user who scored.
  this.usersHandler.update(user, points);
};

var ctrl = {
  STOP_TIMERS: false
};

ctrl.stopTimers = function() {
  var text = 'Stop timers';
  if (!ctrl.STOP_TIMERS) {
    text = 'Start timers';
  }
  $('#stopTimerBtn').val(text);
  ctrl.STOP_TIMERS = !ctrl.STOP_TIMERS;
};

ctrl.init_ = function() {
    window.console.log("ready for damage");
    var curUser = 'sportsguy560';
    var usersHandler = new UsersHandler(curUser);
    usersHandler.register(curUser, 0);
    var solvedWordHandler = new WordHandler(usersHandler);
	// Couple components to this game.
	var boardC = new BoardC(
      $('#wordRacerBoard'), solvedWordHandler); // board controller.
    var rounder = new Round(boardC);
    rounder.start();

    // Hide dev console.
    $('#answers').hide();
    $('#rpcs').hide();

    // Handlers.
    $('#showSolutionBtn').click(function(e) {
      $('#answers').toggle();
    });
    $('#showRpcsBtn').click(function(e) {
      $('#rpcs').toggle();
    });

    var clearWord = function() {
      $('#submissionText').val('');
    };

    $('#clearWordBtn').click(function(e) {
      clearWord();
    });

    var submitWord = function() {
      var word = $('#submissionText').val();
      boardC.submitWord(word);
      clearWord();
    };

    $('#submissionText').keypress(function(e) {
      if (e.which == 13) {
      	submitWord();
      }
    });
    $('#submitWordBtn').click(function(e) {
      submitWord();
    });
};

$(document).ready(ctrl.init_);
