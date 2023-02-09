package controllers

import (
	"encoding/json"
	"net/http"
	"photo-media/pkg/dao"
	"photo-media/pkg/model"
	"photo-media/pkg/model/card"
	"photo-media/pkg/util"
	"strconv"
)

const (
	BasicCard = 0
	Comment   = 1
)

func CardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		createCard(w, r)
	}
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		cardIdStr := r.FormValue("cardId")

		cardId, err := strconv.Atoi(cardIdStr)
		util.CheckError("error parsing string to int", err)

		author := r.FormValue("author")
		body := r.FormValue("body")

		comment := card.Comment{
			CommentId: Comment,
			CardId:    cardId,
			Author:    author,
			Body:      body,
		}
		dao.CreateCardComment(comment)
		http.Redirect(w, r, "http://localhost:8080/home.html", http.StatusMovedPermanently)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		var userAction model.UserActionRequest
		err := decoder.Decode(&userAction)
		util.CheckError("Error decoding json", err)

		//Authenticate user
		if dao.ReadActiveSessionByUser(userAction.UserCredentials) {
			w.Header().Set("Content-Type", "application/json")
			if userAction.Action == "onLoad" {
				//Explore list users
				exploreList := getExploreList(userAction.UserCredentials)
				//Follow list users
				followList := getUsersFollowedByUser(userAction.UserCredentials)
				//Card list
				cardList := getCardList(userAction.UserCredentials)

				var mainContentSlice []model.MainContent
				mainContentSlice = append(mainContentSlice, exploreList)
				mainContentSlice = append(mainContentSlice, followList)
				for i := 0; i < len(cardList); i++ {
					mainContentSlice = append(mainContentSlice, cardList[i])
				}

				json.NewEncoder(w).Encode(mainContentSlice)
			}

		} else {
			msg := model.Message{Category: "error", Content: "we are sorry, your connection expired"}
			sendExpiredConnection(w, msg)
		}
	} else if r.Method == "PUT" {
		decoder := json.NewDecoder(r.Body)
		var userAction model.UserActionRequest
		err := decoder.Decode(&userAction)
		util.CheckError("Error decoding json", err)

		w.Header().Set("Content-Type", "application/json")

		//Authenticate if is valid user
		if dao.ReadActiveSessionByUser(userAction.UserCredentials) {
			if userAction.Action == "follow" {
				createRelationFollow(userAction)

			} else if userAction.Action == "likeCard" {
				cardToUpdate := createRelationCardLiked(userAction)
				json.NewEncoder(w).Encode(cardToUpdate)

			} else if userAction.Action == "likeComment" {
				commentToUpdate := createRelationCommentLiked(userAction)
				json.NewEncoder(w).Encode(commentToUpdate)
			}
		}
	} else if r.Method == "DELETE" {

		decoder := json.NewDecoder(r.Body)
		var userAction model.UserActionRequest
		err := decoder.Decode(&userAction)
		util.CheckError("Error decoding json", err)

		//Authenticate if is valid user
		if dao.ReadActiveSessionByUser(userAction.UserCredentials) {
			if userAction.Action == "unfollow" {
				deleteFollowRelationship(userAction)
			}
		}
	}
}

func createRelationFollow(userAction model.UserActionRequest) bool {

	if dao.ReadUserAlreadyFollowsUser(userAction.Username, userAction.Target) {
		return false
	}

	return dao.CreateFollowRelationship(userAction.Username, userAction.Target)
}

func createRelationCardLiked(userAction model.UserActionRequest) card.BasicCard {
	cardId, err := strconv.Atoi(userAction.Target)
	util.CheckError("Error parsing string to int", err)

	if dao.ReadUserAlreadyLikedCard(userAction.Username, cardId) {
		idToDrop := dao.ReadLikeIdFromCardLikes(userAction.Username, cardId)
		dao.DeleteLikeCardRelation(idToDrop)
	} else {
		dao.CreateRelationCardLiked(userAction.Username, cardId)
	}

	cardToUpdate := card.BasicCard{
		CardId:      cardId,
		Author:      "",
		Image:       "",
		Likes:       dao.ReadLikesById(cardId, BasicCard),
		Liked:       dao.ReadLikedById(cardId, userAction.Username, BasicCard),
		Description: "",
	}

	return cardToUpdate
}

func createRelationCommentLiked(userAction model.UserActionRequest) card.Comment {
	commentId, err := strconv.Atoi(userAction.Target)
	util.CheckError("Error parsing string to int", err)

	if dao.ReadUserAlreadyLikedComment(userAction.Username, commentId) {
		idToDrop := dao.ReadLikeIdFromCommentLikes(userAction.Username, commentId)
		dao.DeleteLikeCommentRelation(idToDrop)
	} else {
		dao.CreateRelationCommentLiked(userAction.Username, commentId)
	}

	commentParentId := dao.ReadCommentParentId(commentId)
	likes := dao.ReadLikesById(commentId, Comment)
	liked := dao.ReadLikedById(commentId, userAction.Username, Comment)

	commentToUpdate := card.Comment{
		CommentId: commentId,
		CardId:    commentParentId,
		Author:    "",
		Body:      "",
		Likes:     likes,
		Liked:     liked,
	}

	return commentToUpdate
}

func deleteFollowRelationship(userAction model.UserActionRequest) bool {
	return dao.DeleteFollowRelationship(userAction.Username, userAction.Target)
}

func getUsersFollowedByUser(user model.UserCredentials) model.UserList {
	usernames := dao.ReadUsersFollowedByUser(user.Username)
	var followingUserList model.UserList
	if usernames != nil {
		relation := "following"
		followingUserList = model.UserList{Relation: relation, Usernames: usernames}
	}
	return followingUserList
}

func getCardList(user model.UserCredentials) []card.Card {
	var cards []card.Card
	basicCards := dao.ReadBasicCardsByAuthor(user.Username)
	for i := 0; i < len(basicCards); i++ {
		//Get and set likes of current card by id
		basicCards[i].Likes = dao.ReadLikesById(basicCards[i].CardId, BasicCard)
		basicCards[i].Liked = dao.ReadLikedById(basicCards[i].CardId, user.Username, BasicCard)

		comments := dao.ReadCommentsByCardId(basicCards[i].CardId)
		for j := 0; j < len(comments); j++ {
			comments[j].Likes = dao.ReadLikesById(comments[j].CommentId, Comment)
			comments[j].Liked = dao.ReadLikedById(comments[j].CommentId, user.Username, Comment)
		}
		card := card.Card{
			BasicCard: basicCards[i],
			Comments:  comments,
		}
		cards = append(cards, card)
	}
	return cards
}

func getExploreList(user model.UserCredentials) model.UserList {
	usernames := dao.ReadNewUserRecommendation(user.Username)
	var exploreUserList model.UserList
	if usernames != nil {
		relation := "explore"
		exploreUserList = model.UserList{Relation: relation, Usernames: usernames}
	}
	return exploreUserList
}

func createCard(w http.ResponseWriter, r *http.Request) {
	//PARSE FORM VALUES INTO GO STRUCTS
	err := r.ParseForm()
	util.CheckError("Error parsing form", err)
	user := model.UserCredentials{
		Username: r.FormValue("author"),
		Token:    r.FormValue("session"),
	}
	imageURL := r.FormValue("image")
	description := r.FormValue("description")

	if imageURL == "" {
		sendImgUrlErrMsg(w)

	} else if dao.ReadActiveSessionByUser(user) { //Authenticate user
		//SEND DATA TO MYSQL
		transaction := dao.CreateBasicCard(user.Username, imageURL, description)
		if transaction {
			http.Redirect(w, r, "http://localhost:8080/home.html", http.StatusMovedPermanently)
		}

	} else {
		msg := model.Message{Category: "error", Content: "we are sorry, your connection expired"}
		sendExpiredConnection(w, msg)
	}

}

func sendImgUrlErrMsg(w http.ResponseWriter) {
	notAcceptableValueMsg := model.Message{
		Category: "error",
		Content:  "Not acceptable image URL",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notAcceptableValueMsg)
}

func sendExpiredConnection(w http.ResponseWriter, message model.Message) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// func sendTransactionResultJSON(transaction bool, w http.ResponseWriter) {
// 	w.Header().Set("Content-Type", "application/json")
// 	transactionMsg := getTransactionMessage(transaction)
// 	json.NewEncoder(w).Encode(transactionMsg)
// }

// func getTransactionMessage(transaction bool) model.Message {
// 	message := model.Message{
// 		Category: "transaction",
// 	}
// 	setTransactionContent(transaction, &message)
// 	return message
// }

// func setTransactionContent(transaction bool, message *model.Message) {
// 	if transaction {
// 		message.Content = "success"
// 	} else {
// 		message.Content = "failed"
// 	}
// }
