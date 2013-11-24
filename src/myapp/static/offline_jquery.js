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