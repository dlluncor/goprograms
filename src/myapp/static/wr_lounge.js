
var ctrl = {};


ctrl.TABLES = {
  'game0': 'Foxy Friday',
  'game1': 'Sloppy Joes',
  'game2': 'Cajun Slide',
  'game3': 'Chocolate Thunder',
  'game4': 'Maple Breeze',
  'game5': 'China Force'
};

var rename = function(table) {
  if (table in ctrl.TABLES) {
    return ctrl.TABLES[table];
  }
  return table;
};

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
  var suffix = 'users';
  if (users.length == 1) {
  	suffix = 'user';
  }

  var div = $('<div>&nbsp;(' + users.length + ' ' + suffix +')</div>');
  div.addClass('usersInfo');
  div.attr('title', users.join(','));
  //div.hover(hoverIn, hoverOut);
  return div;
};

LoungeList.prototype.createJoinDiv = function(tableName) {
  var div = $('<div>[ Join table ]</div>');
  div.addClass('aJoinBtn');
  div.attr('theTableName', tableName);
  div.click(function(e) {
  	var table = e.currentTarget.getAttribute('theTableName');
    ctrl.joinTableClicked(table);
  });
  return div;
};

LoungeList.prototype.loungesCb = function(loungeResp) {
  
  var gameMap = loungeResp.GameInfo;
  var loungeArr = loungeResp.Lounges;
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
      if (table in gameMap) {
      	var tableInfo = gameMap[table];
        users = tableInfo.Users;
      }
      var tableDiv = $('<div></div>'); tableDiv.addClass('aTable');
      var nameDiv = $('<div>' + rename(table) + '</div>'); nameDiv.addClass('aTableName');
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
    var loungeResp = JSON.parse(data);
    window.console.log(loungeResp);
    callback(loungeResp);
  });
};


ctrl.joinTableClicked = function(table) {
  window.console.log('Joining the table.');
  var user = $('#loginUser').html();
  // Open the table in a new tab.
  var url = '/enterTable?t=' + table + '&u=' + user;
  window.open(url,'_blank');
};

ctrl.init_ = function() {
    $('#loginUser').html(localStorage.getItem('wr_username'));
    var ll = new LoungeList($('#loungeList'));
    ctrl.getLoungesAndTables(ll.loungesCb.bind(ll));
};

window.onload = function() {
	ctrl.init_();
};