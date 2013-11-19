
// Represents the current board.

Board = function(el) {
  this.boardEl = el;
  this.lines = null;
  this.graph = null; // Graph used to represent these lines.
  this.positionToEl = {}; // {{row: 0, col: 0}: HTMLDiv}
};

Board.prototype.renderBoard = function() {
  // Create board and write it to the display.
  var table = $('<table class="boardTable"></table>');
  for (var j = 0; j < this.lines.length; j++) {
  	var lineArr = this.lines[j];
  	var lineText = lineArr.join('');
  	var row = $('<tr></tr>');
  	for (var c = 0; c < lineArr.length; c++) {
  	  var character = lineArr[c];
  	  if (character == RndLetter.emptySpace) {
  	  	character = ' ';
  	  }
  	  var div = $('<div>' + character + '</div>');
  	  var position = {
        row: j,
        col: c
  	  };
  	  this.positionToEl[JSON.stringify(position)] = div;
      var td = $('<td></td>');
      td.append(div);
      row.append(td);
  	}
  	table.append(row);
  }
  this.boardEl.append(table);
};

// lines, e.g., [['a', 'b', 'c'], ['X', 'q', 'd']] that make up
// the board.
Board.prototype.resetBoard = function(lines) {
  this.positionToEl = {};
  this.lines = lines;
  this.graph = new Graph(this.lines);
};

Board.prototype.destroy = function() {
  this.boardEl.html('');
};

// Returns the board as a string, each line is \n.
Board.prototype.asStringToSolve = function() {
  // Create board and write it to the display.
  var boardTextLines = [];
  for (var j = 0; j < this.lines.length; j++) {
  	var lineArr = this.lines[j];
  	var lineText = lineArr.join('');
    boardTextLines.push(lineText);
  }
  return boardTextLines.join('\n');
};

Board.prototype.drawPaths = function(word) {
  var characters = word.split('');
  // position {{
  //   row: number,
  //   col: number
  // }}
  // pathObj {{
  //   path: !Array.<position>
  // }}
  window.console.log('Finding path for characters: ' + characters.join(''));
  var pathObjs = this.graph.findPaths(characters);
  for (var i = 0; i < pathObjs.length; i++) {
  	var pathObj = pathObjs[i];
  	for (var p = 0; p < pathObj.path.length; p++) {
  	  var position = pathObj.path[p];
  	  this.highlightPosition(position);
  	}
  }
};

Board.prototype.clearPaths = function() {
  for (var positionStr in this.positionToEl) {
  	this.positionToEl[positionStr].css('border', '2px solid black');
  }
};

Board.prototype.highlightPosition = function(position) {
  var el = this.positionToEl[JSON.stringify(position)];
  el.css('border', '2px solid yellow');
};

// A graph to represent the board.
Graph = function(lines) {
  this.vertices = this.createGraph_(lines);
};


Graph.prototype.accruePaths = function(
	curChars, curVertices, curPath, listOfPaths,
	seenVertices) {
  if (curChars.length == 0) {
    return;
  }

  // If any of the vertices match we can continue a path there.
  for (var v = 0; v < curVertices.length; v++) {
    var vertex = curVertices[v];
    // Can't continue if we've seen this vertex.
    var vAsStr = JSON.stringify(vertex.position);
    if (vAsStr in seenVertices) {
      continue;
    }
    if (vertex.data == curChars[0]) {
      // Add to path and follow all edge vertices.
      var nextChars = curChars.slice(1);
      var nextPath = curPath.slice(0); // array copy.
      nextPath.push(vertex.position);
      listOfPaths.push(nextPath);
      seenVertices[vAsStr] = true;
      this.accruePaths(nextChars, vertex.edges, nextPath, listOfPaths,
        seenVertices);
      delete seenVertices[vAsStr]; // Unvisit when done.
  	}
  }
};

Graph.prototype.findPaths = function(characters) {
  // TODO(dlluncor): Keeping no state about previous paths seen so this
  // could be inefficient if we don't maintain some delay on keyup.
  var listOfPaths = [];
  for (var v = 0; v < this.vertices.length; v++) {
  	var vertex = this.vertices[v];
  	var curPath = [];
  	var seenVertices = {};
  	this.accruePaths(characters, [vertex], curPath, listOfPaths,
        seenVertices);
  }

  var pathObjs = [];
  for (var p = 0; p < listOfPaths.length; p++) {
  	var path = listOfPaths[p];
  	if (path.length != characters.length) {
  		// Only found a partial path.
  		continue;
  	}
  	pathObjs.push({
  	  path: path
  	});
  }
  return pathObjs;
};

// Creates a list of vertices given the board representation of lines.
Graph.prototype.createGraph_ = function(lines) {

  // Iterate through this board using the same indexing scheme for
  // position purposes.
  var vertices = [];
  var vertexMap = {}; // Map of position to vertex object.
  for (var row = 0; row < lines.length; row++) {
  	var line = lines[row];
  	for (var col = 0; col < line.length; col++) {
      var character = line[col];
      if (character == RndLetter.emptySpace) {
      	// Empty space cannot be highlighted.
      	continue;
      }
      var position = {
          row: row,
          col: col
      };
      var vertex = {
        position: position,
        data: character,
        edges: []
      };
      vertexMap[JSON.stringify(position)] = vertex;
      vertices.push(vertex);
  	}
  }

  var numCols = lines[0].length;  // Width of the board in chars.
  var numRows = lines.length;
  // Now need to form edges in the graph based on position.
  for (var positionStr in vertexMap) {
  	var pos = JSON.parse(positionStr);
  	var vertex = vertexMap[positionStr];

    // Explore if there are neighors to add.
    for (var rowMove = -1; rowMove < 2; rowMove++) {
      for (var colMove = -1; colMove < 2; colMove++) {
      	if (rowMove == 0 && colMove == 0) {
      		// That's just yourself.
      		continue;
      	}
        var neighPos = {
          row: pos.row + rowMove,
          col: pos.col + colMove
        };
        var neighKey = JSON.stringify(neighPos);
        if (neighKey in vertexMap) {
          // We've found a neighbor.
          vertex.edges.push(vertexMap[neighKey]);
        }
      }
    }
  }
  return vertices;
};