package users

type User struct {
	ID    string
	Name  string
	Email string
}

var usersDb = []User{
	{ID: "sf342sdf4", Name: "User one", Email: "user.one@example.com"},
	{ID: "ufs43s6f4", Name: "User two", Email: "user.two@example.com"},
	{ID: "hfs3l21f4", Name: "User three", Email: "user.three@example.com"},
	{ID: "hfp3l21f4", Name: "User four", Email: "user.four@example.com"},
}

func AttatchDataPipeline() {
}

func GetUserByEmail(email string) *User {
	var user *User = nil
	for _, u := range usersDb {
		if u.Email == email {
			user = &u
			break
		}
	}
	return user
}
