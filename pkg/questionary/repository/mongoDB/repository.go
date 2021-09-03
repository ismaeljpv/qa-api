package mongoDB

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	repo "github.com/ismaeljpv/qa-api/pkg/questionary/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func RepositoryError(e string) error {
	return errors.New(e)
}

const (
	DBName                 = "questionary"
	QuestionInfoCollection = "questionInfo"
)

type repository struct {
	db     *mongo.Database
	logger log.Logger
}

func initDBConnection(ctx context.Context, logger log.Logger, uri string) *mongo.Database {

	ctxTO, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	client, connErr := mongo.Connect(ctxTO, options.Client().ApplyURI(uri))
	if connErr != nil {
		level.Error(logger).Log("msg", connErr.Error())
	}

	ctxTO, cancel = context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := client.Ping(ctxTO, readpref.Primary())
	if err != nil {
		level.Warn(logger).Log("msg", err.Error())
	}
	database := client.Database(DBName)
	return database
}

func NewRepository(ctx context.Context, logger log.Logger, uri string) repo.Repository {
	return &repository{
		db:     initDBConnection(ctx, logger, uri),
		logger: logger,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]domain.QuestionInfo, error) {
	var results []domain.QuestionInfo
	QICollection := r.db.Collection(QuestionInfoCollection)
	cursor, err := QICollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return []domain.QuestionInfo{}, RepositoryError(err.Error())
	}

	for cursor.Next(context.TODO()) {
		var questionInfo domain.QuestionInfo
		err := cursor.Decode(&questionInfo)
		if err != nil {
			return []domain.QuestionInfo{}, RepositoryError(err.Error())
		}
		results = append(results, questionInfo)
	}

	if err := cursor.Err(); err != nil {
		return []domain.QuestionInfo{}, RepositoryError(err.Error())
	}
	cursor.Close(context.TODO())

	if len(results) == 0 {
		return []domain.QuestionInfo{}, nil
	}

	return results, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (domain.QuestionInfo, error) {
	var result domain.QuestionInfo
	filter := bson.D{{Key: "question.id", Value: id}}
	QICollection := r.db.Collection(QuestionInfoCollection)

	err := QICollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		level.Warn(r.logger).Log("msg", err.Error())
		return domain.QuestionInfo{}, RepositoryError("No Question Found")
	}
	return result, nil
}

func (r *repository) FindByUser(ctx context.Context, userId string) ([]domain.QuestionInfo, error) {
	var results []domain.QuestionInfo
	filter := bson.D{{Key: "question.userid", Value: userId}}
	QICollection := r.db.Collection(QuestionInfoCollection)

	cursor, err := QICollection.Find(context.TODO(), filter)
	if err != nil {
		return []domain.QuestionInfo{}, RepositoryError(err.Error())
	}

	for cursor.Next(context.TODO()) {
		var questionInfo domain.QuestionInfo
		err := cursor.Decode(&questionInfo)
		if err != nil {
			return []domain.QuestionInfo{}, RepositoryError(err.Error())
		}
		results = append(results, questionInfo)
	}

	if err := cursor.Err(); err != nil {
		return []domain.QuestionInfo{}, RepositoryError(err.Error())
	}
	cursor.Close(context.TODO())

	if len(results) == 0 {
		return []domain.QuestionInfo{}, nil
	}
	return results, nil
}

func (r *repository) Create(ctx context.Context, question domain.Question) (domain.Question, error) {
	QICollection := r.db.Collection(QuestionInfoCollection)
	newQuestionInfo := domain.QuestionInfo{Question: question}
	_, err := QICollection.InsertOne(context.TODO(), newQuestionInfo)
	if err != nil {
		level.Warn(r.logger).Log("msg", err.Error())
		return domain.Question{}, RepositoryError("Error creating the question")
	}

	level.Info(r.logger).Log("msg", fmt.Sprintf("New Question created with ID [%v]", question.ID))
	return question, nil
}

func (r *repository) Update(ctx context.Context, questionInfo domain.QuestionInfo) (domain.QuestionInfo, error) {
	var result domain.QuestionInfo
	filter := bson.D{{Key: "question.id", Value: questionInfo.Question.ID}}
	QICollection := r.db.Collection(QuestionInfoCollection)
	er := QICollection.FindOne(context.TODO(), filter).Decode(&result)
	if er != nil {
		level.Warn(r.logger).Log("msg", er.Error())
		return domain.QuestionInfo{}, RepositoryError("No Question Found")
	}

	if result.Answer.ID == "" {
		return domain.QuestionInfo{}, RepositoryError("Question Has No Answer To Updated")
	}

	if strings.Compare(result.Question.Statement, questionInfo.Question.Statement) != 0 {
		result.Question.Statement = questionInfo.Question.Statement
	}

	if result.Answer.ID == questionInfo.Answer.ID && strings.Compare(result.Answer.Answer, questionInfo.Answer.Answer) != 0 {
		result.Answer.Answer = questionInfo.Answer.Answer
	}

	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "question.statement", Value: result.Question.Statement},
			{Key: "answer.answer", Value: result.Answer.Answer},
		}}}

	updated, err := QICollection.UpdateOne(ctx, filter, update)
	if err != nil {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("There Was An Error Updating The Data Of Question With ID [%v]", questionInfo.Question.ID))
		return domain.QuestionInfo{}, RepositoryError("There Was An Error Updating The Question/Answer")
	}

	if updated.ModifiedCount == 0 {
		return domain.QuestionInfo{}, RepositoryError("The Question/Answer Has No Modifications")
	}

	return result, nil
}

func (r *repository) Delete(ctx context.Context, id string) (string, error) {
	filter := bson.D{{Key: "question.id", Value: id}}
	QICollection := r.db.Collection(QuestionInfoCollection)

	deleted, err := QICollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		level.Warn(r.logger).Log("msg", err.Error())
		return "", RepositoryError("There Was An Error Deleting The Question")
	}

	if deleted.DeletedCount == 0 {
		return "", RepositoryError("No Question Found")
	}
	return "Question Deleted Successfully", nil
}

func (r repository) AddAnswer(ctx context.Context, answer domain.Answer) (domain.QuestionInfo, error) {
	var result domain.QuestionInfo
	filter := bson.D{{Key: "question.id", Value: answer.QuestionID}}
	QICollection := r.db.Collection(QuestionInfoCollection)
	er := QICollection.FindOne(context.TODO(), filter).Decode(&result)
	if er != nil {
		level.Warn(r.logger).Log("msg", er.Error())
		return domain.QuestionInfo{}, RepositoryError("No Question Found")
	}

	if result.Question.ID == "" {
		return domain.QuestionInfo{}, RepositoryError("No Question Found")
	}

	if result.Answer.ID != "" {
		return domain.QuestionInfo{}, RepositoryError("Question Is Already Answered")
	}
	result.Answer = answer
	update := bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "answer", Value: answer},
		}}}

	updated, err := QICollection.UpdateOne(ctx, filter, update)
	if err != nil || updated.ModifiedCount == 0 {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("There Was An Error Adding The Answer To Question With ID [%v]", answer.QuestionID))
		return domain.QuestionInfo{}, RepositoryError("There Was An Error Adding The Answer")
	}
	return result, nil
}
