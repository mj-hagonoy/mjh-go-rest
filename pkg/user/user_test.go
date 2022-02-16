package user

import (
	"context"
	"encoding/csv"
	"os"
	"testing"
)

func TestCreateUser(t *testing.T) {
	in := []string{"firstname", "lastname", "fl@gmail.com"}
	u, err := NewUser()
	if err != nil {
		t.Fatal(err)
	}

	err = u.Bind(in)
	if err != nil {
		t.Error(err)
	}

	t.Log(u)
}

func TestBulkUpload(t *testing.T) {
	records, err := _getMockUsers(`D:\workspace\github.com\mjh\codev\pkg\utils\test_users.csv`)
	if err != nil {
		t.Fatal(err)
	}

	err = BulkUpload(context.Background(), records)
	if err != nil {
		t.Error(err)
	}
}

func _getMockUsers(filepath string) ([][]string, error) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0555)
	if err != nil {
		panic(err)
	}
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func TestUserStore(t *testing.T) {
	u, _ := NewUser()
	u.FirstName = "test"
	u.LastName = "test"
	u.Email = "test@gmail.com"

	if e := u.Validate(); e != nil {
		t.Error(e)
	}

	id, e := u.Store(context.Background())
	if e != nil {
		t.Error(e)
	}
	t.Logf("inserted id = %s", id)
}

func TestUsersGetAll(t *testing.T) {
	users, err := GetAll(context.Background())
	if err != nil {
		t.Error(err)
	}
	t.Logf("records \n [%+v]", users)

}
