package user

type Option func(*User) error

func ID(id string) Option {
	return func(u *User) error {
		u.ID = id
		return nil
	}
}

func Active(isActive bool) Option {
	return func(u *User) error {
		u.Activated = isActive
		return nil
	}
}
