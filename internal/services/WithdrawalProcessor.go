package services

import (
	"github.com/naneri/diploma/internal/user"
	"github.com/naneri/diploma/internal/withdrawal"
	"gorm.io/gorm"
)

func PerformWithdraw(db *gorm.DB, userRepo *user.DBRepository, userID uint32, sum float64, orderID uint) error {
	userRepo.Access.Lock()
	defer userRepo.Access.Unlock()
	return db.Transaction(func(tx *gorm.DB) error {
		var dbUser user.User

		if err := tx.Where("id", userID).First(&dbUser).Error; err != nil {
			return err
		}

		if sum > dbUser.Balance {
			return &NotEnoughBalanceError{
				balance: dbUser.Balance,
				sum:     sum,
			}
		}
		dbUser.Balance += sum * -1
		dbUser.WithDrawn += sum

		if err := tx.Save(&dbUser).Error; err != nil {
			return err
		}

		dbWithdrawal := withdrawal.Withdrawal{
			UserID:  userID,
			Sum:     sum,
			OrderID: orderID,
		}

		if err := tx.Save(&dbWithdrawal).Error; err != nil {
			return err
		}
		return nil
	})
}
