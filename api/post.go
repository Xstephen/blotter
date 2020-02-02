package api

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/goutils/time"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

type PostRequest struct {
	URL string `json:"url"`
}

func Post(context *register.HandleContext) (err error) {
	output.Debug("call friends")
	args := PostRequest{}
	context.RequestArgs(&args)

	res := make([]PostUnix, 0)
	_, err = mongo.Aggregate(
		"blotter", "posts",
		[]bson.M{
			{
				"$match": bson.M{
					"url": args.URL,
				},
			},
			{
				"$limit": 1,
			},
			{
				"$lookup": bson.M{
					"from":         "tags",
					"localField":   "tags",
					"foreignField": "_id",
					"as":           "tags",
				},
			},
		},
		nil, &res,
	)
	if err != nil {
		return
	}
	go func() {
		mongo.Update(
			"blotter", "posts", bson.M{"url": args.URL},
			bson.M{"$inc": bson.M{"view": 1}}, nil,
		)
	}()

	if len(res) > 0 {
		context.ReturnJSON(res[0].ToPostTime())
	} else {
		context.Response.WriteHeader(404)
	}
	return
}

type PostsRequest struct {
	Number int64  `json:"number"`
	Offset int64  `json:"offset"`
	Type   string `json:"type"`
	Arg    string `json:"arg"`
}

type PostsResponse struct {
	Total int64          `json:"total"`
	Posts []PostCardTime `json:"posts"`
}

func Posts(context *register.HandleContext) (err error) {
	args := PostsRequest{}
	context.RequestArgs(&args)

	output.Debug("%+v", args)

	res := PostsResponse{}
	posts := make([]PostCardUnix, 10)
	switch args.Type {
	case "tag":
		res.Total, err = mongo.Aggregate("blotter", "posts", []bson.M{
			{"$sort": bson.M{"publish_time": -1}},
			{
				"$lookup": bson.M{
					"from":         "tags",
					"localField":   "tags",
					"foreignField": "_id",
					"as":           "tags",
				},
			},
			{"$set": bson.M{"temp": "$tags.short"}},
			{"$match": bson.M{"temp": args.Arg}},
			{"$limit": args.Offset + args.Number},
			{"$skip": args.Offset},
		}, nil, &posts)
	case "index":
		fallthrough
	default:
		res.Total, err = mongo.Aggregate("blotter", "posts", []bson.M{
			{"$sort": bson.M{"publish_time": -1}},
			{
				"$lookup": bson.M{
					"from":         "tags",
					"localField":   "tags",
					"foreignField": "_id",
					"as":           "tags",
				},
			},
			{"$limit": args.Offset + args.Number},
			{"$skip": args.Offset},
		}, nil, &posts)
	}
	if err != nil {
		return
	}

	res.Posts = make([]PostCardTime, len(posts))
	for idx, post := range posts {
		res.Posts[idx] = post.ToPostCardTime()
	}

	err = context.ReturnJSON(res)
	return
}

// ToPostCardTime transfer PostCardUnix to PostCardTime
func (post PostCardUnix) ToPostCardTime() PostCardTime {
	return PostCardTime{
		Title:       post.Title,
		Abstract:    post.Abstract,
		View:        post.View,
		URL:         post.URL,
		PublishTime: time.ToString(post.PublishTime),
		EditTime:    time.ToString(post.EditTime),
		Tags:        post.Tags,
		HeadImage:   post.HeadImage,
	}
}

// ToPostCardUnix transfer PostCardTime to ToPostCardUnix
func (post PostCardTime) ToPostCardUnix() PostCardUnix {
	return PostCardUnix{
		Title:       post.Title,
		Abstract:    post.Abstract,
		View:        post.View,
		URL:         post.URL,
		PublishTime: time.FromString(post.PublishTime),
		EditTime:    time.FromString(post.EditTime),
		Tags:        post.Tags,
		HeadImage:   post.HeadImage,
	}
}

// ToPostTime transfer PostUnix to PostTime
func (post PostUnix) ToPostTime() PostTime {
	return PostTime{
		Title:       post.Title,
		Abstract:    post.Abstract,
		View:        post.View,
		URL:         post.URL,
		PublishTime: time.ToString(post.PublishTime),
		EditTime:    time.ToString(post.EditTime),
		Tags:        post.Tags,
		HeadImage:   post.HeadImage,
		Content:     post.Content,
	}
}

// ToPostUnix transfer PostTime to ToPostUnix
func (post PostTime) ToPostUnix() PostUnix {
	return PostUnix{
		Title:       post.Title,
		Abstract:    post.Abstract,
		View:        post.View,
		URL:         post.URL,
		PublishTime: time.FromString(post.PublishTime),
		EditTime:    time.FromString(post.EditTime),
		Tags:        post.Tags,
		HeadImage:   post.HeadImage,
		Content:     post.Content,
	}
}
