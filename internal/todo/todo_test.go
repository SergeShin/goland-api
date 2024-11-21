package todo_test

import (
	"my-first-api/internal/todo"
	"reflect"
	"testing"
)

func TestService_Search(t *testing.T) {
	tests := []struct {
		name           string
		toDosToAdd     []string
		query          string
		expectedResult []string
	}{
		{
			name:           "given a todo of shop and a search of sh, i should get shop back",
			toDosToAdd:     []string{"shop"},
			query:          "shop",
			expectedResult: []string{"shop"},
		},
		{
			name:           "still returns shop, even if the case doesnt match  ",
			toDosToAdd:     []string{"Shopping"},
			query:          "sh",
			expectedResult: []string{"Shopping"},
		}, {
			name:           "spaces",
			toDosToAdd:     []string{"go Shopping"},
			query:          "go",
			expectedResult: []string{"go Shopping"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := todo.NewService()
			for _, toAdd := range tt.toDosToAdd {
				err := svc.Add(toAdd)
				if err != nil {
					t.Error(err)
				}
			}

			if got := svc.Search(tt.query); !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("Search() = %v, want %v", got, tt.expectedResult)
			}
		})
	}
}
