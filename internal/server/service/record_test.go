package service

import (
	"context"
	"errors"
	"testing"

	"github.com/lenarlenar/mygokeeper/internal/server/entity"
)

type fakeRecordRepo struct {
	data   map[int64][]entity.Record // userID → []Record
	autoID int64
}

func newFakeRecordRepo() *fakeRecordRepo {
	return &fakeRecordRepo{
		data: make(map[int64][]entity.Record),
	}
}

func (r *fakeRecordRepo) Save(_ context.Context, rec *entity.Record) error {
	r.autoID++
	rec.ID = r.autoID
	r.data[rec.UserID] = append(r.data[rec.UserID], *rec)
	return nil
}

func (r *fakeRecordRepo) GetAllByUser(_ context.Context, userID int64) ([]entity.Record, error) {
	return r.data[userID], nil
}

func (r *fakeRecordRepo) Update(_ context.Context, rec *entity.Record) error {
	recs := r.data[rec.UserID]
	for i := range recs {
		if recs[i].ID == rec.ID {
			recs[i] = *rec
			r.data[rec.UserID] = recs
			return nil
		}
	}
	return errors.New("not found")
}

func (r *fakeRecordRepo) DeleteByID(_ context.Context, userID, id int64) error {
	recs := r.data[userID]
	for i, rec := range recs {
		if rec.ID == id {
			r.data[userID] = append(recs[:i], recs[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func TestRecordService_CRUD(t *testing.T) {
	ctx := context.Background()
	repo := newFakeRecordRepo()
	service := NewRecordService(repo)

	// Save
	rec := &entity.Record{
		UserID: 1,
		Type:   "password",
		Data:   []byte("secret123"),
		Meta:   "github",
	}
	err := service.Save(ctx, rec)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if rec.ID == 0 {
		t.Fatal("record ID not set")
	}

	// Get
	all, err := service.GetAll(ctx, 1)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if len(all) != 1 || all[0].Meta != "github" {
		t.Fatalf("unexpected get result: %+v", all)
	}

	// Update
	rec.Meta = "new-meta"
	err = service.Update(ctx, rec)
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	all, _ = service.GetAll(ctx, 1)
	if all[0].Meta != "new-meta" {
		t.Fatal("record not updated")
	}

	// Delete
	err = service.Delete(ctx, 1, rec.ID)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	all, _ = service.GetAll(ctx, 1)
	if len(all) != 0 {
		t.Fatal("record not deleted")
	}
}
