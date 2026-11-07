package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableMemoryTest(t *testing.T) {
	testTable := []struct {
		name   string
		arrang func(ctx context.Context, t *testing.T, m *Memory)
		assert func(ctx context.Context, t *testing.T, m *Memory)
	}{
		{
			name: "Добавили одну запись в память",
			arrang: func(ctx context.Context, t *testing.T, m *Memory) {
				m.Add(ctx, -1, 1)
			},
			assert: func(ctx context.Context, t *testing.T, m *Memory) {
				ok, _ := m.Exists(ctx, -1, 1)
				assert.True(t, ok, "Значение успешно получили")
			},
		},
		{
			name: "Добавили одну запись в память и удалили ее",
			arrang: func(ctx context.Context, t *testing.T, m *Memory) {
				require.NoError(t, m.Add(ctx, -12, 12), "добавляем без ошибки")
				m.Remove(ctx, -12, 12)
			},
			assert: func(ctx context.Context, t *testing.T, m *Memory) {
				ok, _ := m.Exists(ctx, -12, 12)
				assert.False(t, ok, "Значение успешно НЕ получили")
			},
		},
		{
			name: "Добавили две запись в память и удалили одну",
			arrang: func(ctx context.Context, t *testing.T, m *Memory) {
				require.NoError(t, m.Add(ctx, -123, 123), "добавляем без ошибки")
				require.NoError(t, m.Add(ctx, -1234, 1234), "добавляем без ошибки")
				require.NoError(t, m.Remove(ctx, -123, 123), "удаляем без ошибки")
			},
			assert: func(ctx context.Context, t *testing.T, m *Memory) {
				ok, _ := m.Exists(ctx, -123, 123)
				assert.False(t, ok, "Значение успешно НЕ получили")
				ok, _ = m.Exists(ctx, -1234, 1234)
				assert.True(t, ok, "Значение успешно получили")
			},
		},
		{
			name: "Добавили две запись в память и удалили одну дважды",
			arrang: func(ctx context.Context, t *testing.T, m *Memory) {
				require.NoError(t, m.Add(ctx, -123, 123), "добавляем без ошибки")
				require.NoError(t, m.Add(ctx, -1234, 1234), "добавляем без ошибки")
				require.NoError(t, m.Remove(ctx, -123, 123), "удаляем без ошибки")
				require.NoError(t, m.Remove(ctx, -123, 123), "удаляем без ошибки")
			},
			assert: func(ctx context.Context, t *testing.T, m *Memory) {
				ok, _ := m.Exists(ctx, -123, 123)
				assert.False(t, ok, "Значение успешно НЕ получили")
				ok, _ = m.Exists(ctx, -1234, 1234)
				assert.True(t, ok, "Значение успешно получили")
			},
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(tt *testing.T) {
			ctx := context.Background()
			m := NewMemAdapter()

			testCase.arrang(ctx, tt, m)
			testCase.assert(ctx, tt, m)
		})
	}
}
