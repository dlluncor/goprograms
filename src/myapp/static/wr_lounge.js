
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
  var div = $('<div>Users</div>');
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
    for (var t = 0; t < lounge.Tables.length; t++) {
      var table = lounge.Tables[t];
      var tableDiv = $('<div></div>'); tableDiv.addClass('aTable');
      var nameDiv = $('<div>' + table.Name + '</div>'); nameDiv.addClass('aTableName');
      var usersDiv = this.createUsersDiv(table.Users);
      var joinDiv = this.createJoinDiv(table.Name);
      tableDiv.append(nameDiv);
      tableDiv.append(usersDiv);
      tableDiv.append(joinDiv);
      loungeDiv.append(tableDiv);
    }
    this.el.append(loungeDiv);
  }
};

ctrl.getLoungesAndTables = function(callback) {
  var lounge0 = lounge('Intermediate lounge', [
    aTable('Foxy friends', ['ftuser1', 'ftuser2']),
    aTable('Superstars', ['suser1', 'suser2', 'suser3'])
  ]);
  var lounge1 = lounge('Beginner lounge', [
    aTable('Panda pump', ['ppuser1', 'ppuser2']),
    aTable('Giant astronaut', ['gauser1', 'gauser2', 'gauser3'])
  ]);
  var loungeArr = [lounge0, lounge1];
  callback(loungeArr);
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

El = function(el) {
  this.el = el;
};

$ = function(text) {
  // Enable wrapping just an element if its not a text. How to ensure that an object
  // is of class element?

  var byId = text.indexOf('#') == 0;
  if (byId) {
  	return new El(document.getElementById(text.substring(1)));
  }
  var elName = text.substring(text.lastIndexOf('<')+2, text.length-1);
  var el = document.createElement(elName);

  // Get inner html.
  // TODO(dlluncor): Add the ability to have <div class="cheese"> like JQuery
  // enables.
  var l0 = text.indexOf('>');
  var l1 = text.lastIndexOf('<');
  var innerText = text.substring(l0+1, l1);
  el.innerHTML = innerText;
  return new El(el);
};

El.prototype.addClass = function(className) {
  this.el.classList.add(className);
};

El.prototype.append = function(node) {
  this.el.appendChild(node.el);
};

El.prototype.attr = function(attrName, opt_attrVal) {
  if (opt_attrVal) {
  	this.el.setAttribute(attrName, opt_attrVal);
  }
};

El.prototype.click = function(callback) {
  this.el.onclick = callback;
};

El.prototype.val = function() {
  return this.el.value;
};

window.onload = function() {
	ctrl.init_();
};