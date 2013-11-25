
var ctrl = {};


var lounge = function(name, tables) {
  return {
    Name: name,
    Tables: tables
  };
};

var aTable = function(name, users) {
  return {
    Name: name,
    Users: users
  };
};

LoungeList = function(el) {
  this.el = el;
};

LoungeList.prototype.createUsersDiv = function(users) {
  /*
  var hoverIn = function() {

  };
  var hoverOut = function() {
    
  };
  */
  var div = $('<div>&nbsp;(' + users.length + ' users)</div>');
  div.addClass('usersInfo');
  div.attr('title', users.join(','));
  //div.hover(hoverIn, hoverOut);
  return div;
};

LoungeList.prototype.createJoinDiv = function(tableName) {
  var div = $('<div>[ Join ]</div>');
  div.addClass('aJoinBtn');
  div.attr('title', tableName);
  div.click(function(e) {
  	var table = e.currentTarget.getAttribute('title');
    ctrl.joinTableClicked(table);
  });
  return div;
};

LoungeList.prototype.loungesCb = function(loungeArr) {
  
  for (var l = 0; l < loungeArr.length; l++) {
  	var lounge = loungeArr[l];
  	var loungeDiv = $('<div></div>');
  	loungeDiv.attr('id', lounge.Name);
  	loungeDiv.addClass('aLounge');
    var title = $('<h3>' + lounge.Name + '</h3>');
    loungeDiv.append(title);
    for (var t = 0; t < lounge.Games.length; t++) {
      var table = lounge.Games[t];
      var users = [];
      if ('Users' in lounge) {
        users = lounge.Users[t];
      }
      var tableDiv = $('<div></div>'); tableDiv.addClass('aTable');
      var nameDiv = $('<div>' + table + '</div>'); nameDiv.addClass('aTableName');
      var usersDiv = this.createUsersDiv(users);
      var joinDiv = this.createJoinDiv(table);
      tableDiv.append(nameDiv);
      tableDiv.append(usersDiv);
      tableDiv.append(joinDiv);
      loungeDiv.append(tableDiv);
    }
    this.el.append(loungeDiv);
  }
};

ctrl.getLoungesAndTables = function(callback) {
  /*
  var lounge0 = lounge('Intermediate lounge', [
    aTable('Foxy friends', ['ftuser1', 'ftuser2']),
    aTable('Superstars', ['suser1', 'suser2', 'suser3'])
  ]);
  var lounge1 = lounge('Beginner lounge', [
    aTable('Panda pump', ['ppuser1', 'ppuser2']),
    aTable('Giant astronaut', ['gauser1', 'gauser2', 'gauser3'])
  ]);
  var loungeArr = [lounge0, lounge1];
  */
  $.ajax('/getLounges')
      .done(function(data) {
    var loungeArr = JSON.parse(data);
    callback(loungeArr);
  });
};


ctrl.joinTableClicked = function(table) {
  window.console.log('Joining the table.');
  var user = $('#loginUser').val();
  // Open the table in a new tab.
  var url = '/enterTable?t=' + table + '&u=' + user;
  window.open(url,'_blank');
};

ctrl.init_ = function() {

    var ll = new LoungeList($('#loungeList'));
    ctrl.getLoungesAndTables(ll.loungesCb.bind(ll));
};

window.onload = function() {
	ctrl.init_();
};