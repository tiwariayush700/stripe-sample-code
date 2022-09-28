package repositoryimpl

import (
	`context`
	`log`
	`time`

	`go.mongodb.org/mongo-driver/bson`
	`go.mongodb.org/mongo-driver/mongo`
	`go.mongodb.org/mongo-driver/mongo/options`

	`stripe.com/docs/payments/core/api`
	`stripe.com/docs/payments/core/constant`
	`stripe.com/docs/payments/core/model`
	`stripe.com/docs/payments/repository`
)

type paymentRepositoryImpl struct {
	MongoConn *mongo.Database
}

func (p *paymentRepositoryImpl) Create(ctx context.Context, payment *model.Payment) error {

	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	c := p.MongoConn.Collection(constant.PaymentsMongoCollectionName)

	_, err := c.InsertOne(ctx, payment)

	return err
}

func (p *paymentRepositoryImpl) Get(ctx context.Context, id string) (*model.Payment, error) {
	var payment *model.Payment
	c := p.MongoConn.Collection(constant.PaymentsMongoCollectionName)

	result := c.FindOne(ctx, bson.M{"_id": id})
	err := result.Decode(&payment)
	if err != nil {
		log.Printf("Payments : get err :%v", err)
		if err == mongo.ErrNoDocuments {
			return nil, api.NewHTTPResourceNotFound("payments", id, "No payments found")
		}
		return nil, err
	}

	return payment, nil
}

func (p *paymentRepositoryImpl) Update(ctx context.Context, payment *model.Payment) error {

	c := p.MongoConn.Collection(constant.PaymentsMongoCollectionName)

	payment.UpdatedAt = time.Now()

	selector := bson.M{"_id": bson.M{"$eq": payment.ID}}

	log.Printf("updated payment %+v", *payment)
	update := bson.M{"$set": payment}

	_, err := c.UpdateOne(ctx, selector, update)
	if err != nil {
		log.Printf("Payments : update err :%v", err)
		return err
	}

	return nil
}

func NewPaymentRepositoryImpl(ctx context.Context, mongoUri string, mongoDbName string) (repository.Payment, error) {
	mongoConn, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = mongoConn.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &paymentRepositoryImpl{
		MongoConn: mongoConn.Database(mongoDbName),
	}, nil
}
