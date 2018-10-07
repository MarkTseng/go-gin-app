// models.article.go

package main

import (
	"errors"
	"github.com/night-codes/mgo-ai"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	//"log"
)

/*
type article struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Body string `json:"content"`
}
*/

// Article model
type article struct {
	//Id        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Id        uint64 `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string `json:"title" form:"title" binding:"required" bson:"title"`
	Body      string `json:"body" form:"body" binding:"required" bson:"body"`
	CreatedOn int64  `json:"created_on" bson:"created_on"`
	UpdatedOn int64  `json:"updated_on" bson:"updated_on"`
}

// For this demo, we're storing the article list in memory
// In a real application, this list will most likely be fetched
// from a database or from static files
/*
var articleList = []article{
	article{Id: bson.ObjectId("1"), Title: "Article 1", Body: "Article 1 body"},
	article{Id: bson.ObjectId("2"), Title: "Article 2", Body: "Article 2 body"},
}
*/
// Return a list of all the articles
func getAllArticles() []article {
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("articles_demo_dev").C("articles")
	var results []article
	err = c.Find(nil).Sort("-timestamp").All(&results)

	return results
	//return articleList
}

// Fetch an article based on the Id supplied
func getArticleByID(id uint64) (*article, error) {
	/*
		for _, a := range articleList {
			if a.Id == bson.ObjectId(id) {
				return &a, nil
			}
		}
	*/
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("articles_demo_dev").C("articles")
	result := article{}
	err = c.Find(bson.M{"_id": id}).One(&result)
	//log.Print(result)
	if result.Id == id {
		return &result, nil
	}

	return nil, errors.New("Article not found")
}

// Create a new article with the title and content provided
func createNewArticle(title, content, username string) (*article, error) {
	// Set the Id of a new article to one more than the number of articles
	//a := article{Title: title, Body: content}

	// Add the article to the list of articles
	// articleList = append(articleList, a)

	//log.Printf("createNewArticle: %s\n", username)
	session, err := mgo.Dial(MongoDBUrl)
	if err != nil {
		panic(err)
	}

	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	// connect AutoIncrement to collection "counters"
	ai.Connect(session.DB("articles_demo_dev").C("counters"))
	c := session.DB("articles_demo_dev").C("articles")

	aId := ai.Next("articles")
	a := article{Title: title, Body: content, Id: aId}
	err = c.Insert(bson.M{"_id": aId, "title": title, "body": content})
	//log.Print(aId)
	return &a, nil
}
