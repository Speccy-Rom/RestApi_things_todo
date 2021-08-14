package repository

import (
	"database/sql"
	"github.com/speccy-rom/RestApi_things_todo"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"reflect"
	"testing"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item   todo.TodoItem
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs("title", "description").WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "title",
					Description: "description",
				},
			},
			want: 1,
		},
		{
			name: "Empty Fields",
			mock: func() {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs("", "description").WillReturnError(sqlmock.ErrCancelled)

				mock.ExpectRollback()

			},
			input: args{
				listId: 1,
				item: todo.TodoItem{
					Title:       "",
					Description: "description",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Create(tt.input.listId, tt.input.item)
			if err != nil && !tt.wantErr {
				t.Fatal(err)
			}

			if err == nil && got != tt.want {
				t.Fatalf("Results mismatch; want %d, got %d", tt.want, got)
			}
		})
	}
}

func TestTodoItemPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []todo.TodoItem
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title1", "description1", true).
					AddRow(2, "title2", "description2", false).
					AddRow(3, "title3", "description3", false)

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
			want: []todo.TodoItem{
				{1, "title1", "description1", true},
				{2, "title2", "description2", false},
				{3, "title3", "description3", false},
			},
		},
		{
			name: "No Records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				listId: 1,
				userId: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetAll(tt.input.userId, tt.input.listId)
			if err != nil && !tt.wantErr {
				t.Fatal(err)
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Results mismatch; want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestTodoItemPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		itemId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    todo.TodoItem
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"}).
					AddRow(1, "title1", "description1", true)

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				itemId: 1,
				userId: 1,
			},
			want: todo.TodoItem{1, "title1", "description1", true},
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "done"})

				mock.ExpectQuery("SELECT (.+) FROM todo_items ti INNER JOIN lists_items li on (.+) INNER JOIN users_lists ul on (.+) WHERE (.+)").
					WithArgs(1, 1).WillReturnRows(rows)
			},
			input: args{
				itemId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetById(tt.input.userId, tt.input.itemId)
			if err != nil && !tt.wantErr {
				t.Fatal(err)
			}

			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("Results mismatch; want %v, got %v", tt.want, got)
			}
		})
	}
}

func TestTodoItemPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		itemId int
		userId int
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectExec("DELETE FROM todo_items ti USING lists_items li, users_lists ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
			},
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("DELETE FROM todo_items ti USING lists_items li, users_lists ul WHERE (.+)").
					WithArgs(1, 404).WillReturnError(sql.ErrNoRows)
			},
			input: args{
				itemId: 404,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Delete(tt.input.userId, tt.input.itemId)
			if err != nil && !tt.wantErr {
				t.Fatal(err)
			}
		})
	}
}

func TestTodoItemPostgres_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		itemId int
		userId int
		input  todo.UpdateItemInput
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK_AllFields",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs("new title", "new description", true, 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
				input: todo.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
					Done:        boolPointer(true),
				},
			},
		},
		{
			name: "OK_WithoutDone",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs("new title", "new description", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
				input: todo.UpdateItemInput{
					Title:       stringPointer("new title"),
					Description: stringPointer("new description"),
				},
			},
		},
		{
			name: "OK_WithoutDoneAndDescription",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs("new title", 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
				input: todo.UpdateItemInput{
					Title: stringPointer("new title"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func() {
				mock.ExpectExec("UPDATE todo_items ti SET (.+) FROM lists_items li, users_lists ul WHERE (.+)").
					WithArgs(1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				itemId: 1,
				userId: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Update(tt.input.userId, tt.input.itemId, tt.input.input)
			if err != nil && !tt.wantErr {
				t.Fatal(err)
			}
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}
