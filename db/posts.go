package db

import (
	"errors"
)

type Post struct {
	PostId uint64 `json:"postId"`
	Title  string `json:"title"`
	Owner  string `json:"owner"`
}

type PostDB struct {
	postCounter uint64
	posts       map[uint64]*Post
}

// initiatialise a new PostDB instance with an empty map of posts
func NewPostDB() *PostDB {
	return &PostDB{
		posts: make(map[uint64]*Post),
	}
}

// create a new post
func (pdb *PostDB) CreatePost(owner string, post Post) uint64 {

	pdb.postCounter++
	pdb.posts[pdb.postCounter] = &Post{
		PostId: pdb.postCounter,
		Title:  post.Title,
		Owner:  owner,
	}

	return pdb.postCounter
}

// update a post
func (pdb *PostDB) UpdatePost(postId uint64, post Post) error {

	po, found := pdb.posts[postId]
	if !found {
		return errors.New("post not found")
	}

	// update title
	po.Title = post.Title

	return nil
}

// delete a post
func (pdb *PostDB) DeletePost(postId uint64) error {

	_, found := pdb.posts[postId]
	if !found {
		return errors.New("post not found")
	}

	// delete a post from the map, having the id
	delete(pdb.posts, postId)

	return nil
}

// get a post by Id
func (pdb *PostDB) GetPost(postId uint64) (Post, error) {

	po, found := pdb.posts[postId]
	if !found {
		return Post{}, errors.New("post not found")
	}

	return *po, nil
}
