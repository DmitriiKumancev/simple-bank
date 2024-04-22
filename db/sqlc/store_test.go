package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// чтобы убедиться, что наша транзакция работает хорошо - можно запустить несколько горутин
	n := 5
	amount := int64(10)

	// канал для горутин и их обмена данными без явной блокировки
	errs := make(chan error)               // канал для получения ошибок
	results := make(chan TransferTxResult) // канал для получения результата

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Ammount:       amount,
			})

			// теперь внутри кода горутины мы можем отправить сообщение об ошибке и результате в каналы errs и results(канал слева, а данные для отправки справа)
			errs <- err
			results <- result

		}()
	}

	// затем мы проверяем эти ошибки и результаты извне
	for i := 0; i < n; i++ {
		// чтобы получить ошибку из канала мы используем тот же оператор <-, что и внутри функции, но на этот раз канал находится справа от стрелки, а переменная для хранения полученных данных слева
		err := <-errs
		//полученная ошибка должна быть равна нулю
		require.NoError(t, err)

		// аналогичные действия для результата
		result := <-results
		require.NotEmpty(t, result)

		// поскольку результат содержит несколько объектов внутри, мы должны проверить каждый из них
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		// убедимся что запись о переводе действительно создана в бд 
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// теперь мы проверим записи в таблице accounts, которые указывают на выполнение перевода
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		// попытаемся получить запись учетной записи, чтобы убедиться, что она создалась
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// TODO: проверим балансы учетных записей 
	}
}


