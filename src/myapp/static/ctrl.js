
Board = function(text) {
	this.text = text;
};

Board.prototype.printSolution = function(linesAsText) {
	var words = linesAsText.split(',');
	var html = '';
	for (var i = words.length -1; i >= 0; i--) {
		var word = words[i];
		html += '<div>' + word + '</div>';
	}
	$('#answers').html(html);
};

Board.prototype.solve = function() {
	var text = this.text;
	var lines = text.split('\n');
	// Validate the board works.
	var length = lines[0].length;
	window.console.log(this.text);
	var that = this;
	$.ajax('/wordracer_json?board=' + text + '&length=' + length)
		.done(function(data) {
				that.printSolution(data);
				window.console.log('success');
				window.console.log(data);
				});    
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

Round.prototype.start = function() {
  this.curRound = 1;
  this.startRound(this.curRound);
};

Round.prototype.startRound = function(roundNum) {
  this.roundNumEl.text('' + roundNum);
  this.timeLeftTextEl.text('Starts in: ');

  var start = 3;
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
    start--;
  }.bind(this);
  updateTime();
};

// Start a timer that runs for 2 minutes before shutting down
// the round and advancing to the next one.
Round.prototype.startRoundTimer = function() {
  var left = 5;
  var roundTimer = function() {
    this.timeLeftEl.text('' + left);
    if (left == 0) {
      this.roundOver();
    } else {
      window.setTimeout(roundTimer, 1000);
    }
    left--;
  }.bind(this);
  roundTimer();
};

Round.prototype.roundOver = function() {
  this.curRound++;
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

BoardC = function(boardEl) {
  this.boardEl = boardEl;
  this.curBoard = null; // Cur board as a string joined by \n lines.
};

BoardC.prototype.roundStart = function(curRound) {
  window.console.log('Round about to start.');
  this.generateBoard(curRound);
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

  this.boardEl.html('');

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
  this.curBoard = boardTextLines.join('\n');
};

BoardC.prototype.solveBoard = function() {
	var b = new Board(this.curBoard);
	b.solve();
};

var ctrl = {};

ctrl.init_ = function() {
    window.console.log("ready for damage");
	// Couple components to this game.
	var boardC = new BoardC($('#wordRacerBoard')); // board controller.
    var rounder = new Round(boardC);
    rounder.start();

    // Handlers.
    $('#solveBoardBtn').click(function(e) {
      boardC.solveBoard();
    });
};

$(document).ready(ctrl.init_);
