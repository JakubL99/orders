package handler

import (
	"context"
	pb "orders/proto"

	"github.com/micro/micro/v3/service/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Order struct {
	ID       primitive.ObjectID `bson:"_id"`
	Products Products           `json:"products"`
	Price    string             `json:"price"`
	IdUser   string             `json:"idUser"`
	Name     string             `json:"name"`
	Surname  string             `json:"surname"`
	Address  Address            `json:"address"`
	Status   string             `json:"status"`
}

type Product struct {
	IdProduct string `json:"idProduct"`
	Name      string `json:"name"`
	Price     string `json:"price"`
}

type Products []*Product

type Address struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Post    string `json:"post"`
	Street  string `json:"street"`
	Number  string `json:"number"`
}

type Handler struct {
	Repo
}

type Repo interface {
	Create(ctx context.Context, order *Order) error
}

type MongoRepository struct {
	Collection *mongo.Collection
}

func MarshalOrder(order *pb.Order) *Order {
	return &Order{
		ID:       primitive.NewObjectID(),
		Products: MarshalProductCollection(order.Products),
		Price:    order.Price,
		IdUser:   order.IdUser,
		Name:     order.Name,
		Surname:  order.Surname,
		Address:  MarshalAddress(order.Address),
		Status:   order.Status,
	}
}

func MarshalProduct(product *pb.Product) *Product {
	return &Product{
		IdProduct: product.IdProduct,
		Name:      product.Name,
		Price:     product.Price,
	}
}

func MarshalProductCollection(products []*pb.Product) []*Product {
	collection := make([]*Product, 0)
	for _, product := range products {
		collection = append(collection, MarshalProduct(product))
	}
	return collection
}

func MarshalAddress(address *pb.Address) Address {
	return Address{
		Country: address.Country,
		City:    address.City,
		Post:    address.Post,
		Street:  address.Street,
		Number:  address.Number,
	}
}

func (repo *MongoRepository) Create(ctx context.Context, order *Order) error {
	_, err := repo.Collection.InsertOne(ctx, order)
	a := order.ID
	logger.Info("a: ", a)
	return err
}

func (h *Handler) Create(ctx context.Context, req *pb.Order, rsp *pb.OrderResponse) error {
	Request := MarshalOrder(req)
	NumberOrder := Request.ID

	logger.Info("NumberOrder: ", NumberOrder)

	err := h.Repo.Create(ctx, Request)
	if err != nil {
		logger.Error("Error Create Order: ", err)
	}
	c := primitive.ObjectID.Hex(Request.ID)
	rsp.NumberOrder = c
	logger.Info("rsp", rsp)
	return nil
}
