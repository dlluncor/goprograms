/* 
 * This file was copied from
 * https://github.com/just-another-one-timmy/injackted/blob/master/src/structs/set.go
 * but hey, I asked myself for permission and gave it to myself, so everything is fine.
 */

/*
 Package provides basic data structures to be used in the client program.
*/

package hello

/*
 Set data structure.
*/
type Set struct {
        internal map[interface{}]struct{}
        size     int
}

// Creates new empty set.
func NewSet() *Set {
        return &Set{make(map[interface{}]struct{}), 0}
}

// Returns size of the set.
func (s *Set) Size() int {
        return s.size
}

// Returns true if the set is empty.
func (s *Set) IsEmpty() bool {
        return s.Size() == 0
}

// Adds item to the set.
func (s *Set) Add(item interface{}) {
        if !s.ItemPresent(item) {
                s.internal[item] = struct{}{}
                s.size += 1
        }
}

// Checks if item is present in the set.
func (s *Set) ItemPresent(item interface{}) bool {
        _, present := s.internal[item]
        return present
}

// Removes item from the set.
// Works even when item is not present.
func (s *Set) Remove(item interface{}) {
        if s.ItemPresent(item) {
                delete(s.internal, item)
                s.size -= 1
        }
}

func (s *Set) Equals(otherSet *Set) bool {
        if s.Size() != otherSet.Size() {
                return false
        }
        for item, _ := range s.internal {
                if !otherSet.ItemPresent(item) {
                        return false
                }
        }
        for item, _ := range otherSet.internal {
                if !s.ItemPresent(item) {
                        return false
                }
        }
        return true
}

// Intersects two sets, returns new set as a result.
// Doesn't modify arguments.
func (s *Set) Intersection(otherSet *Set) *Set {
        result := NewSet()

        for item, _ := range s.internal {
                if otherSet.ItemPresent(item) {
                        result.Add(item)
                }
        }

        return result
}

// Unions two sets, returns new set as a result.
// Doesn't modify arguments.
func (s *Set) Union(otherSet *Set) *Set {
        result := NewSet()

        for item, _ := range s.internal {
                result.Add(item)
        }
        for item, _ := range otherSet.internal {
                result.Add(item)
        }

        return result
}

// Finds a difference set - otherSet, returns new set as a result.
// Doesn't modify arguments.
func (s *Set) Difference(otherSet *Set) *Set {
        result := NewSet()

        for item, _ := range s.internal {
                if !otherSet.ItemPresent(item) {
                        result.Add(item)
                }
        }

        return result
}

// Iterator over items.
func (s *Set) Iterator() <-chan interface{} {
        iteratorChannel := make(chan interface{})
        go func() {
                for item, _ := range s.internal {
                        iteratorChannel <- item
                }
                close(iteratorChannel)
        }()
        return iteratorChannel
}