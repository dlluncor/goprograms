
// Represents the current board.

Board = function(el) {
  this.boardEl = el;
  this.lines = null;
  this.graph = null; // Graph used to represent these lines.
  this.edgeToEl = {}; // {{row: 0, col: 0}: HTMLDiv}  // for nodes.
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
  	  if (character == ' ') {
  	  	div.addClass('blank');
  	  }
  	  div.addClass('character');
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

  this.graph = new Graph(this.lines);
  // Draw edges between nodes.
  for (var v = 0; v < this.graph.vertices.length; v++) {
  	var vertex = this.graph.vertices[v];
  	for (var e = 0; e < vertex.edges.length; e++) {
  	  var edgeVertex = vertex.edges[e];
      // Should be an edge between these two positions.
      // Since vertices go from left to right and top to bottom
      // We know we can just grab this first vertex and edge
      // and grab a node from there.

      var fromKey = JSON.stringify(vertex.position);
      // Make sure we haven't added this edge already.
      var toKey = JSON.stringify(edgeVertex.position);
      var edgeKey = fromKey + toKey;
      if (edgeKey in this.edgeToEl) {
      	continue;
      }

      var div = this.positionToEl[fromKey];
      var td = div.parent();
      var edgeDiv = $('<div class="line"></div>');
      this.edgeToEl[edgeKey] = edgeDiv;
      this.edgeToEl[toKey + fromKey] = edgeDiv; // Reverse direction will be the same edge.
      // Need to add an edge to this td depending on the direction
      // of this edge connection.
      var edgeType = Graph.direction(vertex, edgeVertex); // 0 - horiz.
      if (edgeType == 'right') {
      	edgeDiv.addClass('horizontal');
      } else if (edgeType == 'down') {
        edgeDiv.addClass('vertical');
      } else if (edgeType == 'diagDown') {
      	edgeDiv.addClass('horizontal');
        edgeDiv.addClass('diag-down-right');
      } else if (edgeType == 'diagUp') {
      	edgeDiv.addClass('horizontal');
      	edgeDiv.addClass('diag-up-right');
      }
      td.append(edgeDiv);
  	}
  }

  this.boardEl.append(table);
};

// lines, e.g., [['a', 'b', 'c'], ['X', 'q', 'd']] that make up
// the board.
Board.prototype.resetBoard = function(lines) {
  this.positionToEl = {};
  this.edgeToEl = {};
  this.lines = lines;
  this.graph = null;
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
  	var pathLen = pathObj.path.length;
  	for (var p = 0; p < pathObj.path.length - 1; p++) {
  	  var position = pathObj.path[p];
  	  var nextPos = pathObj.path[p+1];
  	  this.highlightPosition(position);
  	  this.highlightEdge(position, nextPos);
  	}
  	// No edges emitted from the last character.
  	this.highlightPosition(pathObj.path[pathLen - 1]);
  }
};

Board.prototype.clearPaths = function() {
  // TODO(dlluncor): Possibly slow just find elements with yellowBorder.
  for (var positionStr in this.positionToEl) {
  	this.positionToEl[positionStr].removeClass('yellowBorder');
  }
  for (var keyStr in this.edgeToEl) {
    this.edgeToEl[keyStr].removeClass('yellowLine');
  }
};

Board.prototype.highlightPosition = function(position) {
  var el = this.positionToEl[JSON.stringify(position)];
  el.addClass('yellowBorder');
};

Board.prototype.highlightEdge = function(fromPos, toPos) {
  var key0 = JSON.stringify(fromPos) + JSON.stringify(toPos);
  if (key0 in this.edgeToEl) {
  	this.edgeToEl[key0].addClass('yellowLine');
  }
}

// A graph to represent the board.
Graph = function(lines) {
  this.vertices = this.createGraph_(lines);
};

// Direction between two different vertices, from the fromVertex
// perspective.
// Valid values: ['down', 'right', 'diagDown', 'diagUp'] or ''
// if we don't care about this direction.
Graph.direction = function(from, to) {
  var fromRow = from.position.row;
  var toRow = to.position.row;
  var fromCol = from.position.col;
  var toCol = to.position.col;

  if (fromCol+1 == toCol) {
  	// pointing to the right.
  	if (fromRow == toRow) {
  		return 'right';
  	} else if (fromRow+1 == toRow) {
  		return 'diagDown';
  	}
  } else if (fromCol == toCol) {
  	if (fromRow+1 == toRow) {
  		return 'down';
  	}
  	// don't care about up.
  } else if (fromCol-1 == toCol) {
  	// pointing to the left or diagonal down.
  	if (fromRow+1 == toRow) {
  	  return 'diagUp';
  	}
  }
  return '';
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
  // TODO(dllunor): There is a bug where the paths have a dependency
  // on each other. Try to type a really long word and some letters
  // stay highlighted when they should not be.
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