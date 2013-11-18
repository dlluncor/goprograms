
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

var ctrl = {};

ctrl.getBoardText = function() {
	return $('#wordRacerBoard').val();
};

ctrl.submitBoard = function() {
	var text = ctrl.getBoardText();
	var b = new Board(text);
	b.solve();
};

ctrl.init_ = function() {
	window.console.log("ready for damage");
};

$(document).ready(ctrl.init_);
