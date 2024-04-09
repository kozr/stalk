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

func GetOnlineFollowers(userId string) ([]string, error) {
	rows, err := DB.Query("SELECT follower_id FROM followers WHERE user_id = $1 AND online = true", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []string
	for rows.Next() {
		var followerId string
		err := rows.Scan(&followerId)
		if err != nil {
			return nil, err
		}
		followers = append(followers, followerId)
	}

	return followers, nil
}

func AddFollower(userId string, followerId string) error {
	_, err := DB.Exec("INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)", userId, followerId)
	if err != nil {
		return err
	}
	return nil
}

func RemoveFollower(userId string, followerId string) error {
	_, err := DB.Exec("DELETE FROM followers WHERE user_id = $1 AND follower_id = $2", userId, followerId)
	if err != nil {
		return err
	}
	return nil
}
