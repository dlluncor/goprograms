Translate = function() {

};

Translate.esDict = {
  'user': 'jugador',
  'users': 'jugadores',
  'Start game': 'Empieza',
  'Word': 'Palabra',
  'Discoverer': 'Persona',
  'User': 'Jugador',
  'Points': 'Puntos',
  'Round': 'Ronda',
  'Time': 'Tiempo',
  'Submit': 'Entra',
  'Clear': 'Borra',
  'Type word': 'Escribe'
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