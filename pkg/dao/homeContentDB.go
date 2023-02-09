package dao

import (
	"photo-media/pkg/model"
	"photo-media/pkg/model/card"
	"photo-media/pkg/util"
)

const (
	create_follow_relationship    = "INSERT INTO `photo-media`.`follows` (`follower_user`, `followed_user`) VALUES (?, ?);"
	create_relation_card_liked    = "INSERT INTO `photo-media`.`card_likes` (`username`, `card_id`) VALUES (?, ?);"
	create_relation_comment_liked = "INSERT INTO `photo-media`.`comment_likes` (`username`, `comment_id`) VALUES (?, ?);"
	create_basic_card             = "INSERT INTO `photo-media`.`card` (`author`, `image`, `description`) VALUES (?, ?, ?);"
	create_card_comment           = "INSERT INTO `photo-media`.`comment` (`card_id`, `author`, `body`) VALUES (?, ?, ?);"

	read_basic_card_by_author       = "SELECT card_id, author, image, description FROM `photo-media`.card WHERE author=? || author IN (SELECT followed_user FROM follows WHERE follower_user = ?) ORDER BY created DESC;"
	read_card_likes_by_id           = "SELECT COUNT(*) FROM `photo-media`.card_likes WHERE card_id = ?"
	read_card_liked_by_id           = "SELECT COUNT(*) FROM `photo-media`.card_likes WHERE card_id = ? && username=?"
	read_comment_likes_by_id        = "SELECT COUNT(*) FROM `photo-media`.comment_likes WHERE comment_id = ?"
	read_comment_liked_by_id        = "SELECT COUNT(*) FROM `photo-media`.comment_likes WHERE comment_id = ? && username=?"
	read_comment_parent_id          = "SELECT card_id FROM `photo-media`.comment WHERE comment_id = ?"
	read_session_by_user            = "SELECT username FROM `photo-media`.sessions WHERE username = ? && session_token=?;"
	read_already_follows            = "SELECT COUNT(*) FROM `photo-media`.follows WHERE follower_user = ? && followed_user=?;"
	read_already_liked_card         = "SELECT COUNT(*) FROM `photo-media`.`card_likes` WHERE username = ? && card_id = ?;"
	read_like_id_from_card_likes    = "SELECT like_id FROM `photo-media`.card_likes WHERE username = ? && card_id = ?"
	read_already_liked_comment      = "SELECT COUNT(*) FROM `photo-media`.`comment_likes` WHERE username = ? && comment_id = ?;"
	read_like_id_from_comment_likes = "SELECT like_id FROM `photo-media`.comment_likes WHERE username = ? && comment_id = ?"

	//DELETE FROM `photo-media`.`card_likes` WHERE like_id = (SELECT like_id FROM `photo-media`.`card_likes` WHERE username = ? && card_id = ?)""
	read_users_followed_by_user = "SELECT followed_user FROM `photo-media`.follows WHERE follower_user = ?;"
	read_users_not_followed     = "SELECT username FROM `photo-media`.accounts WHERE username != ? && username NOT IN (SELECT followed_user from follows where follower_user = ?)"
	read_comments_by_card_id    = "SELECT comment_id, card_id, author, body FROM `photo-media`.comment WHERE card_id=?;"

	delete_follow_relationship   = "DELETE FROM `photo-media`.follows WHERE follower_user=? && followed_user=?;"
	delete_relation_like_card    = "DELETE FROM `photo-media`.`card_likes` WHERE (`like_id` = ?);"
	delete_relation_like_comment = "DELETE FROM `photo-media`.`comment_likes` WHERE (`like_id` = ?);"
	delete_basic_card_by_id      = "DELETE FROM `photo-media`.`card` WHERE (`card_id` = ?);"
)

func CreateCardComment(comment card.Comment) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(create_card_comment)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(comment.CardId, comment.Author, comment.Body)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func CreateBasicCard(author, image string, description string) bool {
	connect()
	defer db.Close()

	card := card.BasicCard{
		CardId:      0,
		Author:      author,
		Image:       image,
		Likes:       0,
		Description: description,
	}

	st, err := db.Prepare(create_basic_card)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(card.Author, card.Image, card.Description)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func CreateFollowRelationship(follower, targerUsername string) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(create_follow_relationship)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(follower, targerUsername)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func CreateRelationCardLiked(username string, cardId int) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(create_relation_card_liked)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(username, cardId)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func CreateRelationCommentLiked(username string, commentId int) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(create_relation_comment_liked)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(username, commentId)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func ReadUserAlreadyFollowsUser(followerUser, followedUser string) bool {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_already_follows)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(followerUser, followedUser)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var followCount int
	for rows.Next() {
		err := rows.Scan(&followCount)
		util.CheckError("Error reading the columns", err)
	}
	return followCount != 0
}

func ReadUserAlreadyLikedCard(username string, cardId int) bool {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_already_liked_card)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(username, cardId)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var likeCount int
	for rows.Next() {
		err := rows.Scan(&likeCount)
		util.CheckError("Error reading the columns", err)
	}
	return likeCount != 0
}

func ReadUserAlreadyLikedComment(username string, commentId int) bool {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_already_liked_comment)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(username, commentId)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var likeCount int
	for rows.Next() {
		err := rows.Scan(&likeCount)
		util.CheckError("Error reading the columns", err)
	}
	return likeCount != 0
}

func ReadCommentParentId(commentId int) int {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_comment_parent_id)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(commentId)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var cardId int
	for rows.Next() {
		err := rows.Scan(&cardId)
		util.CheckError("Error reading the columns", err)
	}
	return cardId
}

func ReadLikedById(cardId int, username string, refType int) bool {
	refTypeQry := getLikedQry(refType)

	connect()
	defer db.Close()
	st, err := db.Prepare(refTypeQry)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(cardId, username)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var likes int
	for rows.Next() {
		err := rows.Scan(&likes)
		util.CheckError("Error reading the columns", err)
	}
	return likes != 0
}

func getLikedQry(refType int) string {
	switch refType {
	case 0:
		return read_card_liked_by_id
	case 1:
		return read_comment_liked_by_id
	default:
		return ""
	}
}

func ReadLikesById(cardId int, refType int) int {
	readLikesQry := getQueryLikesType(refType)

	connect()
	defer db.Close()
	st, err := db.Prepare(readLikesQry)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(cardId)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var likes int
	for rows.Next() {
		err := rows.Scan(&likes)
		util.CheckError("Error reading the columns", err)
	}
	return likes
}

func ReadLikeIdFromCardLikes(username string, cardId int) int {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_like_id_from_card_likes)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(username, cardId)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var likeId int
	for rows.Next() {
		err := rows.Scan(&likeId)
		util.CheckError("Error reading the columns", err)
	}
	return likeId
}

func ReadLikeIdFromCommentLikes(username string, commentId int) int {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_like_id_from_comment_likes)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(username, commentId)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var likeId int
	for rows.Next() {
		err := rows.Scan(&likeId)
		util.CheckError("Error reading the columns", err)
	}
	return likeId
}

func getQueryLikesType(refType int) string {
	switch refType {
	case 0:
		return read_card_likes_by_id
	case 1:
		return read_comment_likes_by_id
	default:
		return ""
	}
}

func ReadCommentsByCardId(cardId int) []card.Comment {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_comments_by_card_id)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(cardId)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var comments []card.Comment
	var comment card.Comment
	for rows.Next() {
		err := rows.Scan(&comment.CommentId, &comment.CardId, &comment.Author, &comment.Body)
		util.CheckError("Error reading the columns", err)
		comments = append(comments, comment)
	}
	return comments
}

func ReadBasicCardsByAuthor(author string) []card.BasicCard {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_basic_card_by_author)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(author, author)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var basicCards []card.BasicCard
	var basicCard card.BasicCard
	for rows.Next() {
		err := rows.Scan(&basicCard.CardId, &basicCard.Author, &basicCard.Image, &basicCard.Description)
		util.CheckError("Error reading the columns", err)
		basicCards = append(basicCards, basicCard)
	}
	return basicCards
}

func ReadActiveSessionByUser(user model.UserCredentials) bool {
	//Returns true if there is at least one active session with this user
	connect()
	defer db.Close()
	st, err := db.Prepare(read_session_by_user)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(user.Username, user.Token)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var usernameFound string
	for rows.Next() {
		err := rows.Scan(&usernameFound)
		util.CheckError("Error reading the columns", err)
	}
	return len(usernameFound) > 0
}

func ReadUsersFollowedByUser(username string) []string {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_users_followed_by_user)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(username)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var usernames []string
	var user string
	for rows.Next() {
		err := rows.Scan(&user)
		util.CheckError("Error reading the columns", err)
		usernames = append(usernames, user)
	}
	return usernames
}

func ReadNewUserRecommendation(username string) []string {
	connect()
	defer db.Close()
	st, err := db.Prepare(read_users_not_followed)
	util.CheckError("Error preparing the statement", err)
	defer st.Close()

	rows, err := st.Query(username, username)
	util.CheckError("Error getting rows from query", err)
	defer rows.Close()

	var usernames []string
	var user string
	for rows.Next() {
		err := rows.Scan(&user)
		util.CheckError("Error reading the columns", err)
		usernames = append(usernames, user)
	}
	return usernames
}

func DeleteBasicCardById(id int) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(delete_basic_card_by_id)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(id)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func DeleteFollowRelationship(follower, targetUser string) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(delete_follow_relationship)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(follower, targetUser)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func DeleteLikeCardRelation(likeId int) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(delete_relation_like_card)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(likeId)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}

func DeleteLikeCommentRelation(likeId int) bool {
	connect()
	defer db.Close()

	st, err := db.Prepare(delete_relation_like_comment)
	util.CheckError("Error preparing qry", err)
	sqlResutl, err := st.Exec(likeId)
	util.CheckError("Error executing the qry", err)
	rows, err := sqlResutl.RowsAffected()
	util.CheckError("Error getting the rows affected", err)
	defer st.Close()
	return rows > 0
}
