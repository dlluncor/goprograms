
      Timer = function() {
      };

      Timer.prototype.start = function(count, prefix, callback) {
      	var myCount = count;
        var dec = function() {
          $('#myTime').html(prefix + myCount);
          if (myCount == 0) {
            callback();
          } else {
            myCount--;
            window.setTimeout(dec, 1000);
          }
        };
        dec();
      };

      timer = new Timer();

      Game = function(json) {
        this.obj = json; // JSON which underlies the game object.
      };

      Game.inArr = function(arr, el) {
        for (var i = 0; i < arr.length; i++) {
          if (arr[i] == el) {
          	return true;
          }	
        }
        return false;
      };

      Game.prototype.isStarted = function() {
        return Game.inArr(this.obj.States, 'justStarted');
      };

      var playRound = function() {
      	  var prefix = 'round ' + state.round + ' ends in ';
          timer.start(5, prefix, startRound);
      };

      var startRound = function() {
        // Now everyone can request the round 1 puzzle and all solvable info
        // needed.
        timer.start(2, 'starting in ', playRound);
        if (state.round != 5) {
          sendMessage('getRoundInfo', {'r': state.round});
        } else {
          // Game is over.
          sendMessage('gameOver');
        }
      };

      updateStatus = function(str) {
        $('#gameStatus').append('<div>' + str + '</div>');
      };
      // Logic that deals with responding to user requests.
      handleMessage = function(resp) {
        if (resp.Action == 'join') {
          updateStatus('Set up initial state of the table for the user.');
          var gameM = new Game(resp.Payload);

          // Update the UI given the game state.
          $('#startGame').prop('disabled', gameM.isStarted());
          // For example we can give the user all known words as well as part
          // of this payload.

          // We can also give the user all of the points that everyone has thus
          // far at this snapshot.

          // Fast forward to the appropriate round and amount of time left
          // in the round. (server keeps track of when game started and what
          // time it is when the user gets a response?)
        }
        else if (resp.Action == 'startGame') {
          updateStatus('Game about to start once we generate tables.');
          // This guy can generate the tables for everyone for each round for now
          // even though the backend should be doing that, and then the server
          // can notify everyone when to start the round 1 (everyone should
          // be synchronized at that point).
          var params = {
          	'table1': state.username + 'thisistable1',
          	'table2': state.username + 'thisistable2',
          	'table3': state.username + 'thisistable3',
          	'table4': state.username + 'thisistable4'
          };
          sendMessage('sendTables', params);
        }
        else if (resp.Action == 'startTimers') {
          updateStatus('Starting my timers.');
          startRound();
        }
        else if (resp.Action == 'aboutToStartRound') {
          var info = resp.Payload;
          updateStatus('Got info for round ' + JSON.stringify(info));
          // The timer is still counting down to 10...but now we have all
          // the solutions for the puzzle, what are complete words.
          state.round++;
        }
        //else if (resp.Action == 'endRound') {
          // The round just ended for someone, so we all end the round at the
          // same time and then start our counters at 10.
        //}
        else if (resp.Action == 'wordUpdate') {
          // Here we can handle situations such as:
          // Who is the update for (user - string).
          // What happened - "rejected" means your word got rejected.
          // Points - this is the delta for what the consequence of this action
          // is.
          // Cur points - the points that the user as a result of this action
          // (better than returning deltas in case clients don't get all
          //  the messages.)
          // So we can handle other users giving us updates or we can handle
          // just us getting these updates.
        } else if (resp.Action == 'gameEnded') {
          // Enable the "Start game" button again and provide some message
          // of how the users did.
        }
      };

      var state = {
        table: '',
        token: '',
        username: '',
        round: 1
      };

      sendMessage = function(path, opt_params) {
        path += '?g=' + state.table;
        path += '&t=' + state.token;
        path += '&u=' + state.username;
        if (opt_params) {
          for (var param in opt_params) {
            path += '&' + param + '=' + opt_params[param];	
          }
        }
        var xhr = new XMLHttpRequest();
        xhr.open('POST', path, true);
        xhr.send();
      };

      onOpened = function() {
      	window.console.log('Client channel opened.');
        sendMessage('/opened');
      };
      
      onMessage = function(m) {
        msg = JSON.parse(m.data);
        window.console.log(msg);
        handleMessage(msg);
      }
      
      openChannel = function(token) {
      	state.token = token;
        var channel = new goog.appengine.Channel(token);
        var handler = {
          'onopen': onOpened,
          'onmessage': onMessage,
          'onerror': function() {
          	window.console.log('Error with the channel.');
          },
          'onclose': function() {
          	window.console.log('Channel closing.');
          }
        };
        var socket = channel.open(handler);
        socket.onopen = onOpened;
        socket.onmessage = onMessage;
      }

      init = function() {
        var params = '?';
        params += 'u=' + state.username;
        params += '?g=' + state.table;
        $.ajax('/getToken' + params).done(function(data) {
          openChannel(data);
        });
      };

      var ctrl = {};
      ctrl.init_ = function() {
      	var ind = Math.floor(Math.random() * 100);
        var user = 'sportsguy' + ind;
        $('#username').val(user);

        $('#enterTable').click(function() {
          state.username = $('#username').val();
          state.table = $('#tablename').val();
          init();
        });

        $('#startGame').click(function() {
          sendMessage('startGame');
        });

        $('#curWord').keypress(function(e) {
          var word = $('#curWord').val();
          if (e.which == 13) {
          	sendMessage('submitWord', {
          		'word': word,
          		'points': 20,
          	});
          }
        });
      };
      $(document).ready(ctrl.init_);