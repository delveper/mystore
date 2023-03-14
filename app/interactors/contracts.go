// By separating the business logic from the repository and the transport layers,
// the Clean architecture approach allows for greater flexibility and scalability of the application.
// The usecase interactors layer is an essential part of this architecture as it ensures
// that the business rules are implemented consistently and independently of any other part of the application.
// The interactors would call the entities to fulfill a use case, where a use case might be something like create, purchase, order .

package interactors

import (
	"context"

	"github.com/delveper/mystore/app/entities"
)

// The ProductRepo interface is responsible for defining the methods
// that interact with the database or any other storage mechanism
// used by the application to store and retrieve product data.
type ProductRepo interface {
	Insert(context.Context, entities.Product) error
	Select(context.Context, entities.Product) (*entities.Product, error)
	SelectMany(context.Context) ([]entities.Product, error)
	Update(context.Context, entities.Product) error
	Delete(context.Context, entities.Product) error
}