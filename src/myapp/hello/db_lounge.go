package hello

// TODO(dlluncor): Need ability to create lounges too...


type MyLounge struct {
  Games []string
  Users []string 
}

// Database for the entire lounge.
type changeLoungeFunc func(l *MyLounge) bool


func defaultLounge() *MyLounge {
  return &MyLounge{
    Games: []string{},
    Users: []string{},
  }
}

// TODO(dlluncor): Merge this with ChangeGame as they are the same except for the key
// and the type of the changeEntityFunc, and the defaultLounge() thing.
// Utility function for reading from and updating a game before then
// doing further processing.
func ChangeLounge(c appengine.Context, loungeId string, clf changeLoungeFunc) *MyLounge {
  // Store all tables as part of the game state and send
  // a "startTimer" response.
  loungeKey := loungeId
  l := defaultLounge()
  err := datastore.RunInTransaction(c, func(c appengine.Context) error {
    k := datastore.NewKey(c, "WrLounge", loungeKey, 0, nil)
    if err := datastore.Get(c, k, l); err != nil {
      return err;
    }

    // Perform special logic here to manipulate game.
    shouldRunUpdate := clf(l)
    if !shouldRunUpdate {
      // Sometimes we don't need to update the database so don't.
      return nil
    }

    if _, err := datastore.Put(c, k, l); err != nil {
      return err
    }
    return nil
  }, nil)

  if err != nil {
    c.Errorf("Err in db transaction %v", err)
  }
  return l
}