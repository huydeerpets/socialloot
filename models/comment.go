package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Comment struct {
	Id string `orm:"pk"`

	// User is the person which wrote the comment
	User *User     `orm:"rel(fk);null;on_delete(do_nothing)"`
	Date time.Time `orm:"auto_now_add"`

	// Text is the comment
	Text string

	// Post that is commented on
	ReplyTo string
}

// Comments creates orm object to get comments from database
func Comments() orm.QuerySeter {
	var table Comment
	return orm.NewOrm().QueryTable(table)
}

// Insert writes the comment to the database
func (c *Comment) Insert() error {
	_, err := orm.NewOrm().Insert(c)
	return err
}

// Read searches the comment in the database by the given field
// If no field is provided the primary key is used
// The result is written to the comment struct
func (c *Comment) Read(fields ...string) error {
	return orm.NewOrm().Read(c, fields...)
}

// Valid checks if the comment struct has valid data
// it only checks the syntax and does not access the database
func (c *Comment) Valid(v *validation.Validation) {
	if c.User == nil {
		v.AddError("User", "Missing user")
	}
	if len(c.Text) < 1 {
		v.AddError("Text", "Missing comment text")
	}
	if len(c.ReplyTo) != ItemIDLength {
		v.AddError("RepleyTo", "Invalid item id")
	}
	if len(c.Id) != ItemIDLength {
		v.AddError("Id", "Invalid comment id")
	}
}

// CommentList is a slice of comments
type CommentList []*Comment

// CommentMetaDataList is a slice of meta data for comments
type CommentMetaDataList []*CommentMetaData

// ToMetaData turns a comment list into a list of meta data
func (c *CommentList) ToMetaData() *CommentMetaDataList {
	metas := make(CommentMetaDataList, len(*c))
	for i, comment := range *c {
		metas[i] = comment.NewMetaData()
	}
	return &metas
}

type CommentMetaData struct {
	*Comment
	Votes   int
	VoteDir VoteDirection
	Replies []*CommentMetaData
}

// NewMetaData creates a new MetaData wrapper from a comment
// This sould be the only way new meta data objects are created!
func (c *Comment) NewMetaData() *CommentMetaData {
	return &CommentMetaData{
		Comment: c,
	}
}

// LoadReplies loads all replies to the comment recursively
// The replies are ordererd by date
func (c *CommentMetaData) LoadReplies(u *User) error {
	var replies CommentList
	if _, err := Comments().Filter("replyto", c.Id).RelatedSel("user").OrderBy("date").All(&replies); err != nil {
		return err
	}
	if err := c.ReadVoteData(u); err != nil {
		beego.Error(err)
	}
	c.Replies = *replies.ToMetaData()
	// load replies recursively
	for i := range c.Replies {
		if err := c.Replies[i].LoadReplies(u); err != nil {
			beego.Error(err)
		}
	}
	return nil
}

// ReadVoteSum reads the sum of votes (upvotes-downvotes) on the comment
func (c *CommentMetaData) ReadVoteSum() error {
	var votes []*Vote
	c.Votes = 0
	if _, err := getVotesOnItem(c.Id).All(&votes); err != nil {
		return err
	}
	for _, v := range votes {
		if v.Action == VoteDirectionUp {
			c.Votes++
		} else if v.Action == VoteDirectionDown {
			c.Votes--
		}
	}
	return nil
}

// ReadVoteOnPost gets the users vote on the given post and safes the result in the post struct
func (c *CommentMetaData) ReadVoteOnPost(p *PostMetaData) error {
	var vote Vote
	if err := getVotesOnItem(c.Id).One(&vote, "action"); err != nil {
		return err
	}
	p.VoteDir = vote.Action
	return nil
}

// ReadVoteData reads the sum of votes and the users vote on the comment
func (c *CommentMetaData) ReadVoteData(u *User) error {
	if u != nil {
		// Get user vote on post
		if err := u.ReadVoteOnComment(c); err != nil {
			return err
		}
	}
	return c.ReadVoteSum()
}
