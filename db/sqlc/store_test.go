package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before:", account1.Balance, account2.Balance)

	// чтобы убедиться, что наша транзакция работает хорошо - можно запустить несколько горутин
	n := 5
	amount := int64(10)

	// канал для горутин и их обмена данными без явной блокировки
	errs := make(chan error)               // канал для получения ошибок
	results := make(chan TransferTxResult) // канал для получения результата

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Ammount:       amount,
			})

			// теперь внутри кода горутины мы можем отправить сообщение об ошибке и результате в каналы errs и results(канал слева, а данные для отправки справа)
			errs <- err
			results <- result

		}()
	}


	existed := make(map[int]bool)
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

		// TODO: проверим выходные счета и их балансы учетных записей 
		// начнем со счетов - точнее с счета с которого уходят деньги
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		// теперь проверим балансы счетов
		fmt.Println(">> tx:", fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // остаток должен быть кратным amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// проверим окончательный баланс каждого счета
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)

	fmt.Println(">> after:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance - int64(n) * amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n) * amount, updatedAccount2.Balance)

}


