package database

import "fmt"

func UpdateUserOnlineStatus(userId string, online bool) error {
	result, err := DB.Exec("UPDATE users SET online = $1 WHERE id = $2", online, userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected for user %s", userId)
	}

	return nil
}
