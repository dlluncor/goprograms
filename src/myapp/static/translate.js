Translate = function() {

};

Translate.esDict = {
  'user': 'jugador',
  'users': 'jugadores'
};

Translate.dict = {
  'es': Translate.esDict
};

Translate.translate = function(word, outputLang) {
  if (!(outputLang in Translate.dict)) {
    return word;
  }
  wordDict = Translate.dict[outputLang];
  if (!(word in wordDict)) {
  	return word;
  }
  return wordDict[word];
};