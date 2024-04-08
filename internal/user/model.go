package user

type User struct {
	ID         string
	UserName   string
	FirstName  string
	LastName   string
	IsBot      bool
	IsJoined   bool
	DateJoined string
	DateLeft   string
}

func (u *User) ToDomain() User {
	c := User{
		u.ID,
		u.UserName,
		u.FirstName,
		u.LastName,
		u.IsBot,
		u.IsBot,
		u.DateJoined,
		u.DateLeft,
	}

	return c
}
