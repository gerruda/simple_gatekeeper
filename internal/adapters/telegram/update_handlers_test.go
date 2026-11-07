package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseDataFromCallback(t *testing.T) {
	tests := []struct {
		name string
		data string
		want int64
	}{
		{
			name: "HappyPath",
			data: "Accept:1295070216",
			want: 1295070216,
		},
		{
			name: "без прекфикса",
			data: "1295070216",
			want: 0,
		},
		{
			name: "пустая строка",
			data: "",
			want: 0,
		},
		{
			name: "без числа",
			data: "Accept:",
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseDataFromCallback(tt.data)
			assert.Equal(t, tt.want, got, "Должны совпадать")
		})
	}
}
