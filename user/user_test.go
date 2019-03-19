package user

import (
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/asdine/storm"

	"gopkg.in/mgo.v2/bson"
)

func TestMain(m *testing.M) {
	m.Run()
	os.Remove(dbPath)
}

func cleanDb(b *testing.B) {
	os.Remove(dbPath)
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "Ryan",
		Role: "Engineer",
	}
	err := u.Save()
	if err != nil {
		b.Fatalf("Error saving a record: %s", err)
	}
	b.ResetTimer()
}
func BenchmarkCreate(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Ryan" + strconv.Itoa(i),
			Role: "Engineer",
		}
		b.StartTimer()
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
	}
}

func BenchmarkRead(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Ryan" + strconv.Itoa(i),
			Role: "Engineer",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		b.StartTimer()
		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error getting a record: %s", err)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Ryan" + strconv.Itoa(i),
			Role: "Engineer",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		b.StartTimer()
		u.Role = "Tester"
		err = u.Save()
		if err != nil {
			b.Fatalf("Error updating a record: %s", err)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	cleanDb(b)
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Ryan" + strconv.Itoa(i),
			Role: "Engineer",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		b.StartTimer()
		err = Delete(u.ID)
		if err != nil {
			b.Fatalf("Error deleting a record: %s", err)
		}
	}
}
func BenchmarkCRUD(b *testing.B) {
	os.Remove(dbPath)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		u := &User{
			ID:   bson.NewObjectId(),
			Name: "Ryan",
			Role: "Engineer",
		}
		err := u.Save()
		if err != nil {
			b.Fatalf("Error saving a record: %s", err)
		}
		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error getting a record: %s", err)
		}
		u.Role = "Tester"
		err = u.Save()
		if err != nil {
			b.Fatalf("Error updating a record: %s", err)
		}
		_, err = One(u.ID)
		if err != nil {
			b.Fatalf("Error getting a record: %s", err)
		}

		err = Delete(u.ID)
		if err != nil {
			b.Fatalf("Error deleting a record: %s", err)
		}
	}
}

func TestCRUD(t *testing.T) {
	t.Log("Create")
	u := &User{
		ID:   bson.NewObjectId(),
		Name: "Ryan",
		Role: "Engineer",
	}
	err := u.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	t.Log("Retrieve")
	u2, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error getting a record: %s", err)
	}
	if !reflect.DeepEqual(u2, u) {
		t.Error("Records do not match")
	}
	t.Log("Update")
	u.Role = "Tester"
	err = u.Save()
	if err != nil {
		t.Fatalf("Error updating a record: %s", err)
	}
	u3, err := One(u.ID)
	if err != nil {
		t.Fatalf("Error getting a record: %s", err)
	}
	if !reflect.DeepEqual(u3, u) {
		t.Error("Records do not match")
	}
	t.Log("Delete")
	err = Delete(u.ID)
	if err != nil {
		t.Fatalf("Error deleting a record: %s", err)
	}
	_, err = One(u.ID)
	if err == nil {
		t.Fatal("Record should be deleted")
	}
	if err != storm.ErrNotFound {
		t.Fatalf("Record should be deleted: %s", err)
	}
	t.Log("Read all")
	u2.ID = bson.NewObjectId()
	u3.ID = bson.NewObjectId()
	err = u.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	err = u2.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	err = u3.Save()
	if err != nil {
		t.Fatalf("Error saving a record: %s", err)
	}
	users, err := All()
	if err != nil {
		t.Fatalf("Error reading records: %s", err)
	}
	if len(users) != 3 {
		t.Errorf("The number of user is not correct. Expected: 3. Actual: %d", len(users))
	}

}
