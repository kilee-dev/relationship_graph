package main

import (
	"github.com/google/uuid"

	"errors"
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
		return nil, errors.New(fmt.Sprintf("Queried user %v doesn't exist", user))
	}

	userList := make([]User, 0)
	for followerID, _ := range graph.FollowerRelationship[user.ID] {
		userList = append(userList, graph.Users[followerID])
	}

	return userList, nil
}

func (graph RelationshipGraph) QueryFollowees(user User) ([]User, error) {
	if _, exist := graph.Users[user.ID]; !exist {
		return nil, errors.New(fmt.Sprintf("Queried user %v doesn't exist", user))
	}

	userList := make([]User, 0)
	for followeeID, _ := range graph.FolloweeRelationship[user.ID] {
		userList = append(userList, graph.Users[followeeID])
	}

	return userList, nil
}

func (graph RelationshipGraph) AddUser(user User) (bool, error) {
	// edge case - user already registered in the graph
	if _, exist := graph.Users[user.ID]; exist {
		return false, errors.New(fmt.Sprintf("Queried user %v already registered", user))
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
		return false, errors.New(fmt.Sprintf("Queried user %v doesn't exist", follower))
	}

	if _, exist := graph.Users[followee.ID]; !exist {
		return false, errors.New(fmt.Sprintf("Queried user %v doesn't exist", followee))
	}

	// edge case when follower is already following followee
	if _, exist := graph.FollowerRelationship[follower.ID][followee.ID]; exist {
		return false, errors.New(fmt.Sprintf("User %v is already following user %v", follower, followee))
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
}
