package main

import (
	"github.com/google/uuid"

	"fmt"
	"os"
)

type User struct {
	ID   uuid.UUID
	Name string
}

type RelationshipGraph struct {
	// map user.id to User
	Users map[uuid.UUID]User
	// first key -> The List of User, second key -> Follower User
	FollowerRelationship map[uuid.UUID]map[uuid.UUID]bool
	// first key -> The List of User, second key -> Followee User
	FolloweeRelationship map[uuid.UUID]map[uuid.UUID]bool
}

func (graph RelationshipGraph) QueryAllUsers() ([]User, error) {
	userList := make([]User, 0)
	for _, user := range graph.Users {
		userList = append(userList, user)
	}

	return userList, nil
}

func (graph RelationshipGraph) QueryFollowers(user User) ([]User, error) {
	// edge case - user not exists
	if _, exist := graph.Users[user.ID]; !exist {
		return nil, fmt.Errorf("queried user %v doesn't exist", user)
	}

	userList := make([]User, 0)
	for followerID, _ := range graph.FollowerRelationship[user.ID] {
		userList = append(userList, graph.Users[followerID])
	}

	return userList, nil
}

func (graph RelationshipGraph) QueryFollowees(user User) ([]User, error) {
	if _, exist := graph.Users[user.ID]; !exist {
		return nil, fmt.Errorf("queried user %v doesn't exist", user)
	}

	userList := make([]User, 0)
	for followeeID, _ := range graph.FolloweeRelationship[user.ID] {
		userList = append(userList, graph.Users[followeeID])
	}

	return userList, nil
}

func (graph RelationshipGraph) QueryFriendRelationship(user1 User, user2 User) (bool, error) {
	if _, exist := graph.Users[user1.ID]; !exist {
		return false, fmt.Errorf("queried user %v not exist", user1)
	}

	if _, exist := graph.Users[user2.ID]; !exist {
		return false, fmt.Errorf("queried user %v not exist", user2)
	}

	if _, firstUserFollowingSecondUser := graph.FollowerRelationship[user1.ID][user2.ID]; !firstUserFollowingSecondUser {
		return false, nil
	}

	if _, secondUserFollowingFirstUser := graph.FollowerRelationship[user2.ID][user1.ID]; !secondUserFollowingFirstUser {
		return false, nil
	}

	return true, nil
}

func (graph RelationshipGraph) QueryMutualFriends(user1 User, user2 User) ([]User, error) {
	if _, exist := graph.Users[user1.ID]; !exist {
		return nil, fmt.Errorf("queried user %v not exist", user1)
	}

	if _, exist := graph.Users[user2.ID]; !exist {
		return nil, fmt.Errorf("queried user %v not exist", user2)
	}

	firstUserFriendMap := make(map[uuid.UUID]bool, 0)
	// user 1 -> follow -> follower
	for followerID, _ := range graph.FollowerRelationship[user1.ID] {
		if _, isFriend := graph.FolloweeRelationship[user1.ID][followerID]; isFriend {
			firstUserFriendMap[followerID] = true
		}
	}

	mutualFriendList := make([]User, 0)
	// get all the list of followees
	for followerID, _ := range graph.FolloweeRelationship[user2.ID] {
		// if follwer is not a friend of user1
		if _, isFirstUserFriend := firstUserFriendMap[followerID]; !isFirstUserFriend {
			continue
		}

		// if follower is not a friend of user2
		if _, isFriend := graph.FolloweeRelationship[user2.ID][followerID]; !isFriend {
			continue
		}

		mutualFriendList = append(mutualFriendList, graph.Users[followerID])
	}

	return mutualFriendList, nil
}

func (graph RelationshipGraph) QueryUnfollowingFriends(user User) ([]User, error) {
	if _, exist := graph.Users[user.ID]; !exist {
		return nil, fmt.Errorf("queried user %v not exist", user)
	}

	unfollowingFriendList := make([]User, 0)
	for followeeID, _ := range graph.FolloweeRelationship[user.ID] {
		if _, exist := graph.FollowerRelationship[user.ID][followeeID]; !exist {
			unfollowingFriendList = append(unfollowingFriendList, graph.Users[followeeID])
		}
	}

	return unfollowingFriendList, nil
}

func (graph RelationshipGraph) AddUser(user User) (bool, error) {
	// edge case - user already registered in the graph
	if _, exist := graph.Users[user.ID]; exist {
		return false, fmt.Errorf("queried user %v already registered", user)
	}

	graph.Users[user.ID] = user

	if _, exist := graph.FollowerRelationship[user.ID]; !exist {
		graph.FollowerRelationship[user.ID] = make(map[uuid.UUID]bool, 0)
	}

	if _, exist := graph.FolloweeRelationship[user.ID]; !exist {
		graph.FolloweeRelationship[user.ID] = make(map[uuid.UUID]bool, 0)
	}

	return true, nil
}

func (graph RelationshipGraph) FollowUser(follower User, followee User) (bool, error) {
	// edge case when follower and followee doesn't exist
	if _, exist := graph.Users[follower.ID]; !exist {
		return false, fmt.Errorf("queried user %v doesn't exist", follower)
	}

	if _, exist := graph.Users[followee.ID]; !exist {
		return false, fmt.Errorf("queried user %v doesn't exist", followee)
	}

	// edge case when follower is already following followee
	if _, exist := graph.FollowerRelationship[follower.ID][followee.ID]; exist {
		return false, fmt.Errorf("User %v is already following user %v", follower, followee)
	}

	graph.FollowerRelationship[follower.ID][followee.ID] = true

	graph.FolloweeRelationship[followee.ID][follower.ID] = true

	return true, nil
}

func main() {
	userA := User{
		ID:   uuid.New(),
		Name: "giung.lee",
	}
	userB := User{
		ID:   uuid.New(),
		Name: "lei.su",
	}
	userC := User{ // non-registered user
		ID:   uuid.New(),
		Name: "phillip.teng",
	}
	relationshipGraph := RelationshipGraph{
		Users:                make(map[uuid.UUID]User, 0),
		FollowerRelationship: make(map[uuid.UUID]map[uuid.UUID]bool, 0),
		FolloweeRelationship: make(map[uuid.UUID]map[uuid.UUID]bool, 0),
	}

	// mutations
	_, err := relationshipGraph.AddUser(userA)
	if err != nil {
		fmt.Printf("relationshipGraph.AddUser err=%v\n", err)
		os.Exit(1)
	}

	_, err = relationshipGraph.AddUser(userB)
	if err != nil {
		fmt.Printf("relationshipGraph.AddUser err=%v\n", err)
		os.Exit(1)
	}

	// edge case - duplicate user registration
	_, err = relationshipGraph.AddUser(userA)
	if err != nil {
		fmt.Printf("relationshipGraph.AddUser err=%v\n", err)
	}

	_, err = relationshipGraph.FollowUser(userA, userB)
	if err != nil {
		fmt.Printf("relationshipGraph.FollowUser err=%v\n", err)
		os.Exit(1)
	}

	// edge case - unregistered user userC
	_, err = relationshipGraph.FollowUser(userB, userC)
	if err != nil {
		fmt.Printf("relationshipGraph.FollowUser err=%v\n", err)
	}

	// query
	followerList, err := relationshipGraph.QueryFollowers(userA)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryFollowers err=%v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryFollowers expected:true, actual:%t\n", len(followerList) == 1)
	}

	followeeList, err := relationshipGraph.QueryFollowees(userA)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryFollowees err=%v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryFollowees expected:true, actual:%t\n", len(followeeList) == 0)
	}

	followerList, err = relationshipGraph.QueryFollowers(userB)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryFollowers err=%v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryFollowers expected:true, actual:%t\n", len(followerList) == 0)
	}

	followeeList, err = relationshipGraph.QueryFollowees(userB)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryFollowees err=%v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryFollowees expected:true, actual:%t\n", len(followeeList) == 1)
	}

	// edge case - user not exist
	_, err = relationshipGraph.QueryFollowers(userC)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryFollowees err=%v\n", err)
	}

	// edge case - user not exist
	_, err = relationshipGraph.QueryFollowees(userC)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryFollowees err=%v\n", err)
	}

	userList, err := relationshipGraph.QueryAllUsers()
	if err != nil {
		fmt.Printf("relationshipGraph.QueryAllUsers err=%v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryAllUsers expected:%d, actual:%d\n", 2, len(userList))
	}

	unfollowingFriendList, err := relationshipGraph.QueryUnfollowingFriends(userA)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryUnfollowingFriends err=%v", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryUnfollowingFriends expected:%d, actual:%d\n", 0, len(unfollowingFriendList))
	}

	unfollowingFriendList, err = relationshipGraph.QueryUnfollowingFriends(userB)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryUnfollowingFriends err=%v", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryUnfollowingFriends expected:%d, actual:%d\n", 1, len(unfollowingFriendList))
	}

	// edge case - not registered user
	_, err = relationshipGraph.QueryUnfollowingFriends(userC)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryUnfollowingFriends err=%v\n", err)
	}

	isFriendEachOther, err := relationshipGraph.QueryFriendRelationship(userA, userB)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryFriendRelationship err=%v\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryFriendRelationship expected:%t, actual:%t\n", false, isFriendEachOther)
	}

	_, err = relationshipGraph.FollowUser(userB, userA)
	if err != nil {
		fmt.Printf("relationshipGraph.FollowUser err=%v\n", err)
		os.Exit(1)
	}

	isFriendEachOther, err = relationshipGraph.QueryFriendRelationship(userA, userB)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryFriendRelationship err=%v", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryFriendRelationship expected:%t, actual:%t\n", true, isFriendEachOther)
	}

	_, err = relationshipGraph.AddUser(userC)
	if err != nil {
		fmt.Printf("relationshipGraph.AddUser err=%v\n", err)
		os.Exit(1)
	}

	_, err = relationshipGraph.FollowUser(userC, userA)
	if err != nil {
		fmt.Printf("relationshipGraph.FollowUser err=%v\n", err)
		os.Exit(1)
	}

	_, err = relationshipGraph.FollowUser(userA, userC)
	if err != nil {
		fmt.Printf("relationshipGraph.FollowUser err=%v\n", err)
		os.Exit(1)
	}

	mutualFriendList, err := relationshipGraph.QueryMutualFriends(userB, userC)
	if err != nil {
		fmt.Printf("relationshipGraph.QueryMutualFriends err=%v", err)
		os.Exit(1)
	} else {
		fmt.Printf("relationshipGraph.QueryMutualFriends expected:%d, actual:%d\n", 1, len(mutualFriendList))
	}
}
