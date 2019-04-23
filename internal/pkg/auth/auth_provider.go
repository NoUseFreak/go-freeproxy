package auth

type AuthProvider struct {
	users []User
}

func (a *AuthProvider) AddUser(user User) {
	a.users = append(a.users, user)
}

func (a *AuthProvider) Authenticate(user, pwd string) bool {
	for _, u := range a.users {
		if u.Username == user {
			return u.Password == pwd
		}
	}

	return false
}
