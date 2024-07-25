package mxorm

// Restorer is an interface that models can implement to allow them to be
// restored to their original state after a database operation.
//
// Restore() is called on a model after it has been loaded from a database
// and allows the model to reset any temporary or computed fields.
//
// The Restorer interface is used by the mxorm.doRestore function to
// restore models that have been loaded from a database.
type Restorer interface {
	// Restore restores the model to its original state after a database
	// operation. This is typically used to reset any temporary or computed
	// fields that were loaded from the database during the load operation.
	Restore()
}
