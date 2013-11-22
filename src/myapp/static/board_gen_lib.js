// Board generation. Should probably be ported to Go eventually.

var BoardGen = {};

RndLetter = function() {
};

RndLetter.randLetter = function() {
  var letters = RndLetter.getLetterMix();
  var possib = letters.length-1;
  var ind = Math.floor((Math.random() * possib));
  return letters[ind];
};

RndLetter.emptySpace = 'X';
RndLetter.letters = [
  'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h',
  'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
  'Q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'
];

RndLetter.boosts = null;
RndLetter.letterMix = null;

RndLetter.createBoosts = function() {
  if (RndLetter.boosts) {
  	return RndLetter.boosts;
  }

  var tiers = {
  	 3: 'Q x z',
     6: 'j v y',
     8: 'k w',
    10: 'b c',
    12: 'f g h m n p',
    16: 'd',
    18: 'l r',
    26: 's t u o',
    30: 'a e i'
  };

  // Tiers to boost (prob that we get that letter).
  var boosts = {};
  var numLetters = 0;
  for (var tier in tiers) {
  	var letters = tiers[tier].split(' ');
  	for (var l = 0; l < letters.length; l++) {
      boosts[letters[l]] = tier;
      numLetters++;
  	}
  }
  if (numLetters != 26) {
  	alert('Mapping is wrong needs 26 characters!');
  }
  RndLetter.boosts = boosts;
  return boosts;
};

RndLetter.getLetterMix = function() {
  if (RndLetter.letterMix) {
  	return RndLetter.letterMix;
  }

  var letters = [];
  var boosts = RndLetter.createBoosts();
  // Add letters with boosts that many extra times.
  for (var i = 0; i < RndLetter.letters.length; i++) {
  	var letter = RndLetter.letters[i];
  	var repeat = 1;
  	if (letter in boosts) {
  		repeat = boosts[letter];
  	}

  	// Now put that letter into the mix that many times.
  	for (var j = 0; j < repeat; j++) {
  		letters.push(letter);
  	}
  }
  RndLetter.letterMix = letters;
  return RndLetter.letterMix;
};

BoardGen.createLine = function(emptyIndices, width) {
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
BoardGen.generateBoard = function(curRound) {
  window.console.log('About to generate the board.');
  var lines = [];
  if (curRound == 1) {
  	// Need a 4 by 4 grid no empty spaces.
  	var width = 4;
  	lines = [
      BoardGen.createLine([], width),
      BoardGen.createLine([], width),
      BoardGen.createLine([], width),
      BoardGen.createLine([], width)
  	];
  } else if (curRound == 2) {
  	var width = 6;
  	lines = [
  	  BoardGen.createLine([0, 1, 4, 5], width),
  	  BoardGen.createLine([0, 5], width),
  	  BoardGen.createLine([], width),
  	  BoardGen.createLine([], width),
  	  BoardGen.createLine([0, 5], width),
  	  BoardGen.createLine([0, 1, 4, 5], width)
  	];
  } else if (curRound == 3) {
  	var width = 6;
    lines = [
  	  BoardGen.createLine([4, 5], width),
  	  BoardGen.createLine([4, 5], width),
  	  BoardGen.createLine([], width),
  	  BoardGen.createLine([], width),
  	  BoardGen.createLine([0, 1], width),
  	  BoardGen.createLine([0, 1], width)
  	];
  } else if (curRound == 4) {
  	var width = 6;
    lines = [
  	  BoardGen.createLine([], width),
  	  BoardGen.createLine([], width),
  	  BoardGen.createLine([2, 3], width),
  	  BoardGen.createLine([2, 3], width),
  	  BoardGen.createLine([], width),
  	  BoardGen.createLine([], width)
  	];
  }
  return lines;
};

// Join lines for a table, what the server expects as one string
// [['a', 'b'], ['X', 'b']] -> 'ab\nXb'
BoardGen.join = function(lines) {
  var charLines = [];
  for (var i = 0; i < lines.length; i++) {
    charLines.push(lines[i].join(''));
  }
  return charLines.join('\n');
};

// Reverse of join. From the server we get one string and need to
// conver it back to an array of strings.
BoardGen.unjoin = function(str) {
  var lines = [];
  var lineStrs = str.split('\n');
  for (var i = 0; i < lineStrs.length; i++) {
  	var lineStr = lineStrs[i];
  	var line = [];
  	for (var j = 0; j < lineStr.length; j++) {
      line.push(lineStr[j]);
    }
    lines.push(line);
  }
  return lines;
};