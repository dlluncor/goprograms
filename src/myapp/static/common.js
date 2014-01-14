// Url: localhost:8081/match?hi=cheese&bye=my. qs('hi') -> 'cheese'
function qs(key) {
    key = key.replace(/[*+?^$.\[\]{}()|\\\/]/g, "\\$&"); // escape RegEx meta chars
    var match = location.search.match(new RegExp("[?&]"+key+"=([^&]+)(&|$)"));
    return match && decodeURIComponent(match[1].replace(/\+/g, " "));
};

var ctrl = ctrl || {};

ctrl.setUserName = function(userName) {
  createCookie('ww-user', userName, 1);
};

ctrl.getUserName = function(shouldSave) {
  var saveUserName = function(userName) {
    return ctrl.setUserName(userName);
  };

  var debugUser = qs('debugUser');
  var user = '';
  if (debugUser) {
    // Debug user gets precedence.
    // Save the user to the cookie.
    user = debugUser;
  } else {
    user = getCookie('ww-user');
    if (!user) {
      user = 'anonymous' + Math.floor(Math.random() * 1000);
    }
  }
  if (shouldSave) {
    saveUserName(user);
  }
  return user;
};

/* Cookie stuff. */
function createCookie(name, value, days) {
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        var expires = "; expires=" + date.toGMTString();
    }
    else var expires = "";
    document.cookie = name + "=" + value + expires + "; path=/";
}
function getCookie(c_name) {
    if (document.cookie.length > 0) {
        c_start = document.cookie.indexOf(c_name + "=");
        if (c_start != -1) {
            c_start = c_start + c_name.length + 1;
            c_end = document.cookie.indexOf(";", c_start);
            if (c_end == -1) {
                c_end = document.cookie.length;
            }
            return unescape(document.cookie.substring(c_start, c_end));
        }
    }
    return "";
}