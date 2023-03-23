// By separating the business logic from the repository and the transport layers,
// the Clean architecture approach allows for greater flexibility and scalability of the application.
// The usecase interactors layer is an essential part of this architecture as it ensures
// that the business rules are implemented consistently and independently of any other part of the application.
// The interactors would call the entities to fulfill a use case, where a use case might be something like create, purchase, order .
// It's worth noting that we don't have any business-specific heavy-weight logic from the start.

package interactors

import (
	"context"

	"github.com/delveper/mystore/app/entities"
)

// The ProductRepo interface is responsible for defining the methods
// that interact with the database or any other storage mechanism
// used by the application to store and retrieve product data.
// The implementation of these methods is delegated to the repository layer,
// which can be swapped out and replaced with another implementation
// without affecting the business logic of the application.
type ProductRepo interface {
	Insert(context.Context, entities.Product) (id int, err error)
	Select(context.Context, entities.Product) (*entities.Product, error)
	SelectMany(context.Context) ([]entities.Product, error)
	Update(context.Context, entities.Product) error
	Delete(context.Context, entities.Product) (id int, err error)
}
