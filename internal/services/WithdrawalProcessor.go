package services

import (
	"github.com/naneri/diploma/internal/user"
	"github.com/naneri/diploma/internal/withdrawal"
	"gorm.io/gorm"
)

func PerformWithdraw(db *gorm.DB, userRepo *user.DbRepository, userId uint32, sum float64, orderId uint) error {
	userRepo.Access.Lock()
	defer userRepo.Access.Unlock()
	return db.Transaction(func(tx *gorm.DB) error {
		var dbUser user.User

		if err := tx.Where("user_id", userId).First(&dbUser).Error; err != nil {
			return err
		}

		if sum > dbUser.Balance {
			return &NotEnoughBalanceError{
				balance: dbUser.Balance,
				sum:     sum,
			}
		}
		dbUser.Balance += sum * -1

		if err := tx.Save(&dbUser).Error; err != nil {
			return err
		}

		dbWithdrawal := withdrawal.Withdrawal{
			UserId:  userId,
			Sum:     sum,
			OrderId: orderId,
		}

		if err := tx.Save(&dbWithdrawal).Error; err != nil {
			return err
		}
		return nil
	})
}
