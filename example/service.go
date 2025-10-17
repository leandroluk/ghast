// example/service.go
package example

import "fmt"

// UserService handles in-memory user storage.
type UserService struct {
	store map[int]*User
	next  int
}

func (s *UserService) OnInit() {
	s.store = make(map[int]*User)
	s.next = 1
	fmt.Println("[user.service] initialized")
}

func (s *UserService) Create(name, email string) *User {
	u := &User{ID: s.next, Name: name, Email: email}
	s.store[s.next] = u
	s.next++
	return u
}

func (s *UserService) FindAll() []*User {
	users := []*User{}
	for _, u := range s.store {
		users = append(users, u)
	}
	return users
}

func (s *UserService) FindByID(id int) *User {
	return s.store[id]
}

func (s *UserService) Update(id int, name, email string) *User {
	if u, ok := s.store[id]; ok {
		u.Name = name
		u.Email = email
		return u
	}
	return nil
}

func (s *UserService) Delete(id int) bool {
	if _, ok := s.store[id]; ok {
		delete(s.store, id)
		return true
	}
	return false
}
