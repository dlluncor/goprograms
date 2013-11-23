
var ctrl = {};

ctrl.init_ = function() {
	$('#joinTableBtn').click(function(e) {
	  window.console.log('Joining the table.');
	  var user = $('#loginUser').val();
	  var table = $('#tableId').val();
	  // Open the table in a new tab.
	  var url = '/enterTable?t=' + table + '&u=' + user;
	  window.open(url,'_blank');

	});
};


$(document).ready(ctrl.init_);